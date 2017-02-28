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
)

type MetaObject struct {
    Meta string `gorm:"size:4096"`
}

func (meta MetaObject) GetValue(key string) (error,interface{}) {
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

func (meta *MetaObject) SetValue(key string, value interface{}) error {
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