/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import "SpaceDock"

/*
 This function creates tables for all datatypes
 */
func init() {
    SpaceDock.CreateTable(&Ability{})
    SpaceDock.CreateTable(&Game{})
    SpaceDock.CreateTable(&GameVersion{})
    SpaceDock.CreateTable(&Mod{})
    SpaceDock.CreateTable(&ModVersion{})
    SpaceDock.CreateTable(&Publisher{})
    SpaceDock.CreateTable(&Rating{})
    SpaceDock.CreateTable(&Role{})
    SpaceDock.CreateTable(&SharedAuthor{})
    SpaceDock.CreateTable(&Token{})
    SpaceDock.CreateTable(&User{})
}