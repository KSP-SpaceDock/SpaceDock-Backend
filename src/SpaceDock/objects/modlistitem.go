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

type ModListItem struct {
    Model

    Mod       Mod `json:"-" spacedock:"lock"`
    ModID     uint `json:"mod" spacedock:"lock"`
    ModList   ModList `json:"-" spacedock:"lock"`
    ModListID uint `json:"mod_list" spacedock:"lock"`
    SortIndex uint `json:"sort_index" spacedock:"lock"`
}

func (s *ModListItem) AfterFind() {
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.Database.Model(s).Related(&(s.ModList), "ModList")
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
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