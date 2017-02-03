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
    "strconv"
    "github.com/jameskeane/bcrypt"
)

/*
 Registers the routes for the account management
 */
func AccountsRegister() {
    Register(POST, "/api/register", register)
    Register(GET, "/api/confirm/:confirmation", confirm)
    Register(POST, "/api/login", login)
    Register(GET, "/api/logout", logout)

    Register(GET, "/api/", func (ctx *iris.Context) {
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
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
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

/*
 Path:   /api/confirm/:confirmation
 Method: GET
 Description: Confirms a newly created useraccount using a random text sequence
 */
func confirm(ctx *iris.Context) {
    // Grab the confirmation sequence
    confirmation := ctx.GetString("confirmation")

    // Try to get a valid user account
    var user objects.User
    SpaceDock.Database.Where("confirmation = ?", confirmation).First(&user)
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
    role.AddAbility("packs-add")
    role.AddAbility("logged-in")
    role.AddParam("user-edit", "userid", strconv.Itoa(int(user.ID)))
    role.AddParam("mods-add", "gameshort", ".*")
    role.AddParam("packs-add", "gameshort", ".*")
    SpaceDock.Database.Save(&role)
    SpaceDock.Database.Save(&user)

    // Follow Mod

    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}

/*
 Path: /api/login
 Method: POST
 */
func login(ctx *iris.Context) {
    // Grab information
    username := utils.GetJSON(ctx, "username").(string)
    password := utils.GetJSON(ctx, "password").(string)

    // Check if the values are valid
    if username == "" || password == "" {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("Missing username or password").Code(2515))
        return
    }
    if middleware.CurrentUser(ctx) != nil {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You are already logged in").Code(3060))
        return
    }
    var user objects.User
    SpaceDock.Database.Where("username = ?", username).First(&user)
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
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": user})
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
    middleware.LogoutUser(ctx, *user)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}
