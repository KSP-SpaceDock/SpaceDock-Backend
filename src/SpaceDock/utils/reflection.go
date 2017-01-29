/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package utils

import "reflect"

func ReadField(value *interface{}, field string) interface{} {
    v := reflect.ValueOf(*value)
    y := v.FieldByName(field)
    return y.Interface()
}
