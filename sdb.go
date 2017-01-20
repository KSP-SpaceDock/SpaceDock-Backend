/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package main

import "SpaceDock"

/*
 The entrypoint for the spacedock application.
 Instead of running significant code here, we pass this task to the spacedock package
*/
func main() {
    SpaceDock.Run()
}
