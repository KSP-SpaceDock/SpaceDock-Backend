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
)

/*
 Registers the routes for the modlist section
 */
func ModlistsRegister() {
    Register(GET, "/api/lists", list_modlists)
    Register(GET, "/api/lists/:gameshort", list_modlists_game)
    Register(GET, "/api/lists/:gameshort/:listid", list_info)
    Register(POST, "/api/lists",
        middleware.NeedsPermission("lists-add", true, "gameshort"),
        lists_add,
    )
    Register(PUT, "/api/lists/:gameshort/:listid",
        middleware.NeedsPermission("lists-edit", true, "gameshort", "listid"),
        list_edit,
    )
    Register(DELETE, "/api/lists/:gameshort",
        middleware.NeedsPermission("lists-remove", true, "gameshort", "listid"),
        list_remove,
    )
}

/*
 Path: /api/lists/
 Method: GET
 Description: Outputs a list of modlists
 */
func list_modlists(ctx *iris.Context) {
    var modlists []objects.ModList
    SpaceDock.Database.Find(&modlists)
    output := make([]map[string]interface{}, len(modlists))
    for i,element := range modlists {
        output[i] = utils.ToMap(element)
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(output), "data": output})
}

/*
 Path: /api/lists/:gameshort
 Method: GET
 Description: Outputs a list of modlists for a specific game
 */
func list_modlists_game(ctx *iris.Context) {
    gameshort := ctx.GetString("gameshort")
    var modlists []objects.ModList
    SpaceDock.Database.Find(&modlists)
    output := []map[string]interface{}{}
    for _,element := range modlists {
        if element.Game.Short == gameshort {
            output = append(output, utils.ToMap(element))
        }
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(output), "data": output})
}

/*
 Path: /api/lists/:gameshort/:listid
 Method: GET
 Description: Returns info for a specific modlist
 */
func list_info(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    listid := cast.ToUint(ctx.GetString("listid"))

    // Get the modlist
    modlist := &objects.ModList{}
    SpaceDock.Database.Where("id = ?", listid).First(modlist)
    if modlist.ID != listid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The pack ID is invalid").Code(2135))
        return
    }
    if modlist.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(modlist)})
}

/*
 Path: /api/lists
 Method: POST
 Description: Creates a new modlist. Required fields: name, gameshort
 */
func lists_add(ctx *iris.Context) {
    // Get params
    name := cast.ToString(utils.GetJSON(ctx, "name"))
    gameshort := cast.ToString(utils.GetJSON(ctx, "gameshort"))

    // Check the vars
    errors := []string{}
    codes := []int{}
    if name == "" {
        errors = append(errors, "Invalid modlist name.")
        codes = append(codes, 2137)
    }
    modlist := &objects.ModList{}
    SpaceDock.Database.Where("name = ?", name).First(modlist)
    if modlist.Name == name {
        errors = append(errors, "A modlist with this name already exists.")
        codes = append(codes, 2025)
    }
    game := &objects.Game{}
    SpaceDock.Database.Where("short = ?", gameshort).First(game)
    if game.Short != gameshort {
        errors = append(errors, "Invalid gameshort.")
        codes = append(codes, 2125)
    }
    if len(errors) > 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
        return
    }

    // Make a new list
    user := middleware.CurrentUser(ctx)
    modlist = objects.NewModList(name, *user, *game)
    SpaceDock.Database.Save(modlist)
    role := user.AddRole(name)
    role.AddAbility("lists-edit")
    role.AddAbility("lists-remove")
    role.AddParam("lists-edit", "listsid", cast.ToString(modlist.ID))
    role.AddParam("lists-remove", "name", name)
    SpaceDock.Database.Save(user).Save(role)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(modlist)})
}

/*
 Path: /api/lists/:gameshort/:listid
 Method: PUT
 Description: Edits a modlist based on patch data. Required fields: data
 */
func list_edit(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    listid := cast.ToUint(ctx.GetString("listid"))

    // Get the modlist
    modlist := &objects.ModList{}
    SpaceDock.Database.Where("id = ?", listid).First(modlist)
    if modlist.ID != listid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The pack ID is invalid").Code(2135))
        return
    }
    if modlist.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    // Edit the modlist
    code := utils.EditObject(modlist, utils.GetFullJSON(ctx))
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
    SpaceDock.Database.Save(modlist)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(modlist)})
}

/*
 Path: /api/lists/:gameshort
 Method: DELETE
 Description: Removes a modlist. Required parameters: listid
 */
func list_remove(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    listid := cast.ToUint(utils.GetJSON(ctx,"listid"))

    // Get the modlist
    modlist := &objects.ModList{}
    SpaceDock.Database.Where("id = ?", listid).First(modlist)
    if modlist.ID != listid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The pack ID is invalid").Code(2135))
        return
    }
    if modlist.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

    // Remove the modlist
    SpaceDock.Database.Delete(modlist)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/lists/:gameshort/:listid/mods
 Method: POST
 Description: Adds a new mod to the modlist. Required fields: modid
 */
func list_add_mod(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    listid := cast.ToUint(ctx.GetString("listid"))
    modid := cast.ToUint(utils.GetJSON(ctx, "modid"))

    // Get the modlist
    modlist := &objects.ModList{}
    SpaceDock.Database.Where("id = ?", listid).First(modlist)
    if modlist.ID != listid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The pack ID is invalid").Code(2135))
        return
    }
    if modlist.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

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

    // Check if the mod is already in the list
    entry := &objects.ModListItem{}
    SpaceDock.Database.Where("mod_list_id = ?", modlist.ID).Where("mod_id = ?", mod.ID).First(entry)
    if entry.ModListID == modlist.ID && entry.ModID == mod.ID {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The mod is already added to the modlist.").Code(2070))
        return
    }

    // Create an entry
    entry = objects.NewModListItem(*mod, *modlist)
    SpaceDock.Database.Save(entry)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/lists/:gameshort/:listid/mods
 Method: DELETE
 Description: Removes a mod from the modlist. Required fields: modid
 */
func list_remove_mod(ctx *iris.Context) {
    // Get params
    gameshort := ctx.GetString("gameshort")
    listid := cast.ToUint(ctx.GetString("listid"))
    modid := cast.ToUint(utils.GetJSON(ctx, "modid"))

    // Get the modlist
    modlist := &objects.ModList{}
    SpaceDock.Database.Where("id = ?", listid).First(modlist)
    if modlist.ID != listid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The pack ID is invalid").Code(2135))
        return
    }
    if modlist.Game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }

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

    // Check if the mod is already in the list
    entry := &objects.ModListItem{}
    SpaceDock.Database.Where("mod_list_id = ?", modlist.ID).Where("mod_id = ?", mod.ID).First(entry)
    if entry.ModListID != modlist.ID && entry.ModID != mod.ID {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The mod is not added to the modlist.").Code(2075))
        return
    }

    // Create an entry
    SpaceDock.Database.Delete(entry)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}