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
    "mime"
    "path/filepath"
)

/*
 Registers the routes for the general routes
 */
func GeneralRegister() {
    Register(GET, "/content/*path", download)
}

/*
 Path: /content/*path
 Method: GET
 Description: Downloads a file from the storage.
 */
func download(ctx *iris.Context) {
    // Get the path
    path := ctx.GetString("path")

    // Check for a CDN
    if SpaceDock.Settings.CdnDomain != "" {
        ctx.Redirect("http://" + SpaceDock.Settings.CdnDomain + "/" + path, iris.StatusMovedPermanently)
        return
    }

    // Check for X-Sendfile
    if SpaceDock.Settings.UseXAccel == "nginx" {
        ctx.SetHeader("Content-Type", mime.TypeByExtension(filepath.Ext(filepath.Join(SpaceDock.Settings.Storage, path))))
        ctx.SetHeader("Content-Disposition", "attachment; filename=" + filepath.Base(path))
        ctx.SetHeader("X-Accel-Redirect", "/internal/" + path)
    } else if SpaceDock.Settings.UseXAccel == "apache" {
        ctx.SetHeader("Content-Type", mime.TypeByExtension(filepath.Ext(filepath.Join(SpaceDock.Settings.Storage, path))))
        ctx.SetHeader("Content-Disposition", "attachment; filename=" + filepath.Base(path))
        ctx.SetHeader("X-Sendfile", filepath.Join(SpaceDock.Settings.Storage, path))
    } else {
        ctx.SendFile(filepath.Join(SpaceDock.Settings.Storage, path), filepath.Base(path))
    }
}