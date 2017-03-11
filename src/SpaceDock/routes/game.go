/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package routes

import (
    "SpaceDock"
    "SpaceDock/middleware"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
    "strconv"
)

/*
 Registers the routes for the game management
 */
func GameRegister() {
    Register(GET, "/api/games/", listgames)
    Register(GET, "/api/games/:gameshort", showgame)
    Register(PUT, "/api/games/:gameshort",
        middleware.NeedsPermission("game-edit", true, "gameshort"),
        editgame,
    )
    Register(POST, "/api/games/",
        middleware.NeedsPermission("game-add", true, "pubid"),
        addgame,
    )
    Register(DELETE, "/api/games/",
        middleware.NeedsPermission("game-remove", true, "pubid"),
        removegame,
    )
}

/*
 Path: /api/games/
 Method: GET
 Description: Displays a list of all games in the database.
 */
func listgames(ctx *iris.Context) {
    var games []objects.Game
    includeInactive := ctx.URLParam("includeInactive")
    val, err  := strconv.ParseBool(includeInactive)
    if (err != nil) && val {
        SpaceDock.Database.Find(&games)
    } else {
        SpaceDock.Database.Where("active = ?", true).Find(&games)
    }
    output := make([]map[string]interface{}, len(games))
    for i,element := range games {
        output[i] = utils.ToMap(element)
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(games), "data": output})
}

/*
 Path: /api/games/:gameshort
 Method: GET
 Description: Displays information about a game.
 */
func showgame(ctx *iris.Context) {
    gameshort := ctx.GetString("gameshort")
    var game objects.Game
    SpaceDock.Database.Where("short = ?", gameshort).First(&game)
    if game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The game does not exist.").Code(2125))
        return
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(game)})
}

/*
 Path: /api/games/:gameshort
 Method: PUT
 Description: Edits a game, based on the request parameters.
 Abilities: game-edit
 */
func editgame(ctx *iris.Context) {
    gameshort := ctx.GetString("gameshort")
    var game objects.Game
    SpaceDock.Database.Where("short = ?", gameshort).First(&game)
    if game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The game does not exist.").Code(2125))
        return
    }

    // Edit the game
    code := utils.EditObject(&game, utils.GetFullJSON(ctx))
    if code == 3 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The value you submitted is invalid").Code(2180))
        return
    } else if code == 2 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that doesn't exist.").Code(3090))
        return
    } else if code == 1 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("You tried to edit a value that is marked as read-only.").Code(3095))
        return
    }
    SpaceDock.Database.Save(&game)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(game)})
}

/*
 Path: /api/games/
 Method: POST
 Description: Adds a new game based on the request parameters. Required fields: name, pubid, short
 Abilities: game-add
 */
func addgame(ctx *iris.Context) {
    name := cast.ToString(utils.GetJSON(ctx, "name"))
    pubid := cast.ToUint(utils.GetJSON(ctx, "pubid"))
    short := cast.ToString(utils.GetJSON(ctx, "short"))

    errors := []string{}
    codes := []int{}

    // Check if the publisher ID is valid
    var publisher objects.Publisher
    SpaceDock.Database.Where("id = ?", pubid).First(&publisher)
    if publisher.ID != pubid {
        errors = append(errors, "The pubid is invalid.")
        codes = append(codes, 2110)
    }
    if name == "" {
        errors = append(errors, "The name is invalid.")
        codes = append(codes, 2117)
    }
    if short == "" {
        errors = append(errors, "The gameshort is invalid.")
        codes = append(codes, 2125)
    }

    // Check if the game already exists
    var game *objects.Game
    SpaceDock.Database.Where("short = ?", short).First(game)
    if game.Short == short {
        errors = append(errors, "The gameshort already exists.")
        codes = append(codes, 2015)
    }
    SpaceDock.Database.Where("name = ?", name).First(&game)
    if game.Name == name {
        errors = append(errors, "The game name already exists.")
        codes = append(codes, 2020)
    }

    // Errors
    if len(errors) > 0 {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error(errors...).Code(codes...))
        return
    }

    // Make a new game
    game = objects.NewGame(name, publisher, short)
    SpaceDock.Database.Save(game)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(*game)})
}

/*
 Path: /api/games/
 Method: DELETE
 Description: Removes a game from existence. Required fields: short
 Abilities: game-remove
 */
func removegame(ctx *iris.Context) {
    short := cast.ToString(utils.GetJSON(ctx, "short"))

    // Check if the game exists
    var game *objects.Game
    SpaceDock.Database.Where("short = ?", short).First(game)
    if game.Short != short {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The game does not exist.").Code(2125))
        return
    }

    // Remove it
    SpaceDock.Database.Delete(game)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}