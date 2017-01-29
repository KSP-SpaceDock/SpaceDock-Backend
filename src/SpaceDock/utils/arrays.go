/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package utils

import (
    "log"
    "reflect"
    "regexp"
)

func ArrayContains(val interface{}, array interface{}) (exists bool, index int) {
    exists = false
    index = -1

    switch reflect.TypeOf(array).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(array)

        for i := 0; i < s.Len(); i++ {
            if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
                index = i
                exists = true
                return
            }
        }
    }

    return
}

func ArrayContainsRe(itr []string, value string) bool {
    if itr == nil || value == "" {
        return false
    }

    r,err := regexp.Compile(value)

    if err != nil {
        log.Fatalf("Invalid regular expression detected: %s", value)
        return false
    }

    for _,element := range itr {
        if r.MatchString(element) {
            return true
        }
    }
    return false
}
