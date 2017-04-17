/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package middleware

import (
    "SpaceDock"
    "gopkg.in/kataras/iris.v6"
)

/*
 Changes the recursion limit of the database for a route
 */
func Recursion(level int) func (ctx *iris.Context) {
    return func (ctx *iris.Context) {
        oldMax := SpaceDock.DBRecursionMax
        SpaceDock.DBRecursionMax = level
        ctx.Next()
        SpaceDock.DBRecursionMax = oldMax
        return
    }
}