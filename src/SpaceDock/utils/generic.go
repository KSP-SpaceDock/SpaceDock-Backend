/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package utils

/*
 Stupidly simple, but I want to inline conditions
 */
func Ternary(condition bool, tVal interface{}, fVal interface{}) interface{} {
    if (condition) {
        return tVal
    }
    return fVal
}