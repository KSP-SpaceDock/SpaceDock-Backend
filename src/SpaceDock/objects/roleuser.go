/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "github.com/jinzhu/gorm"
)

type RoleUser struct {
    gorm.Model

    UserID uint
    RoleID uint
}

func NewRoleUser(user User, role Role) *RoleUser {
    return &RoleUser{ UserID: user.ID, RoleID: role.ID }
}

func (ru RoleUser) GetUser() User {
    user := User {}
    err := user.GetById(ru.UserID)
    if err != nil {
        return User{}
    }
    return user
}

func (ru RoleUser) GetRole() Role {
    role := Role {}
    err := role.GetById(ru.RoleID)
    if err != nil {
        return Role{}
    }
    return role
}
