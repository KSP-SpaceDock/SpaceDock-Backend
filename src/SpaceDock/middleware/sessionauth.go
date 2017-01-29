/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package middleware

import (
    "SpaceDock/objects"
    "github.com/kataras/iris"
)

func LoginRequired(ctx *iris.Context) {
    userID, err := ctx.Session().GetInt("SessionID")
    if err == nil {
        var user objects.User
        err = user.GetById(userID)
        if !(err == nil && user.IsAuthenticated()) {
            ctx.SetStatusCode(iris.StatusUnauthorized)
            return
        }
    } else {
        ctx.SetStatusCode(iris.StatusUnauthorized)
        return
    }
    ctx.Next()
}

func LoginUser(ctx *iris.Context, user objects.User) {
    user.Login()
    ctx.Session().Set("SessionID", user.ID)
}

func LogoutUser(ctx *iris.Context, user objects.User) {
    user.Logout()
    ctx.Session().Delete("SessionID")
}

func CurrentUser(ctx *iris.Context) (objects.User, bool) {
    userID, err := ctx.Session().GetInt("SessionID")
    if err == nil {
        var user objects.User
        err = user.GetById(userID)
        if err != nil {
            return user, true
        } else {
            return nil, false
        }
    }
    return nil, false
}
