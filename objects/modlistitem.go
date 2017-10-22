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

type ModListItem struct {
    Model

    Mod       Mod `json:"-" gorm:"ForeignKey:ModID" spacedock:"lock"`
    ModID     uint `json:"mod" spacedock:"lock"`
    ModList   ModList `json:"-" gorm:"ForeignKey:ModListID" spacedock:"lock"`
    ModListID uint `json:"mod_list" spacedock:"lock"`
    SortIndex uint `json:"sort_index" spacedock:"lock"`
}

func (s *ModListItem) AfterFind() {
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

    app.Database.Model(s).Related(&(s.Mod), "Mod")
    app.Database.Model(s).Related(&(s.ModList), "ModList")

    app.DBRecursionLock.Lock()
    app.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(app.DBRecursion, utils.CurrentGoroutineID())
    }
    app.DBRecursionLock.Unlock()
}

func NewModListItem(mod Mod, list ModList) *ModListItem {
    modlistitem := &ModListItem{
        Mod: mod,
        ModID: mod.ID,
        ModList: list,
        ModListID: list.ID,
        SortIndex:0,
    }
    modlistitem.Meta = "{}"
    return modlistitem
}