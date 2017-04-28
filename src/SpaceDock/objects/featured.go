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

type Featured struct {
    Model

    Mod    Mod `json:"mod" spacedock:"lock"`
    ModID  uint `json:"mod_id" spacedock:"lock"`
}

func (s *Featured) AfterFind() {
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
}

func NewFeatured(mod Mod) *Featured {
    f := &Featured{
        Mod: mod,
        ModID: mod.ID,
    }
    f.Meta = "{}"
    return f
}
