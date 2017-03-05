/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */


package utils

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "github.com/fatih/structs"
)

func RandomHex(n int) (string, error) {
    bytes := make([]byte, n)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return hex.EncodeToString(bytes), nil
}

func LoadJSON(data string) map[string]interface{} {
    var temp map[string]interface{}
    err := json.Unmarshal([]byte(data), &temp)
    if err != nil {
        return nil
    }
    return temp
}

func DumpJSON(data interface{}) string {
    buff, err := json.Marshal(data)
    if err != nil {
        return "{}"
    }
    return string(buff)
}

func ToMap(data interface{}) map[string]interface{} {
    m := LoadJSON(DumpJSON(data))
    m["meta"] = LoadJSON((m["meta"].(string)))
    for _,element := range structs.Fields(data) {
        if element.Tag("spacedock") == "json" {
            m[element.Tag("json")] = LoadJSON((m[element.Tag("json")].(string)))
        }
    }
    return m
}

func FromMap(data interface{}, values map[string]interface{}) error {
    values["meta"] = DumpJSON(values["meta"])
    for _,element := range structs.Fields(data) {
        if element.Tag("spacedock") == "json" {
            values[element.Tag("json")] = DumpJSON(values[element.Tag("json")])
        }
    }
    buff := []byte(DumpJSON(values))
    err := json.Unmarshal(buff, data)
    return err
}