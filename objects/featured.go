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
)

type Featured struct {
    Model

    Mod    Mod `json:"mod" spacedock:"lock;tomap"`
    ModID  uint `json:"mod_id" spacedock:"lock"`
}

func (s *Featured) AfterFind() {
    app.DBRecursionLock.Lock()
    if _, ok := app.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        app.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if app.DBRecursion[utils.CurrentGoroutineID()] >= app.DBRecursionMax {
        return
    }
    isRoot := app.DBRecursion[utils.CurrentGoroutineID()] == 0
    app.DBRecursion[utils.CurrentGoroutineID()] += 1
    app.DBRecursionLock.Unlock()

    app.Database.Model(s).Related(&(s.Mod), "Mod")

    app.DBRecursionLock.Lock()
    app.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(app.DBRecursion, utils.CurrentGoroutineID())
    }
    app.DBRecursionLock.Unlock()
}

func NewFeatured(mod Mod) *Featured {
    f := &Featured{
        Mod: mod,
        ModID: mod.ID,
    }
    f.Meta = "{}"
    return f
}
