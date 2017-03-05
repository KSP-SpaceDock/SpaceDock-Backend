/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package routes

import (
    "SpaceDock"
    "gopkg.in/kataras/iris.v6"
)

/*
 Init function, here we register the routes in iris
 */
func init() {
    AccessRegister()
    AccountsRegister()
    AdminRegister()
    GameRegister()
}

const (
    GET = 1
    POST = 2
)

func Register(mode int, path string, handlers ...iris.HandlerFunc) iris.RouteInfo {
    if mode == 1 {
        return SpaceDock.App.Get(path, handlers...)
    } else {
        return SpaceDock.App.Post(path, handlers...)
    }
}
