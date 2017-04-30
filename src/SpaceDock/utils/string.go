/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */


package utils

import (
    "crypto/rand"
    "encoding/hex"
    "encoding/json"
    "github.com/fatih/structs"
    "fmt"
    "strings"
)

type DataTransformerFunc func(interface{}, map[string]interface{})

var transformers []DataTransformerFunc = []DataTransformerFunc{}

func RegisterDataTransformer(transformer DataTransformerFunc) {
    transformers = append(transformers, transformer)
}

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
    m["meta"] = LoadJSON(m["meta"].(string))
    for _,element := range structs.Fields(data) {
        if strings.Contains(element.Tag("spacedock"),"json") {
            m[element.Tag("json")] = LoadJSON(m[element.Tag("json")].(string))
        }
        if strings.Contains(element.Tag("spacedock"),"tomap") {
            m[element.Tag("json")] = ToMap(element.Value())
        }
    }
    for _,element := range transformers {
        element(data, m)
    }
    return m
}

func FromMap(data interface{}, values map[string]interface{}) error {
    values["meta"] = DumpJSON(values["meta"].(map[string]interface{}))
    for _,element := range structs.Fields(data) {
        if strings.Contains(element.Tag("spacedock"),"json") {
            values[element.Tag("json")] = DumpJSON(values[element.Tag("json")])
        }
    }
    buff := []byte(DumpJSON(values))
    err := json.Unmarshal(buff, data)
    return err
}

func Format(format string, p map[string]interface{}) string {
    args, i := make([]string, len(p)*2), 0
    for k, v := range p {
        args[i] = "{" + k + "}"
        args[i+1] = fmt.Sprint(v)
        i += 2
    }
    return strings.NewReplacer(args...).Replace(format)
}