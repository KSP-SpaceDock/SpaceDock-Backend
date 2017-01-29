/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "github.com/jinzhu/gorm"
)

type RoleAbility struct {
    gorm.Model

    RoleID uint
    AbilityID uint
}

func NewRoleAbility(role Role, ability Ability) *RoleAbility {
    return &RoleAbility{ RoleID: role.ID, AbilityID: ability.ID }
}

func (ra RoleAbility) GetRole() Role {
    role := Role {}
    err := role.GetById(ra.RoleID)
    if err != nil {
        return Role {}
    }
    return role
}

func (ra RoleAbility) GetAbility() Ability {
    ability := Ability {}
    err := ability.GetById(ra.AbilityID)
    if err != nil {
        return Ability {}
    }
    return ability
}
