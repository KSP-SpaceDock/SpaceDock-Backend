/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "SpaceDock"
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
    if SpaceDock.DBRecursion == 2 {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.User), "User").Related(&(s.Mod), "Mod")
    if isRoot {
        SpaceDock.DBRecursion = 0
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
