/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
*/

package SpaceDock

import (
    "flag"
    "log"
    "gopkg.in/gcfg.v1"
)

/*
 All the config variables from either the commandline, or
 a dedicated config file
 */
type SettingsData struct {
    // The displayed name of this site
    SiteName string

    // The email where users who need help write to
    SupportMail string

    // The actual location of your site
    Protocol string
    Domain string

    // Set this to false to disable registration
    Registration bool

    // The address to bind to
    Host string
    Port int

    // Details for sending emails
    SmtpHost string
    SmtpPort int
    SmtpUser string
    SmtpPassword string
    SmtpTls bool

    // Database connection
    ConnectionString string

    // The directory where files are stored
    Storage string

    // Domain for a storage CDN
    CdnDomain string

    // Thumbnail size in WxH format
    ThumbnailSize string

    // Mod URL expression, used for sending emails containing links to the frontend
    // ModUrl string

    // Whether CORS should be enabled
    DisableSameOrigin bool
}

/*
 The instance of the settings store
 */
var settings SettingsData

/*
 The path of the configuration file (if it exists)
 */
var configFile string

/*
 Loads the settings from commandline and from a config file
 */
func LoadSettings() {
    loadFromCommandLine()
    loadFromConfigFile()
}

/*
 Loads the settings from commandline parameters
 */
func loadFromCommandLine() {
    flag.StringVar(&settings.SiteName, "sitename", "", "The displayed name of this site")
    flag.StringVar(&settings.SupportMail, "support-mail", "", "The email where users who need help write to")
    flag.StringVar(&settings.Protocol, "protocol", "http", "The protocol your site is using (http/https)")
    flag.StringVar(&settings.Domain, "domain", "localhost:5000", "The actual location of your site")
    flag.BoolVar(&settings.Registration, "registration", true, "Whether registering new users on the site is allowed")
    flag.StringVar(&settings.Host, "host", "0.0.0.0", "The IP Address to bind to")
    flag.IntVar(&settings.Port, "port", 5000, "The port to bind to")
    flag.StringVar(&settings.SmtpHost, "smtp-host", "", "The hostname of your SMTP server (leave empty if you dont want to send emails)")
    flag.IntVar(&settings.SmtpPort, "smtp-port", 0, "The port your SMTP Server listens on")
    flag.StringVar(&settings.SmtpUser, "smtp-user", "", "The username that should get used to log into your SMTP server")
    flag.StringVar(&settings.SmtpPassword, "smtp-password", "", "The password of your SMTP User")
    flag.BoolVar(&settings.SmtpTls, "smtp-tls", false, "Whether TLS should be used (STARTTLS)")
    flag.StringVar(&settings.ConnectionString, "connection-string", "", "Describes the connection to your SQL Database")
    flag.StringVar(&settings.Storage, "storage", "", "The directory where all modfiles should get stored")
    flag.StringVar(&settings.CdnDomain, "cdn-domain", "", "Whether a custom CDN should be used instead of the local storage")
    flag.StringVar(&settings.ThumbnailSize, "thumbnail-size", "", "Thumbnail size in WxH format")
    flag.BoolVar(&settings.DisableSameOrigin, "disable-same-origin", false, "Enables CORS (Cross Origin Requests)")

    flag.StringVar(&configFile, "config-file", "", "The path for a dedicated configuration file")
    flag.Parse()
}

/*
 Loads the settings from a dedicated configuration file
 */
func loadFromConfigFile() {
    if configFile != "" {
        log.Printf("* Found loading configuration file: %s", configFile)
        err := gcfg.ReadFileInto(&settings, configFile)
        if err != nil {
            log.Fatalf("* Failed to parse configuration file: %s", err)
        }
    }
}
