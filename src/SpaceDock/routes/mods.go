/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "SpaceDock"
    "SpaceDock/middleware"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "os"
    "path/filepath"
    "time"
)

/*
 Registers the routes for the mod section
 */
func ModsRegister() {

}

/*
 Path: /api/mods
 Method: GET
 Description: Returns a list of all mods
 */
func mod_list(ctx *iris.Context) {
    var mods []objects.Mod
    SpaceDock.Database.Find(&mods)
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
    var mods []objects.Mod
    SpaceDock.Database.Find(&mods)
    output := []map[string]interface{}{}
    for _,element := range mods {
        if element.Game.Short == gameshort {
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
func mods_info(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    SpaceDock.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort {
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
func mods_download(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))
    versionname := ctx.GetString("versionname")

    // Get the mod
    mod := &objects.Mod{}
    SpaceDock.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    if !mod.Published && !middleware.IsCurrentUser(ctx, &mod.User) {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("The mod is not published").Code(3020))
        return
    }

    // Get the version
    version := &objects.ModVersion{}
    SpaceDock.Database.Where("friendly_version = ?", versionname).Where("mod_id = ?", mod.ID).First(version)
    if version.FriendlyVersion != versionname {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The version is invalid.").Code(2155))
        return
    }

    // Grab events
    download := &objects.DownloadEvent{}
    SpaceDock.Database.
        Where("mod_id = ?", mod.ID).
        Where("version_id = ?", version.ID).
        Order("download_event.created_at desc").
        First(download)

    // Check whether the path exists
    if _, err := os.Stat(filepath.Join(SpaceDock.Settings.Storage, version.DownloadPath)); os.IsNotExist(err) {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The file you tried to access doesn't exist.").Code(2120))
        return
    }

    if ctx.Header().Get("Range") == "" {
        if download.ID == 0 || (time.Now().Sub(download.CreatedAt).Seconds() / 60 / 60) >= 1 {
            download = objects.NewDownloadEvent(*mod, *version)
            download.Downloads += 1
            SpaceDock.Database.Save(download)
        } else {
            download.Downloads += 1
        }
        mod.DownloadCount += 1
    }
    SpaceDock.Database.Save(mod)

    // Download
    ctx.Redirect("/content/" + version.DownloadPath, iris.StatusTemporaryRedirect)
}

/*
 Path: /api/mods/:gameshort/:modid
 Method: PUT
 Description: Edits a mod, based on the request parameters. Required fields: data
 Abilities: mods-edit
 */
func mods_edit(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    modid := cast.ToUint(ctx.GetString("modid"))

    // Get the mod
    mod := &objects.Mod{}
    SpaceDock.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort {
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
    SpaceDock.Database.Save(mod)

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
    SpaceDock.Database.Where("name = ?", name).First(mod)
    if mod.Name == name {
        errors = append(errors, "A mod with this name already exists.")
        codes = append(codes, 2035)
    }
    game := &objects.Game{}
    SpaceDock.Database.Where("short = ?", gameshort).First(game)
    if game.Short != gameshort {
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
    SpaceDock.Database.Save(mod)
    role := mod.User.AddRole(name)
    role.AddAbility("mods-edit")
    role.AddAbility("mods-remove")
    role.AddParam("mods-edit", "modid", cast.ToString(mod.ID))
    role.AddParam("mods-remove", "name", name)
    SpaceDock.Database.Save(role)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(mod)})
}

/*
 Path: /api/mods
 Method: DELETE
 Description: Removes a mod, based on the request parameters. Required fields: name, gameshort
 Abilities: mods-remove
 */
func mods_remove(ctx *iris.Context) {
    // Get params
    gameshort := utils.GetJSON(ctx, "gameshort")
    modid := cast.ToUint(utils.GetJSON(ctx,"modid"))

    // Get the mod
    mod := &objects.Mod{}
    SpaceDock.Database.Where("id = ?", modid).First(mod)
    if mod.ID != modid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The modid is invalid").Code(2130))
        return
    }
    if mod.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    // Delete the mod
    role := &objects.Role{}
    SpaceDock.Database.Where("name = ?", mod.Name).First(role)
    role.RemoveAbility("mods-edit")
    role.RemoveAbility("mods-remove")
    mod.User.RemoveRole(mod.Name)
    SpaceDock.Database.Delete(mod).Delete(role)

    // Display info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}
