/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package utils

import "gopkg.in/kataras/iris.v6"

type ErrorMessage struct {
    reasons []string
}

func Error(reasons ...string) ErrorMessage {
    return ErrorMessage{reasons:reasons}
}

func (error ErrorMessage) Code(codes ...int) iris.Map {
    return iris.Map{"error": true, "reasons": error.reasons, "codes": codes}
}
