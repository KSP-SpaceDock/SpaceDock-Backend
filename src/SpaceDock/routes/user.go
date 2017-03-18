/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "SpaceDock"
    "SpaceDock/middleware"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "github.com/kennygrant/sanitize"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "io"
    "os"
    "path/filepath"
    "regexp"
    "time"
)

/*
 Registers the routes for the user management
 */
func UserRegister() {
    Register(GET, "/api/users", list_users)
    Register(POST, "/api/users", register)
    Register(GET, "/api/users/:userid", show_user)
    Register(PUT, "/api/users/:userid",
        middleware.NeedsPermission("user-edit", false, "userid"),
        edit_user,
    )
    Register(POST, "/api/users/:userid/update-media",
        middleware.NeedsPermission("user-edit", false, "userid"),
        update_user_media,
    )
}

/*
 Path: /api/users/
 Method: GET
 Description: Returns a list of users.
 */
func list_users(ctx *iris.Context) {
    var users []objects.User
    SpaceDock.Database.Find(&users)
    output := make([]map[string]interface{}, len(users))
    for i,element := range users {
        userid := uint(1)
        if middleware.CurrentUser(ctx) != nil {
            userid = middleware.CurrentUser(ctx).ID
        }
        if element.ID == userid || middleware.UserHasPermission(ctx, "view-users-full", false, []string{}) == 0 {
            output[i] = element.Format(true)
        } else if element.Public {
            output[i] = element.Format(false)
        }
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(users), "data": output})
}

/*
 Path: /api/users/
 Method: POST
 Description: Creates a new useraccount
 */
func register(ctx *iris.Context) {
    // Check if registration is allowed
    if !SpaceDock.Settings.Registration {
        utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("Registration is disabled").Code(3010))
        return
    }

    // Grab parameters from the JSON
    followMod := cast.ToString(utils.GetJSON(ctx,"follow-mod"))
    email := cast.ToString(utils.GetJSON(ctx,"email"))
    username := cast.ToString(utils.GetJSON(ctx,"username"))
    password := cast.ToString(utils.GetJSON(ctx,"password"))
    confirmPassword := cast.ToString(utils.GetJSON(ctx,"repeatPassword"))
    data := cast.ToStringMap(utils.GetJSON(ctx,"userdata"))
    check := cast.ToString(utils.GetJSON(ctx,"check"))

    var errors []string
    var codes []int
    emailError := checkEmailForRegistration(email)
    if emailError != "" {
        errors = append(errors, emailError)
        codes = append(codes, 4000)
        if check == "email" {
            utils.WriteJSON(ctx, iris.StatusOK, utils.Error(emailError).Code(4000))
            return
        }
    }

    usernameError := checkUsernameForRegistration(username)
    if usernameError != "" {
        errors = append(errors, usernameError)
        codes = append(codes, 4000)
        if check == "username" {
            utils.WriteJSON(ctx, iris.StatusOK, utils.Error(usernameError).Code(4000))
            return
        }
    }

    if password == "" {
        errors = append(errors, "Password is required")
        codes = append(codes, 2515)
        if check == "password" {
            utils.WriteJSON(ctx, iris.StatusOK, utils.Error("Password is required").Code(2515))
            return
        }
    } else {
        if password != confirmPassword {
            errors = append(errors, "Passwords do not match")
            codes = append(codes, 3005)
            if check == "password" {
                utils.WriteJSON(ctx, iris.StatusOK, utils.Error("Passwords do not match").Code(3005))
                return
            }
        }
        if len(password) < 5 {
            errors = append(errors, "Your password must be greater than 5 characters")
            codes = append(codes, 2101)
            if check == "password" {
                utils.WriteJSON(ctx, iris.StatusOK, utils.Error("Your password must be greater than 5 characters").Code(2101))
                return
            }
        }
        if len(password) > 256 {
            errors = append(errors, "We admire your dedication to security, but please use a shorter password")
            codes = append(codes, 2102)
            if check == "password" {
                utils.WriteJSON(ctx, iris.StatusOK, utils.Error("We admire your dedication to security, but please use a shorter password").Code(2102))
                return
            }
        }
    }

    if check == "email" || check == "password" || check == "username" {
        utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
        return
    }
    if len(errors) > 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
        return
    }

    // Everything is valid, make them an account
    user := objects.NewUser(username, email, password)
    user.Confirmation,_ = utils.RandomHex(20)

    // Edit user
    if data != nil {
        code := utils.EditObject(&user, data)
        if code == 3 {
            utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The value you submitted is invalid").Code(2180))
            return
        } else if code == 2 {
            utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that doesn't exist.").Code(3090))
            return
        } else if code == 1 {
            utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that is marked as read-only.").Code(3095))
            return
        }
    }

    SpaceDock.Database.Save(user)
    utils.SendConfirmation(user.Confirmation, user.Username, user.Email, followMod)

    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": user.Format(true)})
}

