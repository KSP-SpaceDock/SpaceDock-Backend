/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import "SpaceDock"

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
    if SpaceDock.DBRecursion == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.User), "User")
    SpaceDock.Database.Model(s).Related(&(s.Game), "Game")
    SpaceDock.Database.Model(s).Related(&(s.Mods), "Mods")
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
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