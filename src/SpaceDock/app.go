/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package SpaceDock

import "log"

/*
 Entrypoint wrapper that is called from the main() function.

 This will load the config, establish a connection to the database
 and startup the webserver to serve JSON
*/
func Run() {
    log.Print("SpaceDock-Backend -- Version: {$VERSION}")
    log.Print("* Loading configuration")
    LoadSettings()

    // Debug
    log.Print(settings.SiteName)

}
