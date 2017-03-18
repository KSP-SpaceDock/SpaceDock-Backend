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

type GameVersion struct {
    Model

    Game            Game `json:"-" spacedock:"lock"`
    GameID          uint `json:"game" spacedock:"lock"`
    FriendlyVersion string `gorm:"size:128;" json:"friendly_version" spacedock:"lock"`
    Beta            bool
}

func (s *GameVersion) AfterFind() {
    if SpaceDock.DBRecursion == 2 {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.Game), "Game")
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
}

func NewGameVersion(friendly_version string, game Game, beta bool) *GameVersion {
    gv := &GameVersion{
        FriendlyVersion: friendly_version,
        Beta: beta,
        Game: game,
        GameID: game.ID,
    }
    gv.Meta = "{}"
    return gv
}