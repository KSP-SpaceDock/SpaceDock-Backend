/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
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

    SpaceDock.Database.Model(s).Related(&(s.User), "User")
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")

    SpaceDock.DBRecursionLock.Lock()
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
    SpaceDock.DBRecursionLock.Unlock()
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
