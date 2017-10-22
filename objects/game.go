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
    "time"
)

type Game struct {
    Model

    Name             string `gorm:"size:1024;unique_index;not null" json:"name"`
    Active           bool `json:"active"`
    Altname          string `gorm:"size:1024" json:"altname"`
    Rating           float32 `json:"rating" spacedock:"lock"`
    Releasedate      time.Time `json:"releasedate" spacedock:"lock"`
    Short            string `gorm:"size:1024" json:"short" spacedock:"lock"`
    Publisher        Publisher `json:"-" gorm:"ForeignKey:PublisherID" spacedock:"lock"`
    PublisherID      uint `json:"publisher" spacedock:"lock"`
    Description      string `gorm:"size:100000" json:"description"`
    ShortDescription string `gorm:"size:1000" json:"short_description"`
    Mods             []Mod `json:"-" spacedock:"lock"`
    Modlists         []ModList `json:"-" spacedock:"lock"`
    Versions         []GameVersion `json:"-" spacedock:"lock"`
}

func (s *Game) AfterFind() {
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

    app.Database.Model(s).Related(&(s.Publisher))
    app.Database.Model(s).Related(&(s.Mods))
    app.Database.Model(s).Related(&(s.Modlists))
    app.Database.Model(s).Related(&(s.Versions))

    app.DBRecursionLock.Lock()
    app.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(app.DBRecursion, utils.CurrentGoroutineID())
    }
    app.DBRecursionLock.Unlock()
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