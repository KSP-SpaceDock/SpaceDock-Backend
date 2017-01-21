/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package SpaceDock

import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
    "log"
)

/*
 The database handler that provides the API to interact with the the DB
 */
var database gorm.DB

/*
 Establishes the connection to the database
 */
func LoadDatabase() {
    db, err := gorm.Open(settings.Dialect, settings.ConnectionData)
    if err != nil {
        log.Fatalf("* Failed to connect to the database: %s", err)
    }
    database = *db
    log.Print("* Database connection successfull")

    // Init Tables
}

/*
 Creates a table only if it doesn't exist
 */
func CreateTable(models interface{}) {
    if !database.HasTable(models) {
        database.CreateTable(models)
    }
}
