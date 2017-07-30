/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package middleware

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "gopkg.in/kataras/iris.v6"
)

/*
 Changes the recursion limit of the database for a route
 */
func Recursion(level int) func (ctx *iris.Context) {
    return func (ctx *iris.Context) {
        oldMax := app.DBRecursionMax
        app.DBRecursionMax = level
        ctx.Next()
        app.DBRecursionMax = oldMax
        return
    }
}