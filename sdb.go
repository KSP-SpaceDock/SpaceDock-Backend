/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package main

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    _ "github.com/KSP-SpaceDock/SpaceDock-Backend/routes"
)

/*
 The entrypoint for the spacedock application.
 Instead of running significant code here, we pass this task to the app package
*/
func main() {
    app.Run()
}

