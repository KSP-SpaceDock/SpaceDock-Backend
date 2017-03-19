/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package routes

import (
    "SpaceDock"
    "SpaceDock/middleware"
    "SpaceDock/objects"
    "SpaceDock/utils"
    "github.com/spf13/cast"
    "gopkg.in/kataras/iris.v6"
)

/*
 Registers the routes for the publisher section
 */
func PublisherRegister() {
    Register(GET, "/api/publishers", publishers_list)
    Register(GET, "/api/publishers/:pubid", publishers_info)
    Register(PUT, "/api/publishers/:pubid",
        middleware.NeedsPermission("publisher-edit", true, "pubid"),
        edit_publisher,
    )
    Register(POST, "/api/publishers",
        middleware.NeedsPermission("publisher-add", true),
        add_publisher,
    )
    Register(DELETE, "/api/publishers/:pubid",
        middleware.NeedsPermission("publisher-remove", true),
        remove_publisher,
    )
}

/*
 Path: /api/publishers
 Method: GET
 Description: Outputs all publishers known by the application
 */
func publishers_list(ctx *iris.Context) {
    var publishers []objects.Publisher
    SpaceDock.Database.Find(&publishers)
    output := make([]map[string]interface{}, len(publishers))
    for i,element := range publishers {
        output[i] = utils.ToMap(element)
    }
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": len(publishers), "data": output})
}

/*
 Path: /api/publishers/:pubid
 Method: GET
 Description: Outputs detailed infos for one publisher
 */
func publishers_info(ctx *iris.Context) {
    pubid := cast.ToUint(ctx.GetString("pubid"))

    // Get the publisher
    pub := &objects.Publisher{}
    SpaceDock.Database.Where("id = ?", pubid).First(pub)
    if pub.ID != pubid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The publisher ID is invalid").Code(2110))
        return
    }

    // Return the info
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(pub)})
}

/*
 Path: /api/publishers/:pubid
 Method: PUT
 Description: Edits a publisher, based on the request parameters. Required fields: data
 Abilities: publisher-edit
 */
func edit_publisher(ctx *iris.Context) {
    pubid := cast.ToUint(ctx.GetString("pubid"))

    // Get the publisher
    pub := &objects.Publisher{}
    SpaceDock.Database.Where("id = ?", pubid).First(pub)
    if pub.ID != pubid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The publisher ID is invalid").Code(2110))
        return
    }

    // Edit the publisher
    code := utils.EditObject(pub, utils.GetFullJSON(ctx))
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
    SpaceDock.Database.Save(pub)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(pub)})
}

/*
 Path: /api/publishers/
 Method: POST
 Description: Adds a publisher, based on the request parameters. Required fields: name
 Abilities: publisher-add
 */
func add_publisher(ctx *iris.Context) {
    name := cast.ToString(utils.GetJSON(ctx, "name"))

    // Get the publisher
    pub := &objects.Publisher{}
    SpaceDock.Database.Where("name = ?", name).First(pub)
    if pub.Name == name {
        utils.WriteJSON(ctx, iris.StatusBadRequest, utils.Error("A publisher with this name already exists.").Code(2000))
        return
    }

    // Add the publisher
    pub = objects.NewPublisher(name)
    SpaceDock.Database.Save(pub)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false, "count": 1, "data": utils.ToMap(pub)})
}

/*
 Path: /api/publishers/:pubid
 Method: DELETE
 Description: Removes a game from existence. Required fields: pubid
 Abilities: publisher-remove
 */
func remove_publisher(ctx *iris.Context) {
    pubid := cast.ToUint(ctx.GetString("pubid"))

    // Get the publisher
    pub := &objects.Publisher{}
    SpaceDock.Database.Where("id = ?", pubid).First(pub)
    if pub.ID != pubid {
        utils.WriteJSON(ctx, iris.StatusNotFound, utils.Error("The publisher ID is invalid").Code(2110))
        return
    }

    // Delete the publisher
    SpaceDock.Database.Delete(pub)
    utils.WriteJSON(ctx, iris.StatusOK, iris.Map{"error": false})
}