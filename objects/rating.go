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

    User   User `json:"-" gorm:"ForeignKey:UserID" spacedock:"lock"`
    UserID uint `json:"user" spacedock:"lock"`
    Mod    Mod `json:"-" gorm:"ForeignKey:ModID" spacedock:"lock"`
    ModID  uint `json:"mod" spacedock:"lock"`
    Score  float64 `gorm:"not null" json:"score"`
}

func (s *Rating) AfterFind() {
    app.DBRecursionLock.Lock()
    if _, ok := app.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        app.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if app.DBRecursion[utils.CurrentGoroutineID()] >= app.DBRecursionMax {
        app.DBRecursionLock.Unlock()
        return
    }
    isRoot := app.DBRecursion[utils.CurrentGoroutineID()] == 0
    app.DBRecursion[utils.CurrentGoroutineID()] += 1
    app.DBRecursionLock.Unlock()

    app.Database.Model(s).Related(&(s.User), "User")
    app.Database.Model(s).Related(&(s.Mod), "Mod")

    app.DBRecursionLock.Lock()
    app.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(app.DBRecursion, utils.CurrentGoroutineID())
    }
    app.DBRecursionLock.Unlock()
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
