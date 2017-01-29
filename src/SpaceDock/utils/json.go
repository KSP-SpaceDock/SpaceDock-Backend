/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package utils

import "github.com/kataras/iris"

func WriteJSON(ctx *iris.Context, status int, v interface {}) error {
    if _,ok := ctx.URLParams()["callback"]; ok {
        return ctx.JSONP(status, ctx.URLParam("callback"), v)
    }
    return ctx.JSON(status, v)
}
