/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import "github.com/KSP-SpaceDock/SpaceDock-Backend/app"

/*
 This function creates tables for all datatypes
 */
func init() {
    app.CreateTable(&Ability{})
    app.CreateTable(&DownloadEvent{})
    app.CreateTable(&FollowEvent{})
    app.CreateTable(&ReferralEvent{})
    app.CreateTable(&Featured{})
    app.CreateTable(&Game{})
    app.CreateTable(&GameVersion{})
    app.CreateTable(&Mod{})
    app.CreateTable(&ModList{})
    app.CreateTable(&ModListItem{})
    app.CreateTable(&ModVersion{})
    app.CreateTable(&Publisher{})
    app.CreateTable(&Rating{})
    app.CreateTable(&Role{})
    app.CreateTable(&SharedAuthor{})
    app.CreateTable(&Token{})
    app.CreateTable(&User{})
}