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

type Role struct {
    gorm.Model

    Name string `gorm:"size:128;unique_index;not null"`
    Params string `gorm:"size:4096"`
    RoleUsers []RoleUser
    RoleAbilities []RoleAbility
}

func (role Role) GetById(id interface{}) error {
    SpaceDock.Database.First(&role, id)
    if role.Name != "" {
        return errors.New("Invalid role ID")
    }
    return nil
}

func (role Role) AddAbility(name string) Ability {
    ability := Ability {}
    SpaceDock.Database.Where("name = ?", name).First(&ability)
    if ability.Name == "" {
        ability.Name = name
        SpaceDock.Database.Save(&ability)
    }
    ra := RoleAbility {}
    SpaceDock.Database.Where("roleid = ?", role.ID).Where("abilityid = ?", ability.ID).First(&ra)
    if ra.RoleID != role.ID || ra.AbilityID != ability.ID {
        SpaceDock.Database.Save(NewRoleAbility(role, ability))
    }
    return ability
}

func (role Role) RemoveAbility(name string) {
    ability := Ability {}
    SpaceDock.Database.Where("name = ?", name).First(&ability)
    if ability.Name == "" {
        return
    }
    ra := RoleAbility {}
    SpaceDock.Database.Where("roleid = ?", role.ID).Where("abilityid = ?", ability.ID).First(&ra)
    if ra.RoleID == role.ID && ra.AbilityID == ability.ID {
        SpaceDock.Database.Delete(&ra)
    }
}

func (role Role) GetAbilities() []Ability {
    value := make([]Ability, len(role.RoleAbilities))
    for index,element := range role.RoleAbilities {
        ability := Ability {}
        SpaceDock.Database.First(&ability, element.AbilityID)
        value[index] = ability
    }
    return value
}