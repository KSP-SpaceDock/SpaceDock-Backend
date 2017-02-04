/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package utils

import (
    "SpaceDock"
    "SpaceDock/objects"
    "bytes"
    "github.com/go-gomail/gomail"
    "io/ioutil"
    "log"
    "text/template"
)

func SendMail(sender string, recipients []string, subject string, message string, important bool) {
    if SpaceDock.Settings.SmtpHost == "" {
        return
    }
    srv := gomail.NewDialer(SpaceDock.Settings.SmtpHost, SpaceDock.Settings.SmtpPort, SpaceDock.Settings.SmtpUser, SpaceDock.Settings.SmtpPassword)
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
    defer sc.Close()
    if err != nil {
        log.Fatalf("Error while sending mail: %s", err)
        return
    }
    sc.Send(sender, recipients, m)
    log.Printf("Sending email from %s to %d recipients", sender, len(recipients))
}

func SendConfirmation(user objects.User, followMod string) {
    buffer,err := ioutil.ReadFile("emails/confirm-account")
    if err != nil {
        log.Fatalf("Error while reading Email Template confirm-account: %s", err)
        return
    }
    confirmation := user.Confirmation
    if followMod != "" {
        confirmation += "?f=" + followMod
    }
    data := map[string]interface{}{
        "SiteName":     SpaceDock.Settings.SiteName,
        "Username": user.Username,
        "Domain": SpaceDock.Settings.Domain,
        "Confirmation": confirmation,
    }
    text := string(buffer)
    t := template.Must(template.New("email").Parse(text))
    buf := &bytes.Buffer{}
    if err := t.Execute(buf, data); err != nil {
        log.Fatalf("Error while parsing Email Template confirm-account: %s", err)
        return
    }
    s := buf.String()
    go SendMail(SpaceDock.Settings.SupportMail, []string{user.Email}, "Welcome to " + SpaceDock.Settings.SiteName + "!", s, true)
}

func SendReset(user objects.User) {
    buffer,err := ioutil.ReadFile("emails/password-reset")
    if err != nil {
        log.Fatalf("Error while reading Email Template password-reset: %s", err)
        return
    }
    data := map[string]interface{}{
        "SiteName":     SpaceDock.Settings.SiteName,
        "Username": user.Username,
        "Domain": SpaceDock.Settings.Domain,
        "Confirmation": user.PasswordReset,
    }
    text := string(buffer)
    t := template.Must(template.New("email").Parse(text))
    buf := &bytes.Buffer{}
    if err := t.Execute(buf, data); err != nil {
        log.Fatalf("Error while parsing Email Template password-reset: %s", err)
        return
    }
    s := buf.String()
    go SendMail(SpaceDock.Settings.SupportMail, []string{user.Email}, "Reset your password on " + SpaceDock.Settings.SiteName, s, true)
}