func checkUsernameForRegistration(username string) string {
    if username == "" {
        return "Username is required"
    }
    r,_ := regexp.Compile("^[A-Za-z0-9_]+$")
    var user objects.User
    if !r.MatchString(username) {
        return "Please only use letters, numbers, and underscores"
    }
    if len(username) < 3 || len(username) > 24 {
        return "Usernames must be between 3 and 24 characters"
    }
    if SpaceDock.Database.Where("username = ?", username).First(&user); user.Username != "" {
        return "A user by this name already exists"
    }
    return ""

}

func checkEmailForRegistration(email string) string {
    if email == "" {
        return "Email is required"
    }
    r,_ := regexp.Compile("^[^@]+@[^@]+.[^@]+$")
    var user objects.User
    if !r.MatchString(email) {
        return "Please specify a valid email address."
    } else if SpaceDock.Database.Where("email = ?", email).First(&user); user.Username != "" {
        return "A user with this email already exists."
    }
    return ""
}

/*
 Path: /api/users/:userid
 Method: GET
 Description: Returns more data for one user
 */
func show_user(ctx *iris.Context) {
    userid_ := ctx.GetString("userid")
    userid := uint(0)
    var user objects.User
    if userid_ == "current" {
        if middleware.CurrentUser(ctx) == nil {
            utils.WriteJSON(ctx, iris.StatusForbidden, utils.Error("You need to be logged in to access this page").Code(1035))
            return
        }
        user = *middleware.CurrentUser(ctx)
        userid = user.ID
    } else {
        userid = cast.ToUint(userid_)
        SpaceDock.Database.Where("id = ?", userid).First(&user)
    }
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The userid is invalid.").Code(2145))
        return
    }
    output := map[string]interface{} {}
    if middleware.IsCurrentUser(ctx, &user) || middleware.UserHasPermission(ctx, "view-users-full", false, []string{}) == 0 {
        output = user.Format(true)
    } else if user.Public {
        output = user.Format(false)
    } else {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The userid is invalid.").Code(2145))
        return
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": output})
}

/*
 Path: /api/users/:userid
 Method: PUT
 Description: Edits a user, based on the request parameters. Required fields: data
 Abilities: user-edit
 */
func edit_user(ctx *iris.Context) {
    userid := cast.ToUint(ctx.GetString("userid"))
    var user objects.User
    SpaceDock.Database.Where("id = ?", userid).First(&user)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The userid is invalid.").Code(2145))
        return
    }

    // Everything is ok, edit the user
    code := utils.EditObject(&user, utils.GetFullJSON(ctx))
    if code == 3 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The value you submitted is invalid").Code(2180))
        return
    } else if code == 2 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that doesn't exist.").Code(3090))
        return
    } else if code == 1 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that is marked as read-only.").Code(3095))
        return
    } else {
        utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": user.Format(true)})
        return
    }
}

/*
 Path: /api/users/:userid/update-media
 Method: POST
 Description: Updates a users background. Required fields: image, type
 Abilities: user-edit
 */
func update_user_media(ctx *iris.Context) {
    mediatype := cast.ToString(utils.GetJSON(ctx, "type"))
    userid := cast.ToUint(ctx.GetString("userid"))
    data, info, err := ctx.FormFile("media")
    if err != nil {
        utils.WriteJSON(ctx, iris.StatusInternalServerError, utils.Error(err.Error()).Code(2153))
        return
    }
    var user objects.User
    SpaceDock.Database.Where("id = ?", userid).First(&user)
    if user.ID != userid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The userid is invalid.").Code(2145))
        return
    }

    // Get the file and save it to disk
    ext := filepath.Ext(filepath.Base(info.Filename))
    if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("This file type is not acceptable.").Code(3035))
        return
    }
    filename := sanitize.BaseName(user.Username) + "_" + mediatype + ext
    base_path := filepath.Join(sanitize.BaseName(user.Username) + "-" + time.Now().String() + "_" + cast.ToString(user.ID))
    full_path := filepath.Join(SpaceDock.Settings.Storage, base_path)
    os.MkdirAll(full_path, os.ModePerm)
    path := filepath.Join(full_path, filename)

    // Remove the old file. If it fails, dont care
    err, val := user.GetValue(mediatype)
    if err == nil {
        _ = os.Remove(filepath.Join(SpaceDock.Settings.Storage, cast.ToString(val)))
    }
    out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        utils.WriteJSON(ctx, iris.StatusInternalServerError, utils.Error(err.Error()).Code(2153))
        return
    }
    io.Copy(out, data)
    out.Close()
    data.Close()

    // Edit the user object
    code := utils.EditObject(&user, iris.Map{"meta": iris.Map{mediatype:filepath.Join(base_path, filename)}})
    if code == 3 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The value you submitted is invalid").Code(2180))
        return
    } else if code == 2 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that doesn't exist.").Code(3090))
        return
    } else if code == 1 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that is marked as read-only.").Code(3095))
        return
    } else {
        utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": user.Format(true)})
        return
    }
}