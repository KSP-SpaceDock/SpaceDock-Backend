/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/middleware"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/objects"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "archive/zip"
    "github.com/kennygrant/sanitize"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "io"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "time"
)

/*
 Registers the routes for the mod section
 */
func ModsRegister() {
    Register(GET, "/api/mods", middleware.Recursion(0), middleware.Cache, mod_list)
    Register(GET, "/api/mods/:gameshort", middleware.Recursion(0), middleware.Cache, mod_game_list)
    Register(GET, "/api/mods/:gameshort/:modid", middleware.Cache, mod_info)
    Register(GET, "/api/mods/:gameshort/:modid/download/:versionname", mod_download)
    Register(PUT, "/api/mods/:gameshort/:modid",
        middleware.NeedsPermission("mods-edit", true, "gameshort", "modid"),
        mod_edit,
    )
    Register(POST, "/api/mods",
        middleware.NeedsPermission("mods-add", true),
        mod_add,
    )
    Register(DELETE, "/api/mods",
        middleware.NeedsPermission("mods-remove", true, "gameshort", "modid"),
        mod_remove,
    )
    Register(POST, "/api/mods/:gameshort/:modid/publish",
        middleware.NeedsPermission("mods-edit", true, "gameshort", "modid"),
        mod_publish,
    )
    Register(GET, "/api/mods/:gameshort/:modid/versions", middleware.Recursion(0), middleware.Cache, mod_versions)
    Register(POST, "/api/mods/:gameshort/:modid/versions",
        middleware.NeedsPermission("mod-edit", true, "gameshort", "modid"),
        mod_update,
    )
    Register(DELETE, "/api/mods/:gameshort/:modid/versions",
        middleware.NeedsPermission("mod-edit", true, "gameshort", "modid"),
        mod_version_delete,
    )
    Register(GET, "/api/mods/:gameshort/:modid/follow",
        middleware.NeedsPermission("logged-in", false),
        mod_follow,
    )
    Register(GET, "/api/mods/:gameshort/:modid/unfollow",
        middleware.NeedsPermission("logged-in", false),
        mod_unfollow,
    )
    Register(POST, "/api/mods/:gameshort/:modid/ratings",
        middleware.NeedsPermission("logged-in", false),
        mod_rate,
    )
    Register(DELETE, "/api/mods/:gameshort/:modid/ratings",
        middleware.NeedsPermission("logged-in", false),
        mod_unrate,
    )
    Register(POST, "/api/mods/:gameshort/:modid/grant",
        middleware.NeedsPermission("logged-in", true),
        mod_grant,
    )
    Register(POST, "/api/mods/:gameshort/:modid/grant/accept",
        middleware.NeedsPermission("logged-in", true),
        mod_accept_grant,
    )
    Register(POST, "/api/mods/:gameshort/:modid/grant/reject",
        middleware.NeedsPermission("logged-in", true),
        mod_reject_grant,
    )
    Register(DELETE, "/api/mods/:gameshort/:modid/grant",
        middleware.NeedsPermission("logged-in", true),
        mod_revoke_grant,
    )
}

/*
 Path: /api/mods
 Method: GET
 Description: Returns a list of all mods
 */
func mod_list(ctx *iris.Context) {
    var mods []objects.Mod
    app.Database.Find(&mods)
    output := make([]map[string]interface{}, len(mods))
    for i,element := range mods {
        output[i] = utils.ToMap(element)
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(output), "data": output})
}

/*
 Path: /api/mods/:gameshort
 Method: GET
 Description: Returns a list with all mods for this game.
 */
func mod_game_list(ctx *iris.Context) {
    gameshort := ctx.GetString("gameshort")
    game := &objects.Game{}
    app.Database.Where("short = ?", gameshort).Or("id = ?", cast.ToUint(gameshort)).Find(game)
    var mods []objects.Mod
    app.Database.Find(&mods)
    output := []map[string]interface{}{}
    for _,element := range mods {
        if element.GameID == game.ID {
            output = append(output, utils.ToMap(element))
        }
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(output), "data": output})
}

/*
 Path: /api/mods/:gameshort/:modid
 Method: GET
 Description: Returns information for one mod
 */
func mod_info(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }
    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(mod)})
}

/*
 Path: /api/mods/:gameshort/:modid/download/:versionname
 Method: GET
 Description: Downloads the latest non-beta version of the mod
 */
