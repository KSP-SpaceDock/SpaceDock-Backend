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
    "gopkg.in/kataras/iris.v6"
)

/*
 Registers the routes for the account management
 */
func AccessRegister() {
    Register(GET, "/api/access",
        middleware.NeedsPermission("access-view", true),
        listroles,
    )
    Register(POST, "/api/access/roles/assign",
        middleware.NeedsPermission("access-edit", true),
        assignrole,
    )
    Register(POST, "/api/access/roles/remove",
        middleware.NeedsPermission("access-edit", true),
        removerole,
    )
    Register(GET, "/api/access/abilities",
        middleware.NeedsPermission("access-view", true),
        listabilities,
    )
    Register(POST, "/api/access/abilities/assign",
        middleware.NeedsPermission("access-edit", true),
        assignability,
    )
    Register(POST, "/api/access/abilities/remove",
        middleware.NeedsPermission("access-edit", true),
        removeability,
    )
    Register(POST, "/api/access/params/add/:rolename",
        middleware.NeedsPermission("access-edit", true),
        addparam,
    )
    Register(POST, "/api/access/params/remove/:rolename",
        middleware.NeedsPermission("access-edit", true),
        removeparam,
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



/*
 Path:   /api/access/roles/remove
 Method: POST
 Description: Promotes a user for the given role. Required parameters: userid, rolename
 Abilities: access-edit
 */
func removerole(ctx *iris.Context) {
    // Grab parameters from the JSON
    userid := utils.GetJSON(ctx,"userid").(uint)
    rolename := utils.GetJSON(ctx,"rolename").(string)

    // Try to get the user
    var user objects.User
    user.GetById(userid)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid.").Code(2145))
    }
    if user.HasRole(rolename) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The user doesn't have this role").Code(1015))
    }

    // Everything is valid, remove the role
    user.RemoveRole(rolename)
    SpaceDock.Database.Save(&user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path:   /api/access/abilities
 Method: GET
 Description: Displays a list of all abilities.
 */
func listabilities(ctx *iris.Context) {
    var abilities []objects.Ability
    SpaceDock.Database.Find(&abilities)
    output := make([]map[string]interface{}, len(abilities))
    for i,element := range abilities {
        output[i] = map[string]interface{} {
            "id": element.ID,
            "name": element.Name,
            "meta": utils.LoadJSON(element.Meta),
        }
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(abilities), "data": output})
}

/*
 Path:   /api/access/abilities/assign
 Method: POST
 Description: Adds a permission to a group. Required parameters: rolename, abname
 Abilities: access-edit
 */
func assignability(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := utils.GetJSON(ctx,"rolename").(string)
    abname := utils.GetJSON(ctx,"abname").(string)

    // Try to get the role
    var role objects.Role
    SpaceDock.Database.Where("name = ?", rolename).First(&role)
    if role.Name != rolename {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The role does not exist. Please add it to a user to create it internally.").Code(3030))
    }

    // Role is valid, assign the new ability
    ability := role.AddAbility(abname)
    SpaceDock.Database.Save(&role).Save(&ability)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path:   /api/access/abilities/remove
 Method: POST
 Description: Removes a permission from a group. Required parameters: rolename, abname
 Abilities: access-edit
 */
func removeability(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := utils.GetJSON(ctx,"rolename").(string)
    abname := utils.GetJSON(ctx,"abname").(string)

    // Try to get the role
    var role objects.Role
    SpaceDock.Database.Where("name = ?", rolename).First(&role)
    var ability objects.Ability
    SpaceDock.Database.Where("name = ?", abname).First(&ability)
    errors := []string{}
    codes := []int{}
    if role.Name != rolename {
        errors = append(errors,"The role does not exist.")
        codes = append(codes, 3030)
    }
    if ability.Name != abname {
        errors = append(errors,"The ability does not exist.")
        codes = append(codes, 2107)
    }
    if len(errors) > 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
    }

    // Both objects are valid, check if they are linked
    if !role.HasAbility(abname) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The ability isn't assigned to this role").Code(1010))
    }

    // Remove the ability
    role.RemoveAbility(abname)
    SpaceDock.Database.Save(&role)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path:   /api/access/params/add/:rolename
 Method: POST
 Description: Adds a parameter for an ability. Required parameters: abname, param, value
 Abilities: access-edit
 */
func addparam(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := ctx.GetString("rolename")
    abname := utils.GetJSON(ctx,"abname").(string)
    param := utils.GetJSON(ctx,"param").(string)
    value := utils.GetJSON(ctx,"value").(string)

    // Try to get the role
    var role objects.Role
    SpaceDock.Database.Where("name = ?", rolename).First(&role)
    var ability objects.Ability
    SpaceDock.Database.Where("name = ?", abname).First(&ability)
    errors := []string{}
    codes := []int{}
    if role.Name != rolename {
        errors = append(errors,"The role does not exist.")
        codes = append(codes, 3030)
    }
    if ability.Name != abname {
        errors = append(errors,"The ability does not exist.")
        codes = append(codes, 2107)
    }
    if len(errors) > 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
    }

    // Both objects are valid, check if they are linked
    role.AddParam(abname, param, value)
    SpaceDock.Database.Save(&role)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path:   /api/access/params/remove/:rolename
 Method: POST
 Description: Removes a parameter from an ability. Required parameters: abname, param, value
 Abilities: access-edit
 */
func removeparam(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := ctx.GetString("rolename")
    abname := utils.GetJSON(ctx,"abname").(string)
    param := utils.GetJSON(ctx,"param").(string)
    value := utils.GetJSON(ctx,"value").(string)

    // Try to get the role
    var role objects.Role
    SpaceDock.Database.Where("name = ?", rolename).First(&role)
    var ability objects.Ability
    SpaceDock.Database.Where("name = ?", abname).First(&ability)
    errors := []string{}
    codes := []int{}
    if role.Name != rolename {
        errors = append(errors,"The role does not exist.")
        codes = append(codes, 3030)
    }
    if ability.Name != abname {
        errors = append(errors,"The ability does not exist.")
        codes = append(codes, 2107)
    }
    if len(errors) > 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
    }

    // Both objects are valid, check if the param exists
    role.RemoveParam(abname, param, value)
    SpaceDock.Database.Save(&role)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}