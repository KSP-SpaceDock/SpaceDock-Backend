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
    "errors"
    "crypto/md5"
    "crypto/rand"
)

type Token struct {
    Model

    Token  string `gorm:"size:32" spacedock:"lock"`
}

func NewToken() *Token {
    key := make([]byte, 64)
    _, err := rand.Read(key)
    if err != nil {
        panic(err) // aahhh
    }
    t := &Token{
        Token:string(md5.New().Sum(key)),
    }
    t.Meta = "{}"
    return t
}