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

type ModList struct {
    Model

    User             User `json:"-" spacedock:"lock"`
    UserID           uint `json:"user" spacedock:"lock"`
    Game             Game `json:"-" spacedock:"lock"`
    GameID           uint `json:"game" spacedock:"lock"`
    Description      string `gorm:"size:100000"`
    ShortDescription string `gorm:"size:1000"`
    Name             string `json:"name" gorm:"size:1024;unique_index;not null"`
    Mods             []ModListItem `json:"-" spacedock:"lock"`
}

func (s *ModList) AfterFind() {
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
    SpaceDock.Database.Model(s).Related(&(s.Game), "Game")
    SpaceDock.Database.Model(s).Related(&(s.Mods), "Mods")

    SpaceDock.DBRecursionLock.Lock()
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
    SpaceDock.DBRecursionLock.Unlock()
}

func NewModList(name string, user User, game Game) *ModList {
    modlist := &ModList{
        User: user,
        UserID: user.ID,
        Game: game,
        GameID: game.ID,
        Name: name,
        Description: "",
        ShortDescription: "",
    }
    modlist.Meta = "{}"
    return modlist
}