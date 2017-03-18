/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package middleware

import (
    "SpaceDock"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "github.com/KSP-SpaceDock/limiter"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "net/http"
)

/*
 Creates an iris-style middleware for the limiter HttpFunc
 */
func NewAccessLimiter(limiterObj *limiter.Limiter) iris.HandlerFunc {
    return iris.ToHandler(limiter.NewHTTPMiddleware(limiterObj, limitFunc).ServeHTTP)
}

func limitFunc(context limiter.Context, r *http.Request) bool {
    s_token := r.URL.Query().Get("token")
    if s_token != "" {
        var token objects.Token
        SpaceDock.Database.Where("token = ?", s_token).First(&token)
        if token.Token != s_token {
            return context.Reached
        }
        _,ips := token.GetValue("ips")
        e,_ := utils.ArrayContains(r.RemoteAddr, cast.ToStringSlice(ips))
        return utils.Ternary(e, false, context.Reached).(bool)
    }
    return context.Reached
}