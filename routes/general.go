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
    "os"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "io"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/objects"
    "github.com/spf13/cast"
    "archive/zip"
)

/*
 Registers the routes for the general routes
 */
func GeneralRegister() {
    Register(GET, "/content/*path", download_file)
    Register(POST, "/upload/:token", upload_file)
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

/*
 Path: /upload/:token
 Method: POST
 Description: Uploads a file to the storage. Required parameters: file
 */
func upload_file(ctx *iris.Context) {
    // Get parameters
    file,_,err := ctx.FormFile("file")
    tokenID := ctx.GetString("token")

    // Get the token
    token := &objects.Token{}
    app.Database.Where("token = ?", tokenID).First(token)
    if token.Token != tokenID  {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The token ID is invalid").Code(2131))
        return
    }
    _,isUploading := token.GetValue("isUploading")
    _,mustBeZip := token.GetValue("mustBeZip")
    if !cast.ToBool(isUploading) {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The token ID is invalid").Code(2131))
        return
    }

    // Paths
    _,base_path := token.GetValue("path")
    path := filepath.Join(app.Settings.Storage, cast.ToString(base_path))
    full_path := filepath.Dir(path)

    // Remove the old file. If it fails, dont care
    os.MkdirAll(full_path, os.ModePerm)
    _ = os.Remove(path)
    out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        utils.WriteJSON(ctx, iris.StatusInternalServerError, utils.Error(err.Error()).Code(2153))
        return
    }
    io.Copy(out, file)
    out.Close()

    // Check if the file is a zipfile
    if cast.ToBool(mustBeZip) {
        temp, err := zip.OpenReader(path)
        if err != nil {
            _ = os.Remove(filepath.Join(app.Settings.Storage, path))
            utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("This is not a valid zip file.").Code(2160))
            return
        } else {
            temp.Close()
        }
    }

    // Clear the token
    app.Database.Delete(token)

    // Success
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}