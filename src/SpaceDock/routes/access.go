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
)

/*
 Registers the routes for the account management
 */
func AccessRegister() {
    Register(GET, "/api/access", listroles)
    Register(POST, "/api/access/roles/assign",
        middleware.NeedsPermission("access-edit", true),
        assignrole,
    )
}

/*
 Path:   /api/access
 Method: GET
 Description: Displays  list of all roles with the matching abilities
 */
func listroles(ctx *iris.Context) {
    var roles []objects.Role
    SpaceDock.Database.Find(&roles)
    output := make([]map[string]interface{}, len(roles))
    for i,element := range roles {
        abilities := element.GetAbilities()
        abilitiesout := make([]map[string]interface{}, len(abilities))
        for j,element2 := range abilities {
            abilitiesout[j] = map[string]interface{} {
                "id": element2.ID,
                "name": element2.Name,
                "meta": utils.LoadJSON(element2.Meta),
            }
        }
        output[i] = map[string]interface{} {
            "id": element.ID,
            "name": element.Name,
            "abilities": abilitiesout,
            "meta": utils.LoadJSON(element.Meta),
        }
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(roles), "data": output})
}

/*
 Path:   /api/access/roles/assign
 Method: POST
 Description: Promotes a user for the given role. Required parameters: userid, rolename
 Abilities: access-edit
 */
func assignrole(ctx *iris.Context) {
    // Grab parameters from the JSON
    userid := utils.GetJSON(ctx,"userid").(uint)
    rolename := utils.GetJSON(ctx,"rolename").(string)

    // Try to get the user
    var user objects.User
    user.GetById(userid)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid.").Code(2145))
    }

    // User is valid, assign the new role
    role := user.AddRole(rolename)
    SpaceDock.Database.Save(&role).Save(&user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}