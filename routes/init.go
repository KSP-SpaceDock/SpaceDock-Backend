/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/middleware"
    "github.com/iris-contrib/middleware/cors"
    "github.com/ulule/limiter"
    "github.com/ulule/limiter/drivers/store/memory"
    "gopkg.in/kataras/iris.v6"
    "gopkg.in/kataras/iris.v6/middleware/logger"
    "log"
    "runtime"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
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
    FeaturedRegister()
    GameRegister()
    GeneralRegister()
    ModlistsRegister()
    ModsRegister()
    PublisherRegister()
    TokensRegister()
    UserRegister()
}

const (
    GET = 1
    POST = 2
    PUT = 3
    DELETE = 4
)

var routes []string

func Register(mode int, path string, handlers ...iris.HandlerFunc) iris.RouteInfo {
    if ok, _ := utils.ArrayContains(path, routes); !ok {
        app.App.Options(path, func(ctx *iris.Context){})
        routes = append(routes, path)
    }
    if mode == 2 {
        return app.App.Post(path, handlers...)
    } else if mode == 3 {
        return app.App.Put(path, handlers...)
    } else if mode == 4 {
        return app.App.Delete(path, handlers...)
    }
    return app.App.Get(path, handlers...)
}

func MiddlewareRegister() {
    // Logging
    app.App.Use(logger.New())

    // Request limiting
    rate,err := limiter.NewRateFromFormatted(app.Settings.RequestLimit)
    if err != nil {
        log.Fatal("Failed to parse the request limit")
        return
    }
    store := memory.NewStore()
    limiterInstance := limiter.New(store, rate)
    app.App.Use(middleware.NewAccessLimiter(limiterInstance))

    // Cross-Origin Requests
    if app.Settings.DisableSameOrigin {
        app.App.Use(cors.New(cors.Options{
            AllowedOrigins:   []string{"*"},
            AllowedHeaders:   []string{"*"},
            AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
            AllowCredentials: true,
        }))
    }

    // Force garbage collection
    app.App.Use(iris.HandlerFunc(func (ctx *iris.Context) {
        runtime.GC()
        ctx.Next()
    }))

    // Cache
    middleware.CreateCache()
}
