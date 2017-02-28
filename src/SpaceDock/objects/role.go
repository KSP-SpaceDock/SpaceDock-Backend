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
    "encoding/json"
    "errors"
    "github.com/jinzhu/gorm"
)

type Role struct {
    gorm.Model
    MetaObject

    Name          string `gorm:"size:128;unique_index;not null"`
    Params        string `gorm:"size:4096"`
    roleUsers     []RoleUser
    roleAbilities []RoleAbility
}

func (role *Role) GetById(id interface{}) error {
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
        ability.Meta = "{}"
        SpaceDock.Database.Save(&ability)
    }
    ra := RoleAbility {}
    SpaceDock.Database.Where("role_id = ?", role.ID).Where("ability_id = ?", ability.ID).First(&ra)
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
    SpaceDock.Database.Where("role_id = ?", role.ID).Where("ability_id = ?", ability.ID).First(&ra)
    if ra.RoleID == role.ID && ra.AbilityID == ability.ID {
        SpaceDock.Database.Delete(&ra)
    }
}

func (role Role) GetAbilities() []Ability {
    value := make([]Ability, len(role.roleAbilities))
    for index,element := range role.roleAbilities {
        ability := Ability {}
        SpaceDock.Database.First(&ability, element.AbilityID)
        value[index] = ability
    }
    return value
}

func (role Role) GetParams(ability string, param string) []string {
    var temp map[string]map[string][]string
    err := json.Unmarshal([]byte(role.Params), &temp)
    if err != nil {
        return nil
    }
    if _,ok := temp[ability]; ok {
        if _,ok := temp[ability][param]; ok {
            return temp[ability][param]
        }
    }
    return nil
}

func (role *Role) AddParam(ability string, param string, value string) error  {
    var temp map[string]map[string][]string
    err := json.Unmarshal([]byte(role.Params), &temp)
    if err != nil {
        return err
    }
    if _,ok := temp[ability]; !ok {
        temp[ability] = map[string][]string{}
    }
    if _,ok := temp[ability][param]; !ok {
        temp[ability][param] = []string{}
    }
    if ok,_ := utils.ArrayContains(value, temp[ability][param]); !ok {
        temp[ability][param] = append(temp[ability][param], value)
    }
    val,err := json.Marshal(temp)
    if err != nil {
        return err
    }
    role.Params = string(val)
    return nil
}

func (role *Role) RemoveParam(ability string, param string, value string) error {
    var temp map[string]map[string][]string
    err := json.Unmarshal([]byte(role.Params), &temp)
    if err != nil {
        return err
    }
    if _,ok := temp[ability]; ok {
        if _,ok := temp[ability][param]; ok {
            if ok,index := utils.ArrayContains(value, temp[ability][param]); ok {
                temp[ability][param] = append(temp[ability][param][:index], temp[ability][param][index+1:]...)
            }
        }
    }
    val,err := json.Marshal(temp)
    if err != nil {
        return err
    }
    role.Params = string(val)
    return nil
}