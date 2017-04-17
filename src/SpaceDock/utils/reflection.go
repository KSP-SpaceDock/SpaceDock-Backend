/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package utils

import (
    "github.com/fatih/structs"
    "reflect"
)

func ReadField(value *interface{}, field string) interface{} {
    v := reflect.ValueOf(*value)
    y := v.FieldByName(field)
    return y.Interface()
}

/*
 Return codes:
    0: Everything is ok
    1: Tried to edit a field that is locked
    2: Tried to patch a field that doesnt exist
    3: Invalid type
 */
func EditObject(value interface{}, data map[string]interface{}) int {
    obj := ToMap(value)
    lock := []string{}
    for _,f := range structs.Fields(value) {
        if f.Tag("spacedock") == "lock" {
            lock = append(lock, f.Tag("json"))
        }
    }
    code := editObjectInternal(&obj, data, lock, false)
    if code != 0 {
        return code
    }
    if FromMap(value, obj) != nil {
        return 3
    }
    return 0
}

func editObjectInternal(obj *map[string]interface{}, patch map[string]interface{}, lock []string, allowAdd bool) int {
    c := 0
    for field := range patch {
        if lock != nil {
            if e,_ := ArrayContains(field, lock); e {
                return 1
            }
        }
        if _,ok := (*obj)[field]; ok {
            c = c + 1
            if (*obj)[field] == nil || reflect.TypeOf((*obj)[field]) != reflect.TypeOf(patch[field]) {
                return 3
            }
            if reflect.TypeOf((*obj)[field]).Kind() != reflect.Map {
                (*obj)[field] = patch[field]
            } else {
                o := (*obj)[field].(map[string]interface{})
                code := editObjectInternal(&o, patch[field].(map[string]interface{}), nil, true)
                (*obj)[field] = o
                if code != 0 {
                    return code
                }
            }
        } else if allowAdd {
            (*obj)[field] = patch[field]
            c = c + 1
        }
    }
    if c != len(patch) {
        return 2
    }
    return 0
}
