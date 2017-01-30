/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package routes

import (
    "github.com/kataras/iris"
    "SpaceDock"
    "SpaceDock/utils"
    "regexp"
    "SpaceDock/objects"
    "SpaceDock/middleware"
)

/*
 Registers the routes for the account management
 */
func AccountsRegister() {
    Register(POST, "/api/register", register)

    Register(GET, "/api/", middleware.LoginRequired, func (ctx *iris.Context) {
        utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"hi": true})
    })
}

/*
 Path:   /api/register
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
    //followMod := utils.GetJSON(ctx,"follow-mod")
    email := utils.GetJSON(ctx,"email").(string)
    username := utils.GetJSON(ctx,"username").(string)
    password := utils.GetJSON(ctx,"password").(string)
    confirmPassword := utils.GetJSON(ctx,"repeatPassword").(string)
    //data := utils.GetJSON(ctx,"userdata")
    check := utils.GetJSON(ctx,"check")

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
        utils.WriteJSON(ctx, iris.StatusOK, utils.Error(errors...).Code(codes...))
        return
    }

    // Everything is valid, make them an account
    user := objects.NewUser(username, email, password)
    user.Confirmation,_ = utils.RandomHex(20)

    // Eval userdata

    SpaceDock.Database.Save(&user)

    // Send confirmation mail

    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": user})
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

