/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package SpaceDock

import (
    "github.com/kataras/iris"
    "log"
    "strconv"
)

/*
 The webserver that will listen for Requests
 */
var app iris.Framework

/*
 Entrypoint wrapper that is called from the main() function.

 This will load the config, establish a connection to the database
 and startup the webserver to serve JSON
 */
func Run() {
    log.Print("SpaceDock-Backend -- Version: {$VERSION}")
    log.Print("* Loading configuration")
    LoadSettings()

    // Connect to the database
    log.Print("* Establishing Database connection")

    // Create the App
    log.Print("* Initializing Iris-Framework")
    app = *iris.New(iris.Configuration{IsDevelopment: settings.Debug })

    // Load routes here

    // Start listening
    app.Listen(settings.Host + ":" + strconv.Itoa(settings.Port))
}