func mod_download(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    versionname := ctx.GetString("versionname")

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }
    if versionname == "default" {
        versionname = mod.DefaultVersion.FriendlyVersion
    }
    if versionname == "latest" {
        versionname = mod.Versions[len(mod.Versions) - 1].FriendlyVersion
    }

    // Get the version
    version := &objects.ModVersion{}
    app.Database.Where("friendly_version = ?", versionname).Where("mod_id = ?", mod.ID).First(version)
    if version.FriendlyVersion != versionname {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The version is invalid.").Code(2155))
        return
    }

    // Grab events
    download := &objects.DownloadEvent{}
    app.Database.
        Where("mod_id = ?", mod.ID).
        Where("version_id = ?", version.ID).
        Order("download_event.created_at desc").
        First(download)

    // Check whether the path exists
    if _, err := os.Stat(filepath.Join(app.Settings.Storage, version.DownloadPath)); os.IsNotExist(err) {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The file you tried to access doesn't exist.").Code(2120))
        return
    }

    if ctx.Header().Get("Range") == "" {
        if download.ID == 0 || (time.Now().Sub(download.CreatedAt).Seconds() / 60 / 60) >= 1 {
            download = objects.NewDownloadEvent(*mod, *version)
            download.Downloads += 1
            app.Database.Save(download)
        } else {
            download.Downloads += 1
        }
        mod.DownloadCount += 1
    }
    app.Database.Save(mod)

    // Download
    ctx.Redirect("/content/" + version.DownloadPath, iris.StatusTemporaryRedirect)
}

/*
 Path: /api/mods/:gameshort/:modid
 Method: PUT
 Description: Edits a mod, based on the request parameters. Required fields: data
 Abilities: mods-edit
 */
func mod_edit(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    // Edit the mod
    code := utils.EditObject(mod, utils.GetFullJSON(ctx))
    if code == 3 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The value you submitted is invalid").Code(2180))
        return
    } else if code == 2 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that doesn't exist.").Code(3090))
        return
    } else if code == 1 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that is marked as read-only.").Code(3095))
        return
    }
    app.Database.Save(mod)
    utils.ClearModCache(gameshort, modid)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(mod)})
}

/*
 Path: /api/mods
 Method: POST
 Description: Adds a mod, based on the request parameters. Required fields: name, gameshort, license
 Abilities: mods-add
 */
func mod_add(ctx *iris.Context) {
    // Get params
    name := cast.ToString(utils.GetJSON(ctx, "name"))
    gameshort := cast.ToString(utils.GetJSON(ctx, "gameshort"))
    license := cast.ToString(utils.GetJSON(ctx, "license"))

    // Check the vars
    errors := []string{}
    codes := []int{}
    if name == "" {
        errors = append(errors, "Invalid mod name.")
        codes = append(codes, 2117)
    }
    mod := &objects.Mod{}
    app.Database.Where("name = ?", name).First(mod)
    if mod.Name == name {
        errors = append(errors, "A mod with this name already exists.")
        codes = append(codes, 2035)
    }
    game := &objects.Game{}
    app.Database.Where("short = ?", gameshort).Or("id = ?", cast.ToUint(gameshort)).First(game)
    if game.Short != gameshort && game.ID != cast.ToUint(gameshort) {
        errors = append(errors, "Invalid gameshort.")
        codes = append(codes, 2125)
    }
    if license == "" {
        errors = append(errors, "Invalid License.")
        codes = append(codes, 2190)
    }
    if len(errors) > 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
        return
    }

    // Add new mod
    mod = objects.NewMod(name, *middleware.CurrentUser(ctx), *game, license)
    app.Database.Save(mod)
    role := mod.User.AddRole(name)
    role.AddAbility("mods-edit")
    role.AddAbility("mods-remove")
    role.AddParam("mods-edit", "modid", cast.ToString(mod.ID))
    role.AddParam("mods-remove", "name", name)
    app.Database.Save(role)
    utils.ClearModCache(gameshort, 0)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(mod)})
}

/*
 Path: /api/mods
 Method: DELETE
 Description: Removes a mod, based on the request parameters. Required fields: modid, gameshort
 Abilities: mods-remove
 */
