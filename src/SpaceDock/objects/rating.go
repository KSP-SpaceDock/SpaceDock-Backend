/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "SpaceDock"
    "math"
)

type Rating struct {
    Model

    User   User `json:"-" spacedock:"lock"`
    UserID uint `json:"user" spacedock:"lock"`
    Mod    Mod `json:"-" spacedock:"lock"`
    ModID  uint `json:"mod" spacedock:"lock"`
    Score  float64 `gorm:"not null" json:"score"`
}

func (s *Rating) AfterFind() {
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

func NewRating(user User, mod Mod, score float64) *Rating {
    rating := &Rating{
        User: user,
        UserID: user.ID,
        Mod: mod,
        ModID: mod.ID,
        Score: math.Max(0, math.Min(score, 5)),
    }
    rating.Meta = "{}"
    return rating
}
