/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "SpaceDock"
)

type Ability struct {
    Model

    Name  string `gorm:"size:128;unique_index;not null" json:"name" spacedock:"lock"`
    Roles []Role `gorm:"many2many:role_abilities" json:"-" spacedock:"lock"`
}

func (s *Ability) AfterFind() {
    if SpaceDock.DBRecursion == 2 {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.Roles), "Roles")
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
}