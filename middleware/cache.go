/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package middleware

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "github.com/StollD/iris-cache"
    "gopkg.in/kataras/iris.v6"
    "gopkg.in/redis.v5"
    "time"
)

var Cache iris.HandlerFunc

func CreateCache() {
    var memoryStore cache.CacheStore
    if app.Settings.StoreType == "memory" {
        memoryStore = cache.NewInMemoryStore()
    } else {
        options,err := redis.ParseURL(app.Settings.RedisConnection)
        if err != nil {
            panic(err)
        }
        memoryStore = cache.NewRedisStore(redis.NewClient(options))
    }
    c := cache.NewCache(cache.CacheConfig{
        AutoRemove:        false,
        CacheTimeDuration: time.Duration(app.Settings.CacheTimeout) * time.Minute,
        ContentType:       cache.ContentTypeJSON,
        IrisGzipEnabled:   true,
        CacheKeyFunc:      cache.RequestPathToMD5,
    }, memoryStore)
    Cache = c.Serve
    utils.InvalidFunc = c.Invalidate
}

