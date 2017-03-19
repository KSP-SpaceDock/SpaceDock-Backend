/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import "SpaceDock"

type Mod struct {
    Model

    User             User `json:"-" spacedock:"lock"`
    UserID           uint `json:"user" spacedock:"lock"`
    Game             Game `json:"-" spacedock:"lock"`
    GameID           uint `json:"game" spacedock:"lock"`
    SharedAuthors    []SharedAuthor `json:"-" spacedock:"lock"`
    Name             string `json:"name" gorm:"size:1024;unique_index;not null"`
    Description      string `json:"description" gorm:"size:100000"`
    ShortDescription string `json:"short_description" gorm:"size:1000"`
    Approved         bool `json:"approved" spacedock:"lock"`
    Published        bool `json:"published" spacedock:"lock"`
    License          string `json:"license" gorm:"size:512"`
    DefaultVersion   *ModVersion `json:"-" spacedock:"lock"`
    DefaultVersionID uint `json:"default_version"`
    Versions         []ModVersion `json:"-" spacedock:"lock"`
    // Todo: Tracking API
    Followers  []User `json"-" gorm:"many2many:mod_followers" spacedock:"lock"`
    Ratings    []Rating `json:"-" spacedock:"lock"`
    TotalScore float64 `json:"total_score" gorm:"not null" json:"total_score" spacedock:"lock"`
}

func (s *Mod) AfterFind() {
    if SpaceDock.DBRecursion == 2 {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.User), "User")
    SpaceDock.Database.Model(s).Related(&(s.Game), "Game")
    SpaceDock.Database.Model(s).Related(s.DefaultVersion, "DefaultVersion")
    SpaceDock.Database.Model(s).Related(&(s.Versions), "Versions")
    SpaceDock.Database.Model(s).Related(&(s.Followers), "Followers")
    SpaceDock.Database.Model(s).Related(&(s.Ratings), "Ratings")
    SpaceDock.Database.Model(s).Related(&(s.SharedAuthors), "SharedAuthors")
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
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