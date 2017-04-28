/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "SpaceDock"
    "time"
    "SpaceDock/utils"
)

type Game struct {
    Model

    Name             string `gorm:"size:1024;unique_index;not null" json:"name"`
    Active           bool `json:"active"`
    Altname          string `gorm:"size:1024" json:"altname"`
    Rating           float32 `json:"rating" spacedock:"lock"`
    Releasedate      time.Time `json:"releasedate" spacedock:"lock"`
    Short            string `gorm:"size:1024" json:"short" spacedock:"lock"`
    Publisher        Publisher `json:"-" spacedock:"lock"`
    PublisherID      uint `json:"publisher" spacedock:"lock"`
    Description      string `gorm:"size:100000" json:"description"`
    ShortDescription string `gorm:"size:1000" json:"short_description"`
    Mods             []Mod `json:"-" spacedock:"lock"`
    Modlists         []ModList `json:"-" spacedock:"lock"`
    Versions         []GameVersion `json:"-" spacedock:"lock"`
}

func (s *Game) AfterFind() {
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.Database.Model(s).Related(&(s.Publisher))
    SpaceDock.Database.Model(s).Related(&(s.Mods))
    SpaceDock.Database.Model(s).Related(&(s.Modlists))
    SpaceDock.Database.Model(s).Related(&(s.Versions))
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
}

func NewGame(name string, publisher Publisher, short string) *Game {
    game := &Game {
        Name: name,
        Active: false,
        Altname: "",
        Rating: 0,
        Releasedate: time.Now(),
        Short: short,
        Description: "",
        ShortDescription: "",
        PublisherID: publisher.ID,
    }
    game.Meta = "{}"
    return game
}