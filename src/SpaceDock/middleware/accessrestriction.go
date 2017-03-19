/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package middleware

import (
    "SpaceDock"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "encoding/json"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "log"
)

/*
 Checks if a User has a given permission + parameters

 Return Codes:
    0 - everything is OK
    1 - No user logged in
    2 - Userprofile isn't public
    3 - User has no permission to view this site
    4 - Role params are invalid
 */
func UserHasPermission(ctx *iris.Context, permission string, public bool, params []string) int {
    user := CurrentUser(ctx)
    if user == nil {
        return 1
    } else if public && !user.Public {
        return 2
    }
    if params == nil {
        params = []string{}
    }

    ability := objects.Ability {}
    SpaceDock.Database.Where("name = ?", permission).First(&ability)

    user_abilities := user.GetAbilities()
    user_params := map[string]map[string][]string{}
    for _,element := range user.Roles {
        var temp map[string]map[string][]string
        err := json.Unmarshal([]byte(element.Params), &temp)
        if err != nil {
            log.Fatal("Invalid Role Parameters! (Rolename: " + element.Name + ")")
            return 4
        }
        for k, v := range temp {
            user_params[k] = v
        }
    }

    has := false
    if ok,_ := utils.ArrayContains(ability.Name, user_abilities); ok {
        if len(params) > 0 {
            for _,element := range params {
                if utils.ArrayContainsRe(getParam(ability.Name, element, user_params), ctx.GetString(element)) || utils.ArrayContainsRe(getParam(ability.Name, element, user_params), cast.ToString(utils.GetFullJSON(ctx)[element])) {
                    has = true
                }
            }
        } else {
            has = true
        }
        if has {
            return 0
        }
    }
    return 3
}

func NeedsPermission(permission string, public bool, params ...string) func(ctx *iris.Context) {
    var a objects.Ability
    SpaceDock.Database.FirstOrInit(&a, objects.Ability{Name: permission})
    return func(ctx *iris.Context) {
        status := UserHasPermission(ctx, permission, public, params)
        if status == 0 {
            ctx.Next()
            return
        } else if status == 1 {
            utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("You need to be logged in to access this page").Code(1035))
            return
        } else if status == 2 {
            utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("Only users with public profiles may access this page.").Code(1000))
            return
        } else if status == 3 {
            utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("You don't have access to this page. You need to have the abilities: " + permission).Code(1020))
            return
        } else {
            utils.WriteJSON(ctx, iris.StatusInternalServerError, utils.Error("Invalid Role parameter detected. Please contact the server administrator").Code(1010))
            return
        }
    }
}

func getParam(ability string, param string, p map[string]map[string][]string) []string {
    if _, ok := p[ability]; ok {
        if _,ok := p[ability][param]; ok {
            return p[ability][param]
        }
    }
    return nil
}

func IsCurrentUser(ctx *iris.Context, user *objects.User) bool {
    if CurrentUser(ctx) == nil {
        return false
    } else {
        return CurrentUser(ctx).ID == user.ID
    }
}