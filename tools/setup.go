/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package main

import (
    "SpaceDock"
    "SpaceDock/objects"
    "archive/zip"
    "flag"
    "github.com/kennygrant/sanitize"
    "github.com/spf13/cast"
    "os"
    "path/filepath"
    "regexp"
    "strconv"
    "strings"
)

func main() {

    // Setup an Administrator
    NewDummyUser("Administrator", "admin", "admin@example.com", true)

    // Check if we should add dummy data
    p := flag.Bool("dummy", true, "Adds dummy data")
    flag.Parse()
    if *p {
        space_dock_user := NewDummyUser("SpaceDockUser", "user", "user@example.com", false)

        // Game 1
        ksp := NewDummyGame("Kerbal Space Program", "kerbal-space-program", "Squad MX")
        ksp_1 := NewDummyVersion(ksp, "1.2.1", false)
        ksp_2 := NewDummyVersion(ksp, "1.2.2", false)
        ksp_3 := NewDummyVersion(ksp, "1.2.9", true)

        // Game 2
        fac := NewDummyGame("Factorio", "factorio", "Wube Software")
        NewDummyVersion(fac, "0.12", false)

        // Game Admins
        ksp_game_admin := NewDummyGameAdmin("GameAdminKSP", "gameadminksp", "gameadminksp@example.com", ksp)
        NewDummyGameAdmin("GameAdminFAC", "gameadminfac", "gameadminfac@example.com", fac)

        // Mods
        mod_ksp_1 := NewDummyMod("DarkMultiPlayer", space_dock_user, ksp, "MIT")
        mod_ksp_2 := NewDummyMod("CookieEngine", ksp_game_admin, ksp, "GPL")

        // Versions
        NewDummyModVersion(mod_ksp_1, "0.1", ksp, ksp_1, false)
        NewDummyModVersion(mod_ksp_1, "0.2", ksp, ksp_2, false)
        NewDummyModVersion(mod_ksp_1, "0.3", ksp, ksp_3, true)
        NewDummyModVersion(mod_ksp_2, "1.2", ksp, ksp_2, false)
        NewDummyModVersion(mod_ksp_2, "1.3", ksp, ksp_3, true)
    }

}

func AddAbilityRe(role *objects.Role, expression string) {
    var abilities []objects.Ability
    SpaceDock.Database.Find(&abilities)
    for _,element := range abilities {
        if ok,_ := regexp.MatchString(expression, element.Name); ok {
            role.AddAbility(element.Name)
        }
    }
}

func NewDummyUser(name string, password string, email string, admin bool) *objects.User {
    user := objects.NewUser(name, email, password)
    SpaceDock.Database.Save(user)

    // Setup roles
    role := user.AddRole(user.Username)
    role.AddAbility("user-edit")
    role.AddAbility("mods-add")
    role.AddAbility("lists-add")
    role.AddAbility("logged-in")
    role.AddParam("user-edit", "userid", strconv.Itoa(int(user.ID)))
    role.AddParam("mods-add", "gameshort", ".*")
    role.AddParam("lists-add", "gameshort", ".*")
    SpaceDock.Database.Save(role)

    // Admin roles
    if admin {
        admin_role := user.AddRole("admin")
        AddAbilityRe(admin_role, ".*")
        admin_role.AddAbility("mods-invite")
        admin_role.AddAbility("view-users-full")

        // Params
        admin_role.AddParam("admin-impersonate", "userid", ".*")
        admin_role.AddParam("game-edit", "gameshort", ".*")
        admin_role.AddParam("game-add", "pubid", ".*")
        admin_role.AddParam("game-remove", "short", ".*")
        admin_role.AddParam("mods-feature", "gameshort", ".*")
        admin_role.AddParam("mods-edit", "gameshort", ".*")
        admin_role.AddParam("mods-add", "gameshort", ".*")
        admin_role.AddParam("mods-remove", "gameshort", ".*")
        admin_role.AddParam("lists-add", "gameshort", ".*")
        admin_role.AddParam("lists-edit", "gameshort", ".*")
        admin_role.AddParam("lists-remove", "gameshort", ".*")
        admin_role.AddParam("publisher-edit", "publid", ".*")
        admin_role.AddParam("token-edit", "tokenid", ".*")
        admin_role.AddParam("token-remove", "tokenid", ".*")
        admin_role.AddParam("user-edit", "userid", ".*")

        SpaceDock.Database.Save(admin_role)
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

func NewDummyVersion(game *objects.Game, name string, beta bool) *objects.GameVersion {
    version := objects.NewGameVersion(name, *game, beta)
    SpaceDock.Database.Save(version)
    return version
}

func NewDummyGameAdmin(name string, password string, email string, game *objects.Game) *objects.User {
    user := NewDummyUser(name, password, email, false)

    // Game specific stuff
    role := user.AddRole(game.Name)
    role.AddAbility("game-edit")
    AddAbilityRe(role,"mods-.*")
    role.AddAbility("mods-invite")
    AddAbilityRe(role,"lists-.*")

    // Params
    role.AddParam("mods-feature", "gameshort", game.Name)
    role.AddParam("game-edit", "gameshort", game.Name)
    role.AddParam("mods-edit", "gameshort", game.Name)
    role.AddParam("mods-add", "gameshort", game.Name)
    role.AddParam("mods-remove", "gameshort", game.Name)
    role.AddParam("lists-add", "gameshort", game.Name)
    role.AddParam("lists-remove", "gameshort", game.Name)
    SpaceDock.Database.Save(role).Save(user)
    return user
}

func NewDummyMod(name string, user *objects.User, game *objects.Game, license string) *objects.Mod {
    mod := objects.NewMod(name, *user, *game, license)
    mod.Published = true
    SpaceDock.Database.Save(mod)

    // Roles
    role := user.AddRole(mod.Name)
    role.AddAbility("mods-edit")
    role.AddAbility("mods-remove")
    role.AddParam("mods-edit", "modid", cast.ToString(mod.ID))
    role.AddParam("mods-remove", "name", mod.Name)
    SpaceDock.Database.Save(role).Save(user)
    return mod
}

func NewDummyModVersion(mod *objects.Mod, friendly_version string, game *objects.Game, version *objects.GameVersion, beta bool) *objects.ModVersion {
    // Path
    user := mod.User
    filename := sanitize.BaseName(mod.Name) + "-" + sanitize.BaseName(version.FriendlyVersion) + ".zip"
    base_path := filepath.Join(sanitize.BaseName(user.Username) + "_" + strconv.Itoa(int(user.ID)), sanitize.BaseName(mod.Name))
    full_path := filepath.Join(SpaceDock.Settings.Storage, base_path)
    os.MkdirAll(full_path, os.ModePerm)
    path := filepath.Join(full_path, filename)

    // Create the object
    modversion := objects.NewModVersion(*mod, friendly_version, *version, "/content/" + strings.Replace(base_path, "\\", "/", -1) + "/" + filename, beta)

    // Save data
    out, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
    zip := zip.NewWriter(out)
    w,_ := zip.Create("SUPRISE.txt")
    w.Write([]byte("As it seems, you downloaded " + mod.Name + " " + friendly_version))
    zip.Flush()
    zip.Close()
    out.Close()

    // Commit
    SpaceDock.Database.Save(modversion)
    if !beta {
        mod.DefaultVersion = *modversion
        mod.DefaultVersionID = modversion.ID
        SpaceDock.Database.Save(mod)
    }
    return modversion
}