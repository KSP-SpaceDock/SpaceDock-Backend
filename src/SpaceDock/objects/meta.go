/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import (
    "encoding/json"
    "errors"
    "time"
)

type Model struct {
    ID        uint `gorm:"primary_key" json:"id" spacedock:"lock"`
    CreatedAt time.Time `json:"created" spacedock:"lock"`
    UpdatedAt time.Time `json:"updated" spacedock:"lock"`
    DeletedAt *time.Time `sql:"index" json:"-" spacedock:"lock"`
    Meta      string `gorm:"size:4096" json:"meta"`
}

func (meta Model) GetValue(key string) (error,interface{}) {
    var temp map[string]interface{}
    err := json.Unmarshal([]byte(meta.Meta), &temp)
    if err != nil {
        return err, nil
    }
    if _,ok := temp[key]; !ok {
        return errors.New("Invalid key!"),nil
    }
    return nil,temp[key]
}

func (meta *Model) SetValue(key string, value interface{}) error {
    var temp map[string]interface{}
    err := json.Unmarshal([]byte(meta.Meta), &temp)
    if err != nil {
        return err
    }
    temp[key] = value
    buffer,err := json.Marshal(temp)
    if err != nil {
        return err
    }
    meta.Meta = string(buffer)
    return nil
}