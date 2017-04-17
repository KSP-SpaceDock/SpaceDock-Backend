/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package utils

import (
    "crypto/md5"
    "fmt"
)

func cacheKey(url string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(url)))
}

var InvalidFunc func (cacheKey string)

func InvalidateCache(url string, values ...interface{}) {
    InvalidFunc(cacheKey(fmt.Sprintf(url, values...)))
    InvalidFunc(cacheKey(fmt.Sprintf(url + "/", values...)))
}

func ClearFeaturedCache(gameshort string) {
    InvalidateCache("/api/featured")
    if gameshort != "" {
        InvalidateCache("/api/featured/%s", gameshort)
    }
}

func ClearGameCache(gameshort string, version string) {
    InvalidateCache("/api/games")
    if gameshort != "" {
        InvalidateCache("/api/games/%s", gameshort)
        InvalidateCache("/api/games/%s/versions", gameshort)
    }
    if version != "" {
        InvalidateCache("/api/games/%s/versions/%s", gameshort, version)
    }
}

func ClearModListCache(gameshort string, listid uint) {
    InvalidateCache("/api/lists")
    if gameshort != "" {
        InvalidateCache("/api/lists/%s", gameshort)
    }
    if listid != 0 {
        InvalidateCache("/api/lists/%s/%d", gameshort, listid)
    }
}

func ClearModCache(gameshort string, modid uint) {
    InvalidateCache("/api/mods")
    if gameshort != "" {
        InvalidateCache("/api/mods/%s", gameshort)
    }
    if modid != 0 {
        InvalidateCache("/api/mods/%s/%d", gameshort, modid)
        InvalidateCache("/api/mods/%s/%d/versions", gameshort, modid)
    }
}

func ClearPublisherCache(pubid uint) {
    InvalidateCache("/api/publishers")
    if pubid != 0 {
        InvalidateCache("/api/publishers/%d", pubid)
    }
}

func ClearUserCache(userid uint) {
    InvalidateCache("/api/users")
    if userid != 0 {
        InvalidateCache("/api/users/%d", userid)
    }
}