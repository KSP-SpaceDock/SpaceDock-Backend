/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package utils

import (
    "github.com/kataras/iris"
    "log"
)

func GetFullJSON(ctx *iris.Context) map[string]interface{} {
    var value map[string]interface{}
    err := ctx.ReadJSON(value)
    if err != nil {
        log.Fatal("Tried to parse invalid JSON")
        return nil
    }
    return value
}

func GetJSON(ctx *iris.Context, key string) interface{} {
    full := GetFullJSON(ctx)
    if value,ok := full[key]; ok {
        return value
    }
    return nil
}

func WriteJSON(ctx *iris.Context, status int, v interface {}) error {
    if _,ok := ctx.URLParams()["callback"]; ok {
        return ctx.JSONP(status, ctx.URLParam("callback"), v)
    }
    return ctx.JSON(status, v)
}

