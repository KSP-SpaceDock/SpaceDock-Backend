/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package SpaceDock

import (
    "flag"
    "github.com/jinzhu/configor"
    "log"
    "os"
)

/*
 All the config variables from either the commandline, or
 a dedicated config file
 */
type SettingsData struct {
    // Whether the app should run in debug mode
    Debug bool

    // The displayed name of this site
    SiteName string

    // The email where users who need help write to
    SupportMail string

    // The actual location of your site
    Protocol string
    Domain   string

    // Set this to false to disable registration
    Registration bool

    // The address to bind to
    Host string
    Port int

    // Details for sending emails
    SmtpHost     string
    SmtpPort     int
    SmtpUser     string
    SmtpPassword string
    SmtpTls      bool

    // Database connection
    Dialect string
    ConnectionData string

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

    // Whether the code should generate a dummy database
    CreateDefaultDatabase bool
}

/*
 The instance of the settings store
 */
var Settings SettingsData

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
    flag.BoolVar(&Settings.Debug, "debug", false, "Whether the app should run in debug mode")
    flag.StringVar(&Settings.SiteName, "sitename", "", "The displayed name of this site")
    flag.StringVar(&Settings.SupportMail, "support-mail", "", "The email where users who need help write to")
    flag.StringVar(&Settings.Protocol, "protocol", "http", "The protocol your site is using (http/https)")
    flag.StringVar(&Settings.Domain, "domain", "localhost:5000", "The actual location of your site")
    flag.BoolVar(&Settings.Registration, "registration", true, "Whether registering new users on the site is allowed")
    flag.StringVar(&Settings.Host, "host", "0.0.0.0", "The IP Address to bind to")
    flag.IntVar(&Settings.Port, "port", 5000, "The port to bind to")
    flag.StringVar(&Settings.SmtpHost, "smtphost", "", "The hostname of your SMTP server (leave empty if you dont want to send emails)")
    flag.IntVar(&Settings.SmtpPort, "smtpport", 0, "The port your SMTP Server listens on")
    flag.StringVar(&Settings.SmtpUser, "smtpuser", "", "The username that should get used to log into your SMTP server")
    flag.StringVar(&Settings.SmtpPassword, "smtppassword", "", "The password of your SMTP User")
    flag.BoolVar(&Settings.SmtpTls, "smtptls", false, "Whether TLS should be used (STARTTLS)")
    flag.StringVar(&Settings.Dialect, "dialect", "", "The SQL dialect used by your database")
    flag.StringVar(&Settings.ConnectionData, "connectiondata", "", "Describes the connection to your SQL Database")
    flag.StringVar(&Settings.Storage, "storage", "", "The directory where all modfiles should get stored")
    flag.StringVar(&Settings.CdnDomain, "cdndomain", "", "Whether a custom CDN should be used instead of the local storage")
    flag.StringVar(&Settings.ThumbnailSize, "thumbnailsize", "", "Thumbnail size in WxH format")
    flag.BoolVar(&Settings.DisableSameOrigin, "disablesameorigin", false, "Enables CORS (Cross Origin Requests)")
    flag.BoolVar(&Settings.CreateDefaultDatabase, "createdefaultdatabase", false, "")

    flag.StringVar(&configFile, "configfile", "", "The path for a dedicated configuration file")
    flag.Parse()
}

/*
 Loads the settings from a dedicated configuration file
 */
func loadFromConfigFile() {
    if configFile != "" {
        log.Printf("* Found configuration file: %s", configFile)
        os.Setenv("CONFIGOR_ENV_PREFIX", "SPACEDOCK")
        err := configor.Load(&Settings, configFile)
        if err != nil {
            log.Fatalf("* Failed to parse configuration file: %s", err)
        }
    }
}
