package main

import (
    _ "github.com/StollD/iris-cache"
    _ "github.com/fatih/structs"
    _ "github.com/go-gomail/gomail"
    _ "github.com/iris-contrib/middleware/cors"
    _ "github.com/jameskeane/bcrypt"
    _ "github.com/jinzhu/configor"
    _ "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mssql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/kennygrant/sanitize"
    _ "github.com/spf13/cast"
    _ "github.com/ulule/limiter"
    _ "gopkg.in/kataras/iris.v6"
    _ "gopkg.in/kataras/iris.v6/adaptors/httprouter"
    _ "gopkg.in/kataras/iris.v6/adaptors/sessions"
    _ "gopkg.in/kataras/iris.v6/middleware/logger"
)