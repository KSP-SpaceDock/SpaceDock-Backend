/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package routes

import (
    "SpaceDock"
    "SpaceDock/middleware"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "github.com/kataras/iris"
    "strconv"
)

/*
 Registers the routes for the admin section
 */
func AdminRegister() {
    Register(GET, "/api/admin/impersonate/:userid",
        middleware.NeedsPermission("admin-impersonate", true, "userid"),
        impersonate,
    )
    Register(GET, "/api/admin/manual-confirmation/:userid",
        middleware.NeedsPermission("admin-confirm", true),
        manualConfirmation,
    )
}

/*
 Path: /api/admin/impersonate/:userid
 Method: GET
 Description: Log into another persons account from an admin account
 */
func impersonate(ctx *iris.Context) {
    id, err := ctx.GetInt("userid")
    userid := uint(id)
    if err != nil {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid").Code(2145))
        return
    }
    var user objects.User
    user.GetById(userid)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid").Code(2145))
        return
    }
    middleware.LogoutUser(ctx)
    middleware.LoginUser(ctx, user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/admin/manual-confirmation/:userid
 Method: GET
 */
func manualConfirmation(ctx *iris.Context) {
    id, err := ctx.GetInt("userid")
    userid := uint(id)
    if err != nil {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid").Code(2145))
        return
    }
    var user objects.User
    user.GetById(userid)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid").Code(2145))
        return
    }

    // Everything is valid
    user.Confirmation = ""
    role := user.AddRole(user.Username)
    role.AddAbility("user-edit")
    role.AddAbility("mods-add")
    role.AddAbility("packs-add")
    role.AddAbility("logged-in")
    role.AddParam("user-edit", "userid", strconv.Itoa(int(user.ID)))
    role.AddParam("mods-add", "gameshort", ".*")
    role.AddParam("packs-add", "gameshort", ".*")
    SpaceDock.Database.Save(&role)
    SpaceDock.Database.Save(&user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}