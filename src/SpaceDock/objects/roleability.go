/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

type RoleAbility struct {
    Model

    RoleID uint
    AbilityID uint
}

func NewRoleAbility(role Role, ability Ability) *RoleAbility {
    ra := &RoleAbility{ RoleID: role.ID, AbilityID: ability.ID }
    ra.Meta = "{}"
    return ra
}

func (ra RoleAbility) GetRole() *Role {
    role := Role {}
    err := role.GetById(ra.RoleID)
    if err != nil {
        return nil
    }
    return &role
}

func (ra RoleAbility) GetAbility() *Ability {
    ability := Ability {}
    err := ability.GetById(ra.AbilityID)
    if err != nil {
        return nil
    }
    return &ability
}