func mod_remove(ctx *iris.Context) {
    // Get params
    gameshort := utils.GetJSON(ctx, "gameshort")
    modid := cast.ToUint(utils.GetJSON(ctx,"modid"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    // Delete the mod
    role := &objects.Role{}
    app.Database.Where("name = ?", mod.Name).First(role)
    role.RemoveAbility("mods-edit")
    role.RemoveAbility("mods-remove")
    mod.User.RemoveRole(mod.Name)
    utils.ClearModCache(mod.Game.Short, 0)
    app.Database.Delete(mod).Delete(role)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/mods/:gameshort/:modid/publish
 Method: POST
 Description: Makes a mod public.
 Abilities: mods-edit
 */
func mod_publish(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    // Publish
    mod.Published = true
    app.Database.Save(mod)
    utils.ClearModCache(gameshort, modid)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/mods/:gameshort/:modid/versions
 Method: GET
 Description: Returns a list of mod versions including their data.
 */
func mod_versions(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }

    // Get the mod versions
    output := []map[string]interface{}{}
    for _,element := range mod.Versions {
        output = append(output, utils.ToMap(element))
    }

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(output), "data": output})
}

/*
 Path: /api/mods/:gameshort/:modid/versions
 Method: POST
 Description: Releases a new version of your mod. Required fields: version, game-version, notify-followers, is-beta, zipball. Optional fields: changelog
 */
func mod_update(ctx *iris.Context) {
    // Get params
    version := cast.ToString(utils.GetJSON(ctx, "version"))
    changelog := cast.ToString(utils.GetJSON(ctx, "changelog"))
    friendly_version := cast.ToString(utils.GetJSON(ctx, "game-version"))
    notify := cast.ToBool(utils.GetJSON(ctx, "notify-followers"))
    beta := cast.ToBool(utils.GetJSON(ctx, "is-beta"))
    zipball, _, err := ctx.FormFile("zipball")

    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }

    // Process fields
    if version == "" || friendly_version == "" || err != nil {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("All fields are required.").Code(2505))
        return
    }
    if friendly_version == "default" || friendly_version == "latest" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You cannot use a reserved friendly_version").Code(2503))
    }
    game_version := &objects.GameVersion{}
    app.Database.Where("friendly_version = ?", friendly_version).First(game_version)
    if game_version.FriendlyVersion != friendly_version {
        utils.WriteJSON(ctx, iris.StatusNotFound,  utils.Error("Game version does not exist").Code(2105))
        return
    }

    // Save the file
    user := middleware.CurrentUser(ctx)
    filename := sanitize.BaseName(mod.Name) + "-" + sanitize.BaseName(version) + ".zip"
    base_path := filepath.Join(sanitize.BaseName(user.Username) + "_" + strconv.Itoa(int(user.ID)), sanitize.BaseName(mod.Name))
    full_path := filepath.Join(app.Settings.Storage, base_path)
    os.MkdirAll(full_path, os.ModePerm)
    path := filepath.Join(full_path, filename)
    for _,v := range mod.Versions {
        if v.FriendlyVersion == sanitize.BaseName(version) {
            utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("We already have this version. Did you mistype the version number?").Code(3040))
            return
        }
    }

    // Remove the old file. If it fails, dont care
    _ = os.Remove(filepath.Join(app.Settings.Storage, path))
    out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        utils.WriteJSON(ctx, iris.StatusInternalServerError, utils.Error(err.Error()).Code(2153))
        return
    }
    io.Copy(out, zipball)
    out.Close()
    zipball.Close()

    // Check if the file is a zipfile
    temp,err := zip.OpenReader(path)
    if err != nil {
        _ = os.Remove(filepath.Join(app.Settings.Storage, path))
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("This is not a valid zip file.").Code(2160))
        return
    } else {
        temp.Close()
    }
    modversion := objects.NewModVersion(*mod, sanitize.BaseName(version), *game_version, strings.Replace(filepath.Join(base_path, filename), "\\", "/", -1), beta)
    modversion.Changelog = changelog

    // sort index
    if len(mod.Versions) == 0 {
        modversion.SortIndex = 0
    } else {
        for _, v := range mod.Versions {
            if v.SortIndex > modversion.SortIndex {
                modversion.SortIndex = v.SortIndex
            }
        }
        modversion.SortIndex += 1
    }
    if notify && !beta {
        followers := []string{}
        for _,e := range mod.Followers {
            followers = append(followers, e.Email)
        }
        err, modURL := mod.Game.GetValue("modURL")
        if err != nil {
            modURL = ""
        }
        utils.SendUpdateNotification(followers, changelog, user.Username, friendly_version, mod.Name, mod.ID, cast.ToString(modURL), mod.Game.Name, game_version.FriendlyVersion)
    }
    app.Database.Save(modversion)
    if !beta {
        mod.DefaultVersionID = modversion.ID
        mod.DefaultVersion = *modversion
    }
    app.Database.Save(mod)
    utils.ClearModCache(gameshort, modid)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(modversion)})
}

