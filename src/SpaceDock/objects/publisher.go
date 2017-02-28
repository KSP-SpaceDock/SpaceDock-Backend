/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "SpaceDock"
    "SpaceDock/utils"
    "errors"
    "github.com/jinzhu/gorm"
)

type Publisher struct {
    gorm.Model
    MetaObject

    Name             string `gorm:"size:1024;unique_index;not null"`
    Description      string `gorm:"size:100000"`
    ShortDescription string `gorm:"size:1000"`
}

func NewPublisher(name string) *Publisher {
    pub := &Publisher{ Name: name }
    pub.Meta = "{}"
    return pub
}

func (pub Publisher) GetGames() *gorm.DB {
    return SpaceDock.Database.Where("publisherid = ?", pub.ID)
}

func (pub *Publisher) GetById(id interface{}) error {
    SpaceDock.Database.First(pub, id)
    if pub.Name != "" {
        return errors.New("Invalid user ID")
    }
    return nil
}

func (pub Publisher) Format() map[string]interface{} {
    return map[string]interface{} {
        "id": pub.ID,
        "name": pub.Name,
        "description": pub.Description,
        "short_description": pub.ShortDescription,
        "created": pub.CreatedAt,
        "updated": pub.UpdatedAt,
        "meta": utils.LoadJSON(pub.Meta),
    }
}