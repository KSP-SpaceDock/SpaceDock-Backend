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

func GetJSON(ctx *iris.Context) map[string]interface{} {
    var value map[string]interface{}
    err := ctx.ReadJSON(value)
    if err != nil {
        log.Fatal("Tried to parse invalid JSON")
        return nil
    }
    return value
}
