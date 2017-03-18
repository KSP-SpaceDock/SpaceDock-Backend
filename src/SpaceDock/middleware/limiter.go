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
    "github.com/ulule/limiter"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "strconv"
)

/*
 Creates an iris-style middleware for the limiter HttpFunc
 */
func NewAccessLimiter(limiterObj *limiter.Limiter) iris.HandlerFunc {
    return func (ctx *iris.Context) {
        context, err := limiterObj.Get(limiter.GetIPKey(ctx.Request))
        if err != nil {
            panic(err)
        }

        if limitFunc(context, ctx) {
            utils.WriteJSON(ctx, iris.StatusTooManyRequests, utils.Error("Request limit exceeded").Code(0000))
            return
        } else if !context.Reached {
            ctx.SetHeader("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
            ctx.SetHeader("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
            ctx.SetHeader("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))
        }
        ctx.Next()
        return
    }
}

func limitFunc(context limiter.Context, ctx *iris.Context) bool {
    s_token := ctx.URLParam("token")
    if s_token != "" {
        var token objects.Token
        SpaceDock.Database.Where("token = ?", s_token).First(&token)
        if token.Token != s_token {
            return context.Reached
        }
        _,ips := token.GetValue("ips")
        e,_ := utils.ArrayContains(ctx.RemoteAddr(), cast.ToStringSlice(ips))
        return utils.Ternary(e, false, context.Reached).(bool)
    }
    return context.Reached
}