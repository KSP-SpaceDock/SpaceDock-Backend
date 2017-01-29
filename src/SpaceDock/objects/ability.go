/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "SpaceDock"
    "errors"
    "github.com/jinzhu/gorm"
)

type Ability struct {
    gorm.Model

    Name string `gorm:"size:128;unique_index;not null"`
    RoleAbilities []RoleAbility
}

func (ability Ability) GetById(id interface{}) error {
    SpaceDock.Database.First(&ability, id)
    if ability.Name != "" {
        return errors.New("Invalid ability ID")
    }
    return nil
}