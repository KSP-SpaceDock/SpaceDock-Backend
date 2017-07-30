/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "gopkg.in/kataras/iris.v6"
    "mime"
    "path/filepath"
)

/*
 Registers the routes for the general routes
 */
func GeneralRegister() {
    Register(GET, "/content/*path", download_file)
}

/*
 Path: /content/*path
 Method: GET
 Description: Downloads a file from the storage.
 */
func download_file(ctx *iris.Context) {
    // Get the path
    path := ctx.GetString("path")

    // Check for a CDN
    if app.Settings.CdnDomain != "" {
        ctx.Redirect("http://" + app.Settings.CdnDomain + "/" + path, iris.StatusMovedPermanently)
        return
    }

    // Check for X-Sendfile
    if app.Settings.UseXAccel == "nginx" {
        ctx.SetHeader("Content-Type", mime.TypeByExtension(filepath.Ext(filepath.Join(app.Settings.Storage, path))))
        ctx.SetHeader("Content-Disposition", "attachment; filename=" + filepath.Base(path))
        ctx.SetHeader("X-Accel-Redirect", "/internal/" + path)
    } else if app.Settings.UseXAccel == "apache" {
        ctx.SetHeader("Content-Type", mime.TypeByExtension(filepath.Ext(filepath.Join(app.Settings.Storage, path))))
        ctx.SetHeader("Content-Disposition", "attachment; filename=" + filepath.Base(path))
        ctx.SetHeader("X-Sendfile", filepath.Join(app.Settings.Storage, path))
    } else {
        ctx.SendFile(filepath.Join(app.Settings.Storage, path), filepath.Base(path))
    }
}