/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package utils

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/go-gomail/gomail"
    "github.com/kennygrant/sanitize"
    "io/ioutil"
    "log"
    "strconv"
    "strings"
)

func SendMail(sender string, recipients []string, subject string, message string, important bool) {
    if app.Settings.SmtpHost == "" {
        return
    }
    srv := gomail.NewDialer(app.Settings.SmtpHost, app.Settings.SmtpPort, app.Settings.SmtpUser, app.Settings.SmtpPassword)
    m := gomail.NewMessage()
    if important {
        m.SetHeader("X-MC-Important", "true")
    }
    m.SetHeader("X-MC-PreserveRecipients", "false")
    m.SetHeader("Subject", subject)
    m.SetHeader("From", sender)
    if len(recipients) > 1 {
        m.SetHeader("Precedence", "bulk")
        m.SetHeader("To", "undisclosed-recipients:;")
    } else {
        m.SetHeader("To", recipients[0])
    }
    m.SetBody("text/plain", message)
    sc,err := srv.Dial()
    if err != nil {
        log.Printf("Error while sending mail: %s", err)
        return
    }
    defer sc.Close()
    sc.Send(sender, recipients, m)
    log.Printf("Sending email from %s to %d recipients", sender, len(recipients))
}

func SendConfirmation(userConfirmation string, userUsername string, userEmail string, followMod string) {
    buffer,err := ioutil.ReadFile("emails/confirm-account")
    if err != nil {
        log.Printf("Error while reading Email Template confirm-account: %s", err)
        return
    }
    confirmation := userConfirmation
    if followMod != "" {
        confirmation += "?f=" + followMod
    }
    data := map[string]interface{}{
        "site_name": app.Settings.SiteName,
        "username": userUsername,
        "domain": app.Settings.Domain,
        "confirmation": confirmation,
    }
    text := string(buffer)
    s := Format(text, data)
    go SendMail(app.Settings.SupportMail, []string{userEmail}, "Welcome to " + app.Settings.SiteName + "!", s, true)
}

func SendReset(userUsername string, userPasswordReset string, userEmail string) {
    buffer,err := ioutil.ReadFile("emails/password-reset")
    if err != nil {
        log.Printf("Error while reading Email Template password-reset: %s", err)
        return
    }
    data := map[string]interface{}{
        "site_name": app.Settings.SiteName,
        "username": userUsername,
        "domain": app.Settings.Domain,
        "confirmation": userPasswordReset,
    }
    text := string(buffer)
    s := Format(text, data)
    go SendMail(app.Settings.SupportMail, []string{userEmail}, "Reset your password on " + app.Settings.SiteName, s, true)
}

func SendGrantNotice(userUsername string, modUsername string, modName string, modID uint, userEmail string, modURL string) {
    buffer,err := ioutil.ReadFile("emails/grant-notice")
    if err != nil {
        log.Printf("Error while reading Email Template grant-notice: %s", err)
        return
    }
    data := map[string]interface{}{
        "username": userUsername,
        "mod_username": modUsername,
        "mod_name": modName,
        "site_name": app.Settings.SiteName,
        "domain": app.Settings.Domain,
        "url": create_mod_url(modID, sanitize.BaseName(modName)[:64], modURL),
    }
    text := string(buffer)
    s := Format(text, data)
    go SendMail(app.Settings.SupportMail, []string{userEmail}, "You've been asked to co-author a mod on " + app.Settings.SiteName, s, true)
}

func SendUpdateNotification(followers []string, changelog string, username string, friendly_version string, modname string, modID uint, modURL string, gamename string, gameversion string) {
    buffer,err := ioutil.ReadFile("emails/mod-updated")
    if err != nil {
        log.Printf("Error while reading Email Template mod-updated: %s", err)
        return
    }
    changelog = strings.Replace(changelog, "\n", "\n    ", -1)
    if len(followers) == 0 {
        return
    }
    data := map[string]interface{}{
        "username": username,
        "friendly_version": friendly_version,
        "mod_name": modname,
        "site_name": app.Settings.SiteName,
        "changelog": changelog,
        "domain": app.Settings.Domain,
        "url": create_mod_url(modID, sanitize.BaseName(modname)[:64], modURL),
        "game_name": gamename,
        "gameversion": gameversion,
    }
    text := string(buffer)
    s := Format(text, data)
    go SendMail(app.Settings.SupportMail, followers, username + " has just updated " + modname + "!", s, true)
}

func SendAutoUpdateNotification(followers []string, changelog string, username string, friendly_version string, modname string, modID uint, modURL string, gamename string, gameversion string) {
    buffer,err := ioutil.ReadFile("emails/mod-autoupdated")
    if err != nil {
        log.Printf("Error while reading Email Template mod-autoupdated: %s", err)
        return
    }
    changelog = strings.Replace(changelog, "\n", "\n    ", -1)
    if len(followers) == 0 {
        return
    }
    data := map[string]interface{}{
        "username": username,
        "friendly_version": friendly_version,
        "mod_name": modname,
        "game_name": gamename,
        "gameversion": gameversion,
        "domain": app.Settings.Domain,
        "url": create_mod_url(modID, sanitize.BaseName(modname)[:64], modURL),
    }
    text := string(buffer)
    s := Format(text, data)
    go SendMail(app.Settings.SupportMail, followers, modname + " is compatible with " + gamename + " " + gameversion + "!", s, true)
}

func create_mod_url(id uint, name string, modURL string) string {
    if modURL == "" {
        modURL = app.Settings.ModUrl
    }
    return strings.Replace(strings.Replace(modURL, "{id}", strconv.Itoa(int(id)), -1), "{name}", name, -1)
}