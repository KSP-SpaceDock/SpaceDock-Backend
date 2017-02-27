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