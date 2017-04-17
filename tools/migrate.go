/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package tools

import (
    "database/sql"
    "encoding/json"
    "fmt"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    "time"
)

/* Insert appropriate values here */
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

    // Clear database
    newDB.Exec("DELETE * FROM users")
    newDB.Exec("DELETE * FROM role_users")

    // Featured
    fmt.Print("Migrating featured mods")
    rows, err := oldDB.Query("SELECT * FROM featured")
    if err != nil {
        panic(err)
    }
    data := SQLToMap(rows)
    tx, _ := newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d", element["id"])
        newDB.Exec("INSERT INTO featured (created_at, updated_at, mod_id, meta) VALUES (?, ?, ?, ?)",
            element["created"], element["created"], element["mod_id"], "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Users
    fmt.Print("Migrating users")
    rows, err = oldDB.Query("SELECT * FROM users")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d (%s)", element["id"], element["name"])
        newDB.Exec("INSERT INTO users (id, created_at, updated_at, username, email, show_email, public, password, description, confirmation, password_reset, password_reset_expiry, meta) VALUES (?,?,?,?,?,?,?,?,?,?,?)",
            element["id"], element["created"], element["created"], element["username"], element["email"], false,
            element["public"], element["password"], element["description"], element["confirmation"],
            element["passwordReset"], element["passwordResetExpiry"], DumpJSON(map[string]interface{} {
                "forumUsername": element["forumUsername"],
                "ircNick": element["ircNick"],
                "twitterUsername": element["twitterUsername"],
                "redditUsername": element["redditUsername"],
                "background": element["backgroundMedia"],
            }))
        if element["admin"].(bool) {
            newDB.Exec("INSERT INTO role_users (role_id, user_id) VALUES (?,?)", admin_role_id, element["id"])
        }
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Publisher
    fmt.Print("Migrating publishers")
    rows, err = oldDB.Query("SELECT * FROM publisher")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d (%s)", element["id"], element["name"])
        newDB.Exec("INSERT INTO publishers (id, created_at, updated_at, name, description, short_description, meta) VALUES (?,?,?,?,?,?,?)",
            element["id"], element["created"], element["updated"], element["name"], element["description"],
            element["short_description"], DumpJSON(map[string]interface{} {
                "link": element["link"],
                "background": element["background"],
            }))
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Game
    fmt.Print("Migrating games")
    rows, err = oldDB.Query("SELECT * FROM game")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d (%s)", element["id"], element["short"])
        newDB.Exec("INSERT INTO games (id, created_at, updated_at, name, active, altname, rating, releasedate, short, publisher_id, description, short_description, meta) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
            element["id"], element["created"], element["updated"], element["name"], element["active"], element["altname"],
            element["rating"], element["releasedate"], element["short"], element["publisher_id"], element["description"],
            element["short_description"], DumpJSON(map[string]interface{} {
                "link": element["link"],
                "background": element["background"],
            }))
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Mod
    fmt.Print("Migrating mods")
    rows, err = oldDB.Query("SELECT * FROM mod")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d (%s)", element["id"], element["name"])
        newDB.Exec("INSERT INTO mods (id, created_at, updated_at, user_id, game_id, name, description, short_description, approved, published, license, default_version_id, total_score, download_count, meta) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
            element["id"], element["created"], element["updated"], element["user_id"], element["game_id"], element["name"],
            element["description"], element["short_description"], true, element["published"], element["license"],
            element["default_version_id"], 0, element["download_count"], DumpJSON(map[string]interface{} {
                "ckan": element["ckan"],
                "source_link": element["source_link"],
                "background": element["background"],
            }))
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Mod Followers
    fmt.Print("Migrating followers")
    rows, err = oldDB.Query("SELECT * FROM mod_followers")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d - %d", element["user_id"], element["mod_id"])
        newDB.Exec("INSERT INTO mod_followers (user_id, mod_id) VALUES (?,?)", element["user_id"], element["mod_id"])
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Modlist
    fmt.Print("Migrating mod lists")
    rows, err = oldDB.Query("SELECT * FROM modlist")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d (%s)", element["id"], element["name"])
        newDB.Exec("INSERT INTO mod_lists (id, created_at, updated_at, user_id, game_id, name, description, short_description, meta) VALUES (?,?,?,?,?,?,?,?,?)",
            element["id"], element["created"], element["created"], element["user_id"], element["game_id"], element["name"],
            element["description"], element["short_description"], DumpJSON(map[string]interface{} {
                "background": element["background"],
            }))
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Modlist Item
    fmt.Print("Migrating modlist items")
    rows, err = oldDB.Query("SELECT * FROM modlistitem")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d", element["id"])
        newDB.Exec("INSERT INTO mod_list_items (created_at, updated_at, mod_id, mod_list_id, sort_index, meta) VALUES (?,?,?,?,?,?)",
            time.Now(), time.Now(), element["mod_id"], element["mod_list_id"], element["sort_index"], "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Shared Authors
    fmt.Print("Migrating shared authors")
    rows, err = oldDB.Query("SELECT * FROM sharedauthor")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d", element["id"])
        newDB.Exec("INSERT INTO shared_authors (created_at, updated_at, mod_id, user_id, accepted, meta) VALUES (?,?,?,?,?,?)",
            time.Now(), time.Now(), element["mod_id"], element["user_id"], element["accepted"], "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Download events
    fmt.Print("Migrating download events")
    rows, err = oldDB.Query("SELECT * FROM downloadevent")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d", element["id"])
        newDB.Exec("INSERT INTO download_event (created_at, updated_at, mod_id, version_id, downloads, meta) VALUES (?,?,?,?,?,?)",
            element["created"], element["created"], element["mod_id"], element["version_id"], element["downloads"], "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Follow events
    fmt.Print("Migrating follow events")
    rows, err = oldDB.Query("SELECT * FROM followevent")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d", element["id"])
        newDB.Exec("INSERT INTO follow_event (created_at, updated_at, mod_id, events, delta, meta) VALUES (?,?,?,?,?,?)",
            element["created"], element["created"], element["mod_id"], element["events"], element["delta"], "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Referral events
    fmt.Print("Migrating referral events")
    rows, err = oldDB.Query("SELECT * FROM referralevent")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d", element["id"])
        newDB.Exec("INSERT INTO referral_event (created_at, updated_at, mod_id, events, host, meta) VALUES (?,?,?,?,?,?)",
            element["created"], element["created"], element["mod_id"], element["events"], element["host"], "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Mod versions
    fmt.Print("Migrating mod versions")
    rows, err = oldDB.Query("SELECT * FROM modversion")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d (%s)", element["id"], element["friendly_version"])
        newDB.Exec("INSERT INTO mod_versions (created_at, updated_at, mod_id, friendly_version, beta, game_version_id, download_path, changelog, sort_index, file_size, meta) VALUES (?,?,?,?,?,?,,?,?,?,?)",
            element["created"], element["created"], element["mod_id"], element["friendly_version"], false, element["gameversion_id"], element["download_path"], element["changelog"], element["sort_index"], 0, "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    // Game versions
    fmt.Print("Migrating game versions")
    rows, err = oldDB.Query("SELECT * FROM gameversion")
    if err != nil {
        panic(err)
    }
    data = SQLToMap(rows)
    tx, _ = newDB.Begin()
    for _,element := range data {
        fmt.Printf("   Migrating Entry %d (%s)", element["id"], element["friendly_version"])
        newDB.Exec("INSERT INTO game_versions (created_at, updated_at, game_id, friendly_version, beta, meta) VALUES (?,?,?,?,?,?)",
            element["created"], element["created"], element["game_id"], element["friendly_version"], false, "{}")
    }
    tx.Commit()
    rows.Close()
    fmt.Print("")

    fmt.Print("Migration completed. You might have to fix the auto_increments of the tables. Sorry.")
    newDB.Close()
    oldDB.Close()
}