/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package SpaceDock

import (
    "github.com/kataras/iris"
    "github.com/kataras/iris/adaptors/httprouter"
    "github.com/kataras/iris/adaptors/sessions"
    "log"
    "strconv"
)

/*
 The webserver that will listen for Requests
 */
var App *iris.Framework

/*
 Startup function for the app

 This will load the config, establish a connection to the database
 and startup the webserver to serve JSON
 */
func init() {
    log.Print("SpaceDock-Backend -- Version: {$VERSION}")
    log.Print("* Loading configuration")
    LoadSettings()

    // Connect to the database
    log.Print("* Establishing Database connection")
    LoadDatabase()

    // Create the App
    log.Print("* Initializing Iris-Framework")
    App = iris.New()
    App.Adapt(httprouter.New())
    App.Adapt(iris.DevLogger())
    mySessions := sessions.New(sessions.Config{
        Cookie: "spacedocksid",
        DecodeCookie: false,
        Expires: 0,
        CookieLength: 32,
        DisableSubdomainPersistence: false,
    })
    App.Adapt(mySessions)
    App.Config.Gzip = true
    App.Config.DisableBodyConsumptionOnUnmarshal = true
}

/*
 Entrypoint wrapper that is called from the main() function
 */
func Run() {
    // Start listening
    App.Listen(Settings.Host + ":" + strconv.Itoa(Settings.Port))
}