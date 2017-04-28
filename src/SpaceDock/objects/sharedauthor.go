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
)

type SharedAuthor struct {
    Model

    User   User `json:"-" spacedock:"lock"`
    UserID uint `json:"user" spacedock:"lock"`
    Mod    Mod `json:"-" spacedock:"lock"`
    ModID  uint `json:"mod" spacedock:"lock"`
    Accepted  bool `gorm:"not null" json:"accepted"`
}

func (s *SharedAuthor) AfterFind() {
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.Database.Model(s).Related(&(s.User), "User")
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
}

func NewSharedAuthor(user User, mod Mod) *SharedAuthor {
    author := &SharedAuthor{
        User: user,
        UserID: user.ID,
        Mod: mod,
        ModID: mod.ID,
        Accepted: false,
    }
    author.Meta = "{}"
    return author
}
