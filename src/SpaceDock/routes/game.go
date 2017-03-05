/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package routes

import (
    "SpaceDock"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "gopkg.in/kataras/iris.v6"
    "strconv"
)

/*
 Registers the routes for the game management
 */
func GameRegister() {
    Register(GET, "/api/games/", listgames)
    Register(GET, "/api/games/:gameshort", showgame)
    Register(POST, "/api/games/:gameshort/edit", editgame)
}

/*
 Path:   /api/games/
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
 Path:   /api/games/:gameshort
 Method: GET
 Description: Displays information about a game.
 */
func showgame(ctx *iris.Context) {
    gameshort := ctx.GetString("gameshort")
    var game objects.Game
    SpaceDock.Database.Where("short = ?", gameshort).First(&game)
    if game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
        return
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(game)})
}

/*
 Path:   /api/games/:gameshort/edit
 Method: GET
 Description: Edits a game, based on the request parameters. Required fields: data
 */
func editgame(ctx *iris.Context) {
    gameshort := ctx.GetString("gameshort")
    var game objects.Game
    SpaceDock.Database.Where("short = ?", gameshort).First(&game)
    if game.Short != gameshort {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("The gameshort is invalid.").Code(2125))
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