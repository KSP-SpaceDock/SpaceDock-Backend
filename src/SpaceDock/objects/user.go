/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "github.com/jameskeane/bcrypt"
    "github.com/jinzhu/gorm"
    "time"
    "SpaceDock"
    "errors"
)

type User struct {
    gorm.Model

    Username string `gorm:"size:128;unique_index;not null"`
    Email string `gorm:"size:256;unique_index;not null"`
    Public bool
    Password string `gorm:"size:128"`
    Description string `gorm:"size:10000"`
    Confirmation string `gorm:"size:128"`
    PasswordReset string `gorm:"size:128"`
    PasswordResetExpiry time.Time
    Authed bool
}

func (user User) SetPassword(password string) {
    salt, _ := bcrypt.Salt()
    user.Password, _ = bcrypt.Hash(password, salt)
    SpaceDock.Database.Save(user)
}

func (user User) IsAuthenticated() bool {
    return user.Authed
}

func (user User) Login() {
    user.Authed = true
    SpaceDock.Database.Save(user)
}

func (user User) Logout() {
    user.Authed = false
    SpaceDock.Database.Save(user)
}

func (user User) UniqueId() interface{} {
    return user.ID
}

func (user User) GetById(id interface{}) error {
    SpaceDock.Database.First(&user, id)
    if user.Username != "" {
        return errors.New("Invalid user ID")
    }
    return nil
}