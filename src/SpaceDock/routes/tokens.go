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
    "strconv"
)

/*
 Registers the routes for the token section
 */
func TokensRegister() {
    Register(POST, "/api/tokens",
        middleware.NeedsPermission("token-generate", true),
        generate_token,
    )
    Register(PUT, "/api/tokens",
        middleware.NeedsPermission("token-edit", true, "tokenid"),
        edit_token,
    )
    Register(DELETE, "/api/tokens",
        middleware.NeedsPermission("token-revoke", true, "tokenid"),
        revoke_token,
    )
}

/*
 Path: /api/tokens
 Method: POST
 Description: Generates a new API token
 Abilities: token-generate
 */
func generate_token(ctx *iris.Context) {
    token := objects.NewToken()
    token.SetValue("ips", []string{})
    SpaceDock.Database.Save(token)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(token)})
}

/*
 Path: /api/tokens
 Method: PUT
 Description: Edits the IP-Adresses of a token
 Abilities: token-edit
 */
func edit_token(ctx *iris.Context) {
    tokenid := cast.ToUint(utils.GetJSON(ctx, "tokenid"))
    ips := cast.ToStringSlice(utils.GetJSON(ctx, "ips"))

    // Get the token
    var token objects.Token
    SpaceDock.Database.Where("id = ?", tokenid).First(&token)
    if token.ID != tokenid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The token ID is invalid").Code(2131))
        return
    }
    if ips == nil {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The list of IP Addresses is invalid.").Code(2132))
        return
    }

    // Edit the token
    token.SetValue("ips", ips)
    SpaceDock.Database.Save(&token)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(token)})
}

/*
 Path: /api/tokens
 Method: DELETE
 Description: Removes a token completely
 Abilities: token-revoke
 */
func revoke_token(ctx *iris.Context) {
    tokenid := cast.ToUint(utils.GetJSON(ctx, "tokenid"))

    // Get the token
    var token objects.Token
    SpaceDock.Database.Where("id = ?", tokenid).First(&token)
    if token.ID != tokenid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The token ID is invalid").Code(2131))
        return
    }

    // Delete the token
    SpaceDock.Database.Delete(&token)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}