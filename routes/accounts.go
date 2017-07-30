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
    "github.com/KSP-SpaceDock/SpaceDock-Backend/objects"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "github.com/jameskeane/bcrypt"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "strconv"
    "time"
)

/*
 Registers the routes for the account management (aka. stuff that doesn't fit anywhere)
 */
func AccountsRegister() {
    Register(GET, "/api/confirm/:confirmation", confirm) // Maybe switch to POST too?
    Register(POST, "/api/login", login)
    Register(POST, "/api/logout", logout)
    Register(POST, "/api/reset", reset)
    Register(POST, "/api/reset/:username/:confirmation", reset_confirm)
}

/*
 Path: /api/confirm/:confirmation
 Method: GET
 Description: Confirms a newly created useraccount using a random text sequence
 */
func confirm(ctx *iris.Context) {
    // Grab the confirmation sequence
    confirmation := ctx.GetString("confirmation")

    // Try to get a valid user account
    user := &objects.User{}
    app.Database.Where("confirmation = ?", confirmation).First(user)
    if user.Confirmation != confirmation {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("User does not exist or it is already confirmed. Did you mistype the confirmation?").Code(2165))
        return
    }

    // Everything is valid
    user.Confirmation = ""
    middleware.LoginUser(ctx, user)
    role := user.AddRole(user.Username)
    role.AddAbility("user-edit")
    role.AddAbility("mods-add")
    role.AddAbility("lists-add")
    role.AddAbility("logged-in")
    role.AddParam("user-edit", "userid", strconv.Itoa(int(user.ID)))
    role.AddParam("mods-add", "gameshort", ".*")
    role.AddParam("packs-add", "gameshort", ".*")
    app.Database.Save(role)
    app.Database.Save(user)

    // Follow Mod

    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/login
 Method: POST
 */
func login(ctx *iris.Context) {
    // Grab information
    username := cast.ToString(utils.GetJSON(ctx, "username"))
    password := cast.ToString(utils.GetJSON(ctx, "password"))

    // Check if the values are valid
    if username == "" || password == "" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Missing username or password").Code(2515))
        return
    }
    if middleware.CurrentUser(ctx) != nil {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You are already logged in").Code(3060))
        return
    }
    user := &objects.User{}
    app.Database.Where("username = ?", username).First(user)
    if user.Username != username {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Username or password is incorrect").Code(2175))
        return
    }
    if s,_ := bcrypt.Hash(password, user.Password); s != user.Password {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Username or password is incorrect").Code(2175))
        return
    }
    if user.Confirmation != "" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("User is not confirmed").Code(3055))
        return
    }
    middleware.LoginUser(ctx, user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": user.Format(true)})
}

/*
 Path: /api/logout
 Method: GET
 */
func logout(ctx *iris.Context) {
    // Check if a user is logged in
    user := middleware.CurrentUser(ctx)
    if user == nil {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You are not logged in. Logging out now would be a bit difficult, right?").Code(3070))
        return
    }
    middleware.LogoutUser(ctx)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/reset
 Method: POST
 */
func reset(ctx *iris.Context) {
    // Get the email
    email := cast.ToString(utils.GetJSON(ctx, "email"))

    // Check if the values are valid
    if email == "" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("No email address").Code(2520))
        return
    }
    user := &objects.User{}
    app.Database.Where("email = ?", email).First(user)
    if user.Email != email {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("No user for provided email address").Code(2115))
        return
    }
    user.PasswordReset,_ = utils.RandomHex(20)
    user.PasswordResetExpiry = time.Now().Add(time.Hour * 24)
    app.Database.Save(user)
    utils.SendReset(user.Username, user.PasswordReset, user.Email)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/reset/:username/:confirmation
 Method: POST
 */
func reset_confirm(ctx *iris.Context) {
    username := ctx.GetString("username")
    confirmation := ctx.GetString("confirmation")
    user := &objects.User{}
    app.Database.Where("username = ?", username).First(user)
    if user.Username != username {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Username is incorrect").Code(2170))
        return
    }
    if user.PasswordReset == "" || user.PasswordResetExpiry.Before(time.Now()) || user.PasswordReset != confirmation {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Password reset invalid").Code(3000))
        return
    }
    password := utils.GetJSON(ctx, "password").(string)
    password2 := utils.GetJSON(ctx, "password2").(string)
    if password == "" || password2 == "" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Passwords not provided").Code(2525))
        return
    }
    if password != password2 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Passwords do not match").Code(3005))
        return
    }
    user.SetPassword(password)
    user.PasswordReset = ""
    user.PasswordResetExpiry = time.Now()
    app.Database.Save(user)
    if middleware.CurrentUser(ctx) != nil {
        middleware.LogoutUser(ctx)
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}
