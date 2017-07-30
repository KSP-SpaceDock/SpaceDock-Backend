/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
*/

package main

import (
    "os"
    "flag"
    "fmt"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    _ "github.com/KSP-SpaceDock/SpaceDock-Backend/routes"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/tools"
)

/*
 The entrypoint for the spacedock application.
 Instead of running significant code here, we pass this task to the app package
*/
func main() {
    args := os.Args[1:]

    // Define subcommands
    helpCommand := flag.NewFlagSet("help", flag.ExitOnError)
    setupCommand := flag.NewFlagSet("setup", flag.ExitOnError)
    migrateCommand := flag.NewFlagSet("migrate", flag.ExitOnError)

    // Setup subcommand flags
    dummyData := setupCommand.Bool("dummy", true, "Populates the database with dummy data")

    // Help text that will be displayed when the help flag is passed
    flag.Usage = func() {
        fmt.Println("usage: sdb [command] [options]")
        fmt.Println("")
        fmt.Println("SpaceDock backend application for handling database operations and http routes.")
        fmt.Println("")
        fmt.Println("Use \"sdb help <command>\" for more information about a command.")
        fmt.Println("")
        fmt.Println("    Commands:")
        fmt.Println("")
        fmt.Println("        migrate     converts a pre-split SpaceDock database to the new backend database format")
        fmt.Println("        setup       populates the database with dummy data and an administrator account")
        fmt.Println("")
        fmt.Println("If no subcommand is specified, the backend application will run.")
        fmt.Println("")
    }

    // Check if we passed an argument to the application and try to parse it as a subcommand.
    // If no arguments were passed, simply run the application.
   if len(args) > 0 {
        switch args[0] {
        case "setup":
            setupCommand.Parse(args[1:])
        case "migrate":
            migrateCommand.Parse(args[1:])
        case "help":
            helpCommand.Parse(args[1:])
        default:
            flag.Usage()
            os.Exit(1)
        }
    } else {
        app.Run()
    }

    // Check if the help subcommand was parsed, so it can display some useful information.
    if helpCommand.Parsed() {
        helpArgs := helpCommand.Args()
        helpCommand.Usage = func() {
            // Default help text for help subcommand
            defaultUsage := func() {
                fmt.Println("usage: sdb help <command>")
                fmt.Println("")
                fmt.Println("    Commands:")
                fmt.Println("")
                fmt.Println("        migrate     converts a pre-split SpaceDock database to the new backend database format")
                fmt.Println("        setup       populates the database with dummy data and an administrator account")
                fmt.Println("")
            }

            // Check if we passed a valid subcommand as argument to the help command.
            // If we didn't, print the default help text.
            if len(helpArgs)  == 1 {
                switch helpArgs[0] {
                case "migrate":
                    fmt.Println("usage: sdb migrate")
                    fmt.Println("")
                    fmt.Println("The migrate subcommand will convert an old database to the new backend format.")
                    fmt.Println("")
                case "setup":
                    fmt.Println("usage: sdb setup [-dummy=true|false]")
                    fmt.Println("")
                    fmt.Println("The setup subcommand will add an administrator account, a normal user,")
                    fmt.Println("publisher, game, game admin and some dummy mods to the database.")
                    fmt.Println("")
                    fmt.Println("If you set the dummy flag to false, only an admin account will be added.")
                    fmt.Println("")
                default:
                    defaultUsage()
                }
            } else if len(helpArgs) == 0 {
                defaultUsage()
            }
        }

        helpCommand.Usage()
        os.Exit(1)
    }

    // Check if the setup subcommand was parsed so the app can setup the database
    if setupCommand.Parsed() {
        tools.Setup(*dummyData)
    }

    // Check if the migrate subcommand was parsed so the app can start its migration process
    if migrateCommand.Parsed() {
        tools.MigrateDB()
    }
}

