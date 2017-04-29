/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "SpaceDock"
    "SpaceDock/utils"
    "errors"
    "github.com/jameskeane/bcrypt"
    "time"
)

type User struct {
    Model

    Username            string `gorm:"size:128;unique_index;not null" json:"username"`
    Email               string `gorm:"size:256;unique_index;not null" json:"email"`
    ShowEmail           bool `json:"showEmail"`
    Public              bool `json:"public"`
    Password            string `gorm:"size:128" json:"-" spacedock:"lock"`
    Description         string `gorm:"size:10000" json:"description"`
    Confirmation        string `gorm:"size:128" json:"-" spacedock:"lock"`
    PasswordReset       string `gorm:"size:128" json:"-" spacedock:"lock"`
    PasswordResetExpiry time.Time `json:"-" spacedock:"lock"`
    Roles               []Role `gorm:"many2many:role_users" json:"-" spacedock:"lock"`
    authed              bool
    SharedAuthors       []SharedAuthor `json:"-" spacedock:"lock"`
    Following           []Mod `json"-" gorm:"many2many:mod_followers" spacedock:"lock"`
}

func (s *User) AfterFind() {
    SpaceDock.DBRecursionLock.Lock()
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] >= SpaceDock.DBRecursionMax {
        SpaceDock.DBRecursionLock.Unlock()
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.DBRecursionLock.Unlock()

    SpaceDock.Database.Model(s).Related(&(s.Roles), "Roles")
    SpaceDock.Database.Model(s).Related(&(s.SharedAuthors), "SharedAuthors")
    SpaceDock.Database.Model(s).Related(&(s.Following), "Following")

    SpaceDock.DBRecursionLock.Lock()
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
    SpaceDock.DBRecursionLock.Unlock()
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
        authed: false,
        Roles: []Role{},
    }
    user.SetPassword(password)
    user.Meta = "{}"
    return user
}

func (user *User) SetPassword(password string) {
    salt, _ := bcrypt.Salt()
    user.Password, _ = bcrypt.Hash(password, salt)
}

/* Login Interface */

func (user *User) IsAuthenticated() bool {
    return user.authed
}

func (user *User) Login() {
    user.authed = true
}

func (user *User) Logout() {
    user.authed = false
}

func (user *User) GetById(id uint) error {
    SpaceDock.Database.Where("id = ?", id).First(user)
    if user.ID != id {
        return errors.New("Invalid user ID")
    }
    return nil
}

/* Login Interface End */

func (user *User) AddRole(name string) *Role {
    role := &Role {}
    SpaceDock.Database.Where("name = ?", name).First(role)
    if role.Name != name {
        role.Name = name
        role.Params = "{}"
        role.Meta = "{}"
        SpaceDock.Database.Save(role)
    }
    SpaceDock.Database.Model(user).Related(&(user.Roles), "Roles")
    user.Roles = append(user.Roles, *role)
    SpaceDock.Database.Save(user).Save(role)
    return &user.Roles[len(user.Roles) - 1]
}

func (user *User) RemoveRole(name string) {
    role := &Role{}
    SpaceDock.Database.Where("name = ?", name).First(role)
    if role.Name == "" {
        return
    }
    if e,i := utils.ArrayContains(role, user.Roles); e {
        user.Roles = append(user.Roles[:i], user.Roles[i + 1:]...)
        SpaceDock.Database.Save(user)
    }
}

func (user *User) HasRole(name string) bool {
    role := &Role {}
    SpaceDock.Database.Where("name = ?", name).First(role)
    if role.Name == "" {
        return false
    }
    e,_ := utils.ArrayContains(role, &(user.Roles))
    return e
}

func (user *User) GetAbilities() []string {
    SpaceDock.Database.Model(user).Related(&(user.Roles), "Roles")
    value := []string{}
    for _,element := range user.Roles {
        for _,element2 := range element.Abilities {
            value = append(value, element2.Name)
        }
    }
    return value
}

func (user *User) Format(admin bool) map[string]interface{} {
    if (admin) {
        roles := user.Roles
        names := make([]string, len(roles))
        for i,element := range roles {
            names[i] = element.Name
        }
        return map[string]interface{}{
            "id": user.ID,
            "created": user.CreatedAt,
            "updated": user.UpdatedAt,
            "username": user.Username,
            "email": user.Email,
            "showEmail": user.ShowEmail,
            "public": user.Public,
            "description": user.Description,
            "roles": names,
            "meta": utils.LoadJSON(user.Meta),
        }
    } else {
        roles := user.Roles
        names := make([]string, len(roles))
        for i,element := range roles {
            names[i] = element.Name
        }
        meta := utils.LoadJSON(user.Meta)
        if _,ok := meta["private"]; ok {
            meta["private"] = map[string]string {}
        }
        return map[string]interface{}{
            "id": user.ID,
            "created": user.CreatedAt,
            "updated": user.UpdatedAt,
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