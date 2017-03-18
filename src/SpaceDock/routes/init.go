/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "SpaceDock"
    "gopkg.in/kataras/iris.v6"
    "gopkg.in/kataras/iris.v6/middleware/logger"
    "github.com/iris-contrib/middleware/cors"
    "github.com/ulule/limiter"
    "SpaceDock/middleware"
    "log"
)

/*
 Init function, here we register the routes in iris
 */
func init() {
    // Middlewares
    MiddlewareRegister()

    AccessRegister()
    AccountsRegister()
    AdminRegister()
    GameRegister()
    TokensRegister()
    UserRegister()
}

const (
    GET = 1
    POST = 2
    PUT = 3
    DELETE = 4
)

func Register(mode int, path string, handlers ...iris.HandlerFunc) iris.RouteInfo {
    if mode == 2 {
        return SpaceDock.App.Post(path, handlers...)
    } else if mode == 3 {
        return SpaceDock.App.Put(path, handlers...)
    } else if mode == 4 {
        return SpaceDock.App.Delete(path, handlers...)
    }
    return SpaceDock.App.Get(path, handlers...)
}

func MiddlewareRegister() {
    // Logging
    SpaceDock.App.Use(logger.New())

    // Request limiting
    rate,err := limiter.NewRateFromFormatted(SpaceDock.Settings.RequestLimit)
    if err != nil {
        log.Fatal("Failed to parse the request limit")
        return
    }
    store := limiter.NewMemoryStore()
    limiterInstance := limiter.NewLimiter(store, rate)
    SpaceDock.App.Use(middleware.NewAccessLimiter(limiterInstance))

    // Cross-Origin Requests
    if SpaceDock.Settings.DisableSameOrigin {
        SpaceDock.App.Use(cors.New(cors.Options{
            AllowedOrigins:   []string{"*"},
            AllowedHeaders:   []string{"*"},
            AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
            AllowCredentials: true,
        }))
    }
}
