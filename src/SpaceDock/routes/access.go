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
 Registers the routes for the account management
 */
func AccessRegister() {
    Register(GET, "/api/access/roles",
        middleware.NeedsPermission("access-view", true),
        list_roles,
    )
    Register(POST, "/api/access/roles",
        middleware.NeedsPermission("access-edit", true),
        assign_role,
    )
    Register(DELETE, "/api/access/roles",
        middleware.NeedsPermission("access-edit", true),
        remove_role,
    )
    Register(GET, "/api/access/abilities",
        middleware.NeedsPermission("access-view", true),
        list_abilities,
    )
    Register(POST, "/api/access/abilities",
        middleware.NeedsPermission("access-edit", true),
        assign_ability,
    )
    Register(DELETE, "/api/access/abilities",
        middleware.NeedsPermission("access-edit", true),
        remove_ability,
    )
    Register(POST, "/api/access/roles/:rolename/params",
        middleware.NeedsPermission("access-edit", true),
        add_param,
    )
    Register(DELETE, "/api/access/roles/:rolename/params",
        middleware.NeedsPermission("access-edit", true),
        remove_param,
    )
}

/*
 Path: /api/access/roles/
 Method: GET
 Description: Displays  list of all roles with the matching abilities
 Abilities: access-view
 */
func list_roles(ctx *iris.Context) {
    roles := []objects.Role{}
    SpaceDock.Database.Find(&roles)
    output := make([]map[string]interface{}, len(roles))
    for i,element := range roles {
        abilities := element.Abilities
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
 Path: /api/access/roles/
 Method: POST
 Description: Promotes a user for the given role. Required parameters: userid, rolename
 Abilities: access-edit
 */
func assign_role(ctx *iris.Context) {
    // Grab parameters from the JSON
    userid := cast.ToUint(utils.GetJSON(ctx,"userid"))
    rolename := cast.ToString(utils.GetJSON(ctx,"rolename"))

    // Try to get the user
    user := &objects.User{}
    SpaceDock.Database.Where("id = ?", userid).First(user)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid.").Code(2145))
    }

    // User is valid, assign the new role
    role := user.AddRole(rolename)
    SpaceDock.Database.Save(role).Save(user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}



/*
 Path: /api/access/roles/
 Method: DELETE
 Description: Promotes a user for the given role. Required parameters: userid, rolename
 Abilities: access-edit
 */
func remove_role(ctx *iris.Context) {
    // Grab parameters from the JSON
    userid := cast.ToUint(utils.GetJSON(ctx,"userid"))
    rolename := cast.ToString(utils.GetJSON(ctx,"rolename"))

    // Try to get the user
    user := &objects.User{}
    SpaceDock.Database.Where("id = ?", userid).First(user)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The userid is invalid.").Code(2145))
        return
    }
    if user.HasRole(rolename) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The user doesn't have this role").Code(1015))
        return
    }

    // Everything is valid, remove the role
    user.RemoveRole(rolename)
    SpaceDock.Database.Save(user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/access/abilities
 Method: GET
 Description: Displays a list of all abilities.
 Abilities: access-view
 */
func list_abilities(ctx *iris.Context) {
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
 Path: /api/access/abilities/
 Method: POST
 Description: Adds a permission to a group. Required parameters: rolename, abname
 Abilities: access-edit
 */
func assign_ability(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := cast.ToString(utils.GetJSON(ctx,"rolename"))
    abname := cast.ToString(utils.GetJSON(ctx,"abname"))

    // Try to get the role
    role := &objects.Role{}
    SpaceDock.Database.Where("name = ?", rolename).First(role)
    if role.Name != rolename {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The role does not exist. Please add it to a user to create it internally.").Code(3030))
        return
    }

    // Role is valid, assign the new ability
    ability := role.AddAbility(abname)
    SpaceDock.Database.Save(role).Save(ability)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/access/abilities/
 Method: DELETE
 Description: Removes a permission from a group. Required parameters: rolename, abname
 Abilities: access-edit
 */
func remove_ability(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := cast.ToString(utils.GetJSON(ctx,"rolename"))
    abname := cast.ToString(utils.GetJSON(ctx,"abname"))

    // Try to get the role
    role := &objects.Role{}
    SpaceDock.Database.Where("name = ?", rolename).First(role)
    ability := &objects.Ability{}
    SpaceDock.Database.Where("name = ?", abname).First(ability)
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
        return
    }

    // Both objects are valid, check if they are linked
    if !role.HasAbility(abname) {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The ability isn't assigned to this role").Code(1010))
        return
    }

    // Remove the ability
    role.RemoveAbility(abname)
    SpaceDock.Database.Save(role)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/access/roles/:rolename/params/
 Method: POST
 Description: Adds a parameter for an ability. Required parameters: abname, param, value
 Abilities: access-edit
 */
func add_param(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := ctx.GetString("rolename")
    abname := cast.ToString(utils.GetJSON(ctx,"abname"))
    param := cast.ToString(utils.GetJSON(ctx,"param"))
    value := cast.ToString(utils.GetJSON(ctx,"value"))

    // Try to get the role
    role := &objects.Role{}
    SpaceDock.Database.Where("name = ?", rolename).First(role)
    ability := &objects.Ability{}
    SpaceDock.Database.Where("name = ?", abname).First(ability)
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
        return
    }

    // Both objects are valid, check if they are linked
    role.AddParam(abname, param, value)
    SpaceDock.Database.Save(role)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/access/roles/:rolename/params/
 Method: DELETE
 Description: Removes a parameter from an ability. Required parameters: abname, param, value
 Abilities: access-edit
 */
func remove_param(ctx *iris.Context) {
    // Grab parameters from the JSON
    rolename := ctx.GetString("rolename")
    abname := cast.ToString(utils.GetJSON(ctx,"abname"))
    param := cast.ToString(utils.GetJSON(ctx,"param"))
    value := cast.ToString(utils.GetJSON(ctx,"value"))

    // Try to get the role
    role := &objects.Role{}
    SpaceDock.Database.Where("name = ?", rolename).First(role)
    ability := &objects.Ability{}
    SpaceDock.Database.Where("name = ?", abname).First(ability)
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
        return
    }

    // Both objects are valid, check if the param exists
    role.RemoveParam(abname, param, value)
    SpaceDock.Database.Save(role)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}