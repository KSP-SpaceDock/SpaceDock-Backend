/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package SpaceDock

import (
    "github.com/jinzhu/configor"
    "log"
    "os"
)

/*
 All the config variables from the config file
 */
type SettingsData struct {
    // Whether the app should run in debug mode
    Debug bool

    // The displayed name of this site
    SiteName string `yaml:"site-name" json:"site-name"`

    // The email where users who need help write to
    SupportMail string `yaml:"support-mail" json:"support-mail"`

    // The actual location of your site
    Protocol string
    Domain   string

    // Set this to false to disable registration
    Registration bool

    // The address to bind to
    Host string
    Port int

    // Details for sending emails
    SmtpHost     string `yaml:"smtp-host" json:"smtp-host"`
    SmtpPort     int `yaml:"smtp-port" json:"smtp-port"`
    SmtpUser     string `yaml:"smtp-user" json:"smtp-user"`
    SmtpPassword string `yaml:"smtp-password" json:"smtp-password"`
    SmtpTls      bool `yaml:"smtp-tls" json:"smtp-tls"`

    // Database connection
    Dialect        string
    ConnectionData string `yaml:"connection-data" json:"connection-data"`

    // The directory where files are stored
    Storage string

    // Domain for a storage CDN
    CdnDomain string `yaml:"cdn-domain" json:"cdn-domain"`

    // Thumbnail size in WxH format
    ThumbnailSize string `yaml:"thumbnail-size" json:"thumbnail-size"`

    // Mod URL expression, used for sending emails containing links to the frontend
    // ModUrl string

    // Whether CORS should be enabled
    DisableSameOrigin bool `yaml:"disable-same-origin" json:"disable-same-origin"`

    // How many requests can be made in a defined time span
    RequestLimit string `yaml:"request-limit" json:"request-limit"`

    // Support for X-Accel
    UseXAccel string `yaml:"use-x-accel" json:"use-x-accel"`

    // The default mod url format
    ModUrl string `yaml:"mod-url" json:"mod-url"`

    // Whether to use a memory based store, or redis
    StoreType string `yaml:"store-type" json:"mod-url"`

    // The connection settings for a redis server
    RedisConnection string `yaml:"redis-connection" json:"redis-connection"`

    // How long should a response get cached
    CacheTimeout int `yaml:"cache-timeout" json:"cache-timeout"`
}

/*
 The instance of the settings store
 */
var Settings SettingsData

/*
 Loads the settings from the config file
 */
func LoadSettings() {
    LoadFromConfigFile(&Settings, "config.yml")
}

/*
 Loads the settings from a configuration file
 */
func LoadFromConfigFile(data interface{}, configFile string) {
    log.Printf("* Loading configuration file: config/%s", configFile)
    os.Setenv("CONFIGOR_ENV_PREFIX", "SPACEDOCK")
    err := configor.Load(data, "config/" + configFile)
    if err != nil {
        log.Fatalf("* Failed to parse configuration file: %s", err)
    }
}