/*
 Path: /api/mods/:gameshort/:modid/versions
 Method: DELETE
 Description: Deletes a released version of the mod. Required fields: version-id
 */
func mod_version_delete(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    versionID := cast.ToUint(utils.GetJSON(ctx, "version-id"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }

    // Get the version
    version := &objects.ModVersion{}
    app.Database.Where("id = ?", versionID).Where("mod_id = ?", mod.ID).First(version)
    if version.ID != versionID {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The version is invalid.").Code(2155))
        return
    }

    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The mod is not published.").Code(3020))
        return
    }

    // Checks
    if len(mod.Versions) == 1 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("There is only one version left. You cant delete this one.").Code(3025))
        return
    }
    if version.ID == mod.DefaultVersionID {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You cannot delete the default version of a mod.").Code(3080))
        return
    }
    app.Database.Delete(version)
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/mods/:gameshort/:modid/follow
 Method: GET
 Description: Registers a user for automated email sending when a new mod version is released
 */
func mod_follow(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    user := middleware.CurrentUser(ctx)

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }
    if e, _ := utils.ArrayContains(mod, user.Following); e {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You are already following this mod.").Code(3050))
        return
    }

    // Grab events
    follow := &objects.FollowEvent{}
    app.Database.
        Where("mod_id = ?", mod.ID).
        Order("follow_event.created_at desc").
        First(follow)

    if follow.ID == 0 || (time.Now().Sub(follow.CreatedAt).Seconds()/60/60) >= 1 {
        follow = objects.NewFollowEvent(*mod)
        follow.Delta += 1
        follow.Events = 1
        app.Database.Save(follow)
    } else {
        follow.Delta += 1
        follow.Events += 1
    }
    mod.Followers = append(mod.Followers, *user)
    user.Following = append(user.Following, *mod)
    app.Database.Save(mod).Save(user)
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/mods/:gameshort/:modid/unfollow
 Method: GET
 Description: Unregisters a user for automated email sending when a new mod version is released
 */
func mod_unfollow(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    user := middleware.CurrentUser(ctx)

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }
    if e,_ := utils.ArrayContains(mod, user.Following); !e {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You are not following this mod.").Code(3050))
        return
    }

    // Grab events
    follow := &objects.FollowEvent{}
    app.Database.
        Where("mod_id = ?", mod.ID).
        Order("follow_event.created_at desc").
        First(follow)

    if follow.ID == 0 || (time.Now().Sub(follow.CreatedAt).Seconds()/60/60) >= 1 {
        follow = objects.NewFollowEvent(*mod)
        follow.Delta -= 1
        follow.Events = 1
        app.Database.Save(follow)
    } else {
        follow.Delta -= 1
        follow.Events += 1
    }
    _,i := utils.ArrayContains(user, mod.Followers)
    _,j := utils.ArrayContains(mod, user.Following)
    mod.Followers = append(mod.Followers[:i], mod.Followers[i+1:]...)
    user.Following = append(user.Following[:j], user.Following[j+1:]...)
    app.Database.Save(mod).Save(user)
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/mods/:gameshort/:modid/ratings
 Method: POST
 Description: Rates a mod. Required fields: rating
 */
