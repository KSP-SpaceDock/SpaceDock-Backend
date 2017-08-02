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

type Mod struct {
    Model

    User             User `json:"-" spacedock:"lock"`
    UserID           uint `json:"user" spacedock:"lock"`
    Game             Game `json:"-" spacedock:"lock"`
    GameID           uint `json:"game" spacedock:"lock"`
    SharedAuthors    []SharedAuthor `json:"shared_authors" spacedock:"lock"`
    Name             string `json:"name" gorm:"size:1024;unique_index;not null"`
    Description      string `json:"description" gorm:"size:100000"`
    ShortDescription string `json:"short_description" gorm:"size:1000"`
    Approved         bool `json:"approved" spacedock:"lock"`
    Published        bool `json:"published" spacedock:"lock"`
    License          string `json:"license" gorm:"size:512"`
    DefaultVersion   ModVersion `json:"default_version" gorm:"ForeignKey:DefaultVersionID" spacedock:"lock;tomap"`
    DefaultVersionID uint `json:"default_version_id"`
    Versions         []ModVersion `json:"versions" spacedock:"lock"`
    DownloadEvents   []DownloadEvent `json:"-" spacedock:"lock"`
    FollowEvents     []FollowEvent `json:"-" spacedock:"lock"`
    ReferralEvents   []ReferralEvent `json:"-" spacedock:"lock"`
    Followers        []User `json:"-" gorm:"many2many:mod_followers" spacedock:"lock"`
    Ratings          []Rating `json:"-" spacedock:"lock"`
    TotalScore       float64 `json:"total_score" gorm:"not null" spacedock:"lock"`
    DownloadCount    int64 `json:"download_count" spacedock:"lock"`
}

func (s *Mod) AfterFind() {
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
    app.Database.Model(s).Related(&(s.Game), "Game")
    app.Database.Model(s).Related(&(s.DefaultVersion), "DefaultVersion")
    app.Database.Model(s).Related(&(s.Versions), "Versions")
    app.Database.Model(s).Related(&(s.Followers), "Followers")
    app.Database.Model(s).Related(&(s.Ratings), "Ratings")
    app.Database.Model(s).Related(&(s.SharedAuthors), "SharedAuthors")

    app.DBRecursionLock.Lock()
    app.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(app.DBRecursion, utils.CurrentGoroutineID())
    }
    app.DBRecursionLock.Unlock()
}

func (mod *Mod) CalculateScore() {
    score := float64(0)
    for _,element := range mod.Ratings {
        score = score + element.Score
    }
    mod.TotalScore = score / float64(len(mod.Ratings))
}

func NewMod(name string, user User, game Game, license string) *Mod {
    mod := &Mod{
        User: user,
        UserID: user.ID,
        Game: game,
        GameID: game.ID,
        Name: name,
        Description: "",
        ShortDescription: "",
        Approved: true, // because hey
        Published: false,
        License: license,
        DefaultVersionID: 0,
        TotalScore: 0,
    }
    mod.Meta = "{}"
    return mod
}