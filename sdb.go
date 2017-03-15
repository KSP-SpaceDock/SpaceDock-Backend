/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package main

import (
    "SpaceDock"
     "SpaceDock/objects"
    _ "SpaceDock/routes"
    "strconv"
)

/*
 The entrypoint for the spacedock application.
 Instead of running significant code here, we pass this task to the spacedock package
*/
func main() {
    if (SpaceDock.Settings.CreateDefaultDatabase) {
        CreateDefaultData()
    }
    SpaceDock.Run()
}

func CreateDefaultData() {

    // Setup users
    NewDummyUser("Administrator", "admin", "admin@example.com", true)
    NewDummyUser("SpaceDockUser", "user", "user@example.com", false)

    // Setup games
    NewDummyGame("Kerbal Space Program", "kerbal-space-program", "Squad MX")
    NewDummyGame("Factorio", "factorio", "Wube Software")
}

func NewDummyUser(name string, password string, email string, admin bool) *objects.User {
    user := objects.NewUser(name, email, password)
    SpaceDock.Database.Save(user)

    // Setup roles
    role := user.AddRole(user.Username)
    role.AddAbility("user-edit")
    role.AddAbility("mods-add")
    role.AddAbility("packs-add")
    role.AddAbility("logged-in")
    role.AddParam("user-edit", "userid", strconv.Itoa(int(user.ID)))
    role.AddParam("mods-add", "gameshort", ".*")
    role.AddParam("packs-add", "gameshort", ".*")
    SpaceDock.Database.Save(&role)

    // Admin roles
    if admin {
        admin_role := user.AddRole("admin")

        // access.go
        admin_role.AddAbility("access-view")
        admin_role.AddAbility("access-edit")

        // admin.go
        admin_role.AddAbility("admin-impersonate")
        admin_role.AddAbility("admin-confirm")

        // game.go
        admin_role.AddAbility("game-add")
        admin_role.AddAbility("game-edit")
        admin_role.AddAbility("game-remove")

        // Params
        admin_role.AddParam("admin-impersonate", "userid", ".*")
        admin_role.AddParam("game-edit", "gameshort", ".*")
        admin_role.AddParam("game-add", "pubid", ".*")
        admin_role.AddParam("game-remove", "short", ".*")
        /*admin_role.AddParam("mods-feature", "gameshort", ".*")
        admin_role.AddParam("mods-edit", "gameshort", ".*")
        admin_role.AddParam("mods-add", "gameshort", ".*")
        admin_role.AddParam("mods-remove", "gameshort", ".*")
        admin_role.AddParam("packs-add", "gameshort", ".*")
        admin_role.AddParam("packs-remove", "gameshort", ".*")
        admin_role.AddParam("publisher-edit", "publid", ".*")
        admin_role.AddParam("token-edit", "tokenid", ".*")
        admin_role.AddParam("token-remove", "tokenid", ".*")
        admin_role.AddParam("user-edit", "userid", ".*")*/

        SpaceDock.Database.Save(&admin_role)
    }

    // Confirmation
    user.Confirmation = ""
    user.Public = true
    SpaceDock.Database.Save(user)
    return user
}

func NewDummyGame(name string, short string, publisher string) *objects.Game {
    pub := objects.NewPublisher(publisher)
    SpaceDock.Database.Save(pub)

    // Create the game
    game := objects.NewGame(name, *pub, short)
    game.Active = true
    SpaceDock.Database.Save(game)
    return game
}
