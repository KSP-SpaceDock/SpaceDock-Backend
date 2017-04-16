/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package tools

import (
    "database/sql"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "encoding/json"
)

/* Insert appropreate values here */
const driver_old_SD = "postgres"
const connection_old_SD = ""
const driver_new_SD = "postgres"
const connection_new_SD = ""
const admin_role_id = 1

func SQLToMap(rows *sql.Rows) []map[string]interface{} {
    cols,_ := rows.Columns()
    result := []map[string]interface{}{}
    for rows.Next() {
        // Create a slice of interface{}'s to represent each column,
        // and a second slice to contain pointers to each item in the columns slice.
        columns := make([]interface{}, len(cols))
        columnPointers := make([]interface{}, len(cols))
        for i, _ := range columns {
            columnPointers[i] = &columns[i]
        }

        // Scan the result into the column pointers...
        if err := rows.Scan(columnPointers...); err != nil {
            panic(err)
        }

        // Create our map, and retrieve the value for each column from the pointers slice,
        // storing it in the map with the name of the column as the key.
        m := make(map[string]interface{})
        for i, colName := range cols {
            val := columnPointers[i].(*interface{})
            m[colName] = *val
        }

        // Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
        result = append(result, m)
    }
    return result
}

func DumpJSON(data interface{}) string {
    buff, err := json.Marshal(data)
    if err != nil {
        return "{}"
    }
    return string(buff)
}

func main() {
    oldDB, err := sql.Open(driver_old_SD, connection_old_SD)
    if err != nil {
        panic(err)
    }
    newDB, err := sql.Open(driver_new_SD, connection_new_SD)
    if err != nil {
        panic(err)
    }

    // Transfer entries

    // Featured
    rows, err := oldDB.Query("SELECT * FROM featured")
    if err != nil {
        panic(err)
    }
    data := SQLToMap(rows)
    tx, _ := newDB.Begin()
    for _,element := range data {
        newDB.Exec("INSERT INTO featured (created_at, updated_at, mod_id, meta) VALUES (?, ?, ?, ?)", element["created"], element["created"], element["mod_id"], "{}")
    }
    tx.Commit()

    // Users
    rows, err = oldDB.Query("SELECT * FROM users")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        newDB.Exec("INSERT INTO users (created_at, updated_at, username, email, show_email, public, password, description, confirmation, password_reset, password_reset_expiry, meta) VALUES (?,?,?,?,?,?,?,?,?,?)",
            element["created"], element["created"], element["username"], element["email"], element["showEmail"] /* TODO: Check */, element["public"], element["password"], element["description"], element["confirmation"],
            element["passwordReset"], element["passwordResetExpiry"], DumpJSON(map[string]interface{} {
                "forumUsername": element["forumUsername"],
                "ircNick": element["ircNick"],
                "twitterUsername": element["twitterUsername"],
                "redditUsername": element["redditUsername"],
                "youtubeUsername": element["youtubeUsername"],
                "twitchUsername": element["twitchUsername"],
                "location": element["location"],
                "facebookUsername": element["facebookUsername"],
                "background": element["backgroundMedia"],
            }))
        if element["admin"].(bool) {
            newDB.Exec("INSERT INTO role_user (role_id, user_id) VALUES (?,?)", admin_role_id, element["id"])
        }
    }
    tx.Commit()

    // Ratings
    rows, err = oldDB.Query("SELECT * FROM rating")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        newDB.Exec("INSERT INTO ratings (created_at, updated_at, user_id, mod_id, score) VALUES (?,?,?,?,?)", element["created"], element["updated"], element["user_id"], element["mod_id"], element["score"])
    }
    tx.Commit()
}