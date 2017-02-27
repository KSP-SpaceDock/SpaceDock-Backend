/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "github.com/jinzhu/gorm"
    "time"
    "SpaceDock/utils"
)

type Game struct {
    gorm.Model
    MetaObject

    Name string `gorm:"size:1024;unique_index;not null"`
    Active bool
    Fileformats string `gorm:"size:1024"`
    Altname string `gorm:"size:1024"`
    Rating float32
    Releasedate time.Time
    Short string `gorm:"size:1024"`
    publisherID uint
    Description string `gorm:"size:100000"`
    ShortDescription string `gorm:"size:1000"`
    // Mods []Mod
    // Modlists []ModList
    // Versions []GameVersion

}

func NewGame(name string, publisher Publisher, short string) *Game {
    game := &Game {
        Name: name,
        Active: false,
        Fileformats: "{\"zip\": \"application/zip\"}",
        Altname: "",
        Rating: 0,
        Releasedate: time.Now(),
        Short: short,
        Description: "",
        ShortDescription: "",
        publisherID: publisher.ID,
    }
    game.Meta = "{}"
    return game
}

func (game Game) GetPublisher() *Publisher {
    pub := &Publisher{}
    err := pub.GetById(game.publisherID)
    if err != nil {
        return nil
    }
    return pub
}

func (game Game) Format() map[string]interface{} {
    return map[string]interface{} {
        "id": game.ID,
        "name": game.Name,
        "active": game.Active,
        "fileformats": utils.LoadJSON(game.Fileformats),
        "rating": game.Rating,
        "releasedate": game.Releasedate,
        "short": game.Short,
        "publisher": game.publisherID,
        "description": game.Description,
        "short_description": game.ShortDescription,
        "created": game.CreatedAt,
        "updated": game.UpdatedAt,
        "meta": utils.LoadJSON(game.Meta),
    }
}