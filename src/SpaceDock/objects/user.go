/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "SpaceDock"
    "SpaceDock/utils"
    "errors"
    "github.com/jameskeane/bcrypt"
    "gopkg.in/kataras/iris.v6"
    "time"
)

type User struct {
    Model

    Username            string `gorm:"size:128;unique_index;not null"`
    Email               string `gorm:"size:256;not null"`
    ShowEmail           bool
    Public              bool
    Password            string `gorm:"size:128"`
    Description         string `gorm:"size:10000"`
    Confirmation        string `gorm:"size:128"`
    PasswordReset       string `gorm:"size:128"`
    PasswordResetExpiry time.Time
    Authed              bool

    roleUsers []RoleUser
}

func NewUser(name string, email string, password string) *User {
    user := &User {
        Username: name,
        Email: email,
        ShowEmail: false,
        Public: false,
        Description: "",
        Confirmation: "",
        PasswordReset: "",
        PasswordResetExpiry: time.Now(),
        Authed: false,
    }
    user.SetPassword(password)
    user.Meta = "{}"
    return user
}

func (user *User) SetPassword(password string) {
    salt, _ := bcrypt.Salt()
    user.Password, _ = bcrypt.Hash(password, salt)
}

func (user User) IsAuthenticated() bool {
    return user.Authed
}

func (user User) Login() {
    user.Authed = true
    SpaceDock.Database.Save(&user)
}

func (user User) Logout() {
    user.Authed = false
    SpaceDock.Database.Save(&user)
}

func (user User) UniqueId() interface{} {
    return user.ID
}

func (user *User) GetById(id interface{}) error {
    SpaceDock.Database.First(&user, id)
    if user.Username != "" {
        return errors.New("Invalid user ID")
    }
    return nil
}

func (user User) AddRole(name string) Role {
    role := Role {}
    SpaceDock.Database.Where("name = ?", name).First(&role)
    if role.Name == "" {
        role.Name = name
        role.Params = "{}"
        role.Meta = "{}"
        SpaceDock.Database.Save(&role)
    }
    ru := RoleUser{}
    SpaceDock.Database.Where("role_id = ?", role.ID).Where("user_id = ?", user.ID).First(&ru)
    if ru.RoleID != role.ID || ru.UserID != user.ID {
        SpaceDock.Database.Save(NewRoleUser(user, role))
    }
    return role
}

func (user User) RemoveRole(name string) {
    role := Role{}
    SpaceDock.Database.Where("name = ?", name).First(&role)
    if role.Name == "" {
        return
    }
    ru := RoleUser{}
    SpaceDock.Database.Where("role_id = ?", role.ID).Where("user_id = ?", user.ID).First(&ru)
    if ru.RoleID == role.ID && ru.UserID == user.ID {
        SpaceDock.Database.Delete(&ru)
    }
}

func (user User) HasRole(name string) bool {
    role := Role {}
    SpaceDock.Database.Where("name = ?", name).First(&role)
    if role.Name == "" {
        return false
    }
    ru := RoleUser{}
    SpaceDock.Database.Where("role_id = ?", role.ID).Where("user_id = ?", user.ID).First(&ru)
    return ru.RoleID == role.ID && ru.UserID == user.ID
}

func (user User) GetRoles() []Role {
    value := make([]Role, len(user.roleUsers))
    for index,element := range user.roleUsers {
        role := Role {}
        SpaceDock.Database.First(&role, element.RoleID)
        value[index] = role
    }
    return value
}

func (user User) GetAbilities() []Ability {
    count := 0
    for _,element := range user.GetRoles() {
        count = count + len(element.GetAbilities())
    }
    value := make([]Ability, count)
    c := 0
    for _,element := range user.GetRoles() {
        for _,element2 := range element.GetAbilities() {
            value[c] = element2
            c = c + 1
        }
    }
    return value
}

func (user User) Format(ctx *iris.Context, admin bool) map[string]interface{} {
    if (admin) {
        roles := user.GetRoles()
        names := make([]string, len(roles))
        for i,element := range roles {
            names[i] = element.Name
        }
        return map[string]interface{}{
            "id": user.ID,
            "username": user.Username,
            "email": user.Email,
            "showEmail": user.ShowEmail,
            "public": user.Public,
            "description": user.Description,
            "roles": names,
            "meta": utils.LoadJSON(user.Meta),
        }
    } else {
        roles := user.GetRoles()
        names := make([]string, len(roles))
        for i,element := range roles {
            names[i] = element.Name
        }
        meta := utils.LoadJSON(user.Meta)
        userID,_ := ctx.Session().GetInt("SessionID")
        if _,ok := meta["private"]; ok && user.ID != uint(userID) {
            meta["private"] = map[string]string {}
        }
        return map[string]interface{}{
            "id": user.ID,
            "username": user.Username,
            "email": utils.Ternary(user.ShowEmail, user.Email, ""),
            "showEmail": user.ShowEmail,
            "public": user.Public,
            "description": user.Description,
            "roles": names,
            "meta": utils.LoadJSON(user.Meta),
        }
    }
}