func mod_rate(ctx *iris.Context) {
    // Get variables
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    score := cast.ToFloat64(utils.GetJSON(ctx, "rating"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }
    rating := &objects.Rating{}
    user := middleware.CurrentUser(ctx)
    app.Database.Where("mod_id = ?", modid).Where("user_id = ?", user.ID).First(rating)
    if rating.UserID == user.ID && rating.ModID == modid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You already have a rating for this mod.").Code(2040))
        return
    }

    // Create a rating
    rating = objects.NewRating(*user, *mod, score)
    app.Database.Save(rating)

    // Add rating to user and mod
    mod.Ratings = append(mod.Ratings, *rating)
    mod.CalculateScore()
    app.Database.Save(mod)
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/mods/:gameshort/:modid/ratings
 Method: DELETE
 Description: Removes a rating for a mod.
 */
func mod_unrate(ctx *iris.Context) {
    // Get variables
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }
    rating := &objects.Rating{}
    user := middleware.CurrentUser(ctx)
    app.Database.Where("mod_id = ?", modid).Where("user_id = ?", user.ID).First(rating)
    if rating.UserID != user.ID {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You can't remove a rating you don't have, right?").Code(3013))
        return
    }

    // Remove the rating
    _,i := utils.ArrayContains(*rating, mod.Ratings)
    mod.Ratings = append(mod.Ratings[:i], mod.Ratings[i+1:]...)
    mod.CalculateScore()
    app.Database.Save(mod)
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/mods/:gameshort/:modid/grant
 Method: POST
 Description: Adds a new author to a mod. Required fields: username
 */
func mod_grant(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    username := cast.ToString(utils.GetJSON(ctx, "username"))

    // Check the params
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("ou have to pubish your mod in order to add contributors.").Code(3043))
        return
    }
    user := &objects.User{}
    app.Database.Where("username = ?", username).First(user)
    if user.Username != username {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The username is invalid").Code(2150))
        return
    }
    shared := &objects.SharedAuthor{}
    app.Database.Where("mod_id = ?", mod.ID).Where("user_id = ?", user.ID).First(shared)
    if mod.User.ID == user.ID && shared.UserID == user.ID {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("This user has already been added.").Code(2010))
        return
    }
    if !user.Public {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("This user has not made their profile public.").Code(3040))
        return
    }
    current := middleware.CurrentUser(ctx)
    if mod.User.ID != current.ID && middleware.UserHasPermission(ctx, "mods-invite", true, []string{}) != 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You dont have the permission to add new authors.").Code(1025))
        return
    }

    // Add SharedAuthor
    shared = objects.NewSharedAuthor(*user, *mod)
    mod.SharedAuthors = append(mod.SharedAuthors, *shared)
    app.Database.Save(shared)
    err, modURL := mod.Game.GetValue("modURL")
    if err != nil {
        modURL = ""
    }
    utils.SendGrantNotice(user.Username, mod.User.Username, mod.Name, mod.ID, user.Email, cast.ToString(modURL))

    // Display info
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(mod)})
}

/*
 Path: /api/mods/:gameshort/:modid/grant/accept
 Method: POST
 Description: Accepts a pending authorship grant for a mod.
 */
func mod_accept_grant(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }

    for _,element := range mod.SharedAuthors {
        if middleware.IsCurrentUser(ctx, &element.User) && !element.Accepted {
            element.Accepted = true
            element.User.AddRole(mod.Name)
            app.Database.Save(&(element.User)).Save(&element)
            utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(mod)})
            return
        }
    }
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You do not have a pending authorship invite.").Code(3085))
}

/*
 Path: /api/mods/:gameshort/:modid/grant/reject
 Method: POST
 Description: Rejects a pending authorship grant for a mod.
 */
func mod_reject_grant(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    // Get the mod
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }

    for _,element := range mod.SharedAuthors {
        if middleware.IsCurrentUser(ctx, &element.User) && !element.Accepted {
            app.Database.Delete(&element)
            utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
            return
        }
    }
    utils.ClearModCache(gameshort, modid)
    utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You do not have a pending authorship invite.").Code(3085))
}

/*
 Path: /api/mods/:gameshort/:modid/grant
 Method: DELETE
 Description: Removes an author from a mod. Required fields: username
 */
func mod_revoke_grant(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    username := cast.ToString(utils.GetJSON(ctx, "username"))

    // Check the params
    mod := &objects.Mod{}
    app.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort && mod.GameID != cast.ToUint(gameshort) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    user := &objects.User{}
    app.Database.Where("username = ?", username).First(user)
    if user.Username != username {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The username is invalid").Code(2150))
        return
    }
    shared := &objects.SharedAuthor{}
    app.Database.Where("mod_id = ?", mod.ID).Where("user_id = ?", user.ID).First(shared)
    if shared.ModID != mod.ID && shared.UserID != user.ID {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("This user is not an author.").Code(3073))
        return
    }
    if mod.User.ID == user.ID {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You can't remove this user.").Code(3075))
        return
    }
    current := middleware.CurrentUser(ctx)
    if mod.User.ID != current.ID && middleware.UserHasPermission(ctx, "mods-invite", true, []string{}) != 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You dont have the permission to remove authors.").Code(1030))
        return
    }

    // Remove SharedAuthor
    shared.User.RemoveRole(mod.Name)
    app.Database.Save(&(shared.User))
    _, i := utils.ArrayContains(shared, mod.SharedAuthors)
    mod.SharedAuthors = append(mod.SharedAuthors[:i], mod.SharedAuthors[i+1:]...)
    app.Database.Delete(shared)
    utils.ClearModCache(gameshort, modid)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}
