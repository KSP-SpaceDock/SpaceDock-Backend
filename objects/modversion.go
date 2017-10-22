/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "github.com/KSP-SpaceDock/SpaceDock-Backend/app"
    "github.com/KSP-SpaceDock/SpaceDock-Backend/utils"
    "os"
    "path/filepath"
)

type ModVersion struct {
    Model

    ModID           uint `json:"mod" spacedock:"lock"`
    FriendlyVersion string `json:"friendly_version" gorm:"size:64;" spacedock:"lock"`
    Beta            bool `json:"beta"`
    GameVersion     GameVersion `json:"gameversion" gorm:"ForeignKey:GameVersionID" spacedock:"lock;tomap"`
    GameVersionID   uint `json:"gameversion_id" spacedock:"lock"`
    DownloadPath    string `json:"download_path" gorm:"size:512" spacedock:"lock"`
    Changelog       string `json:"changelog" gorm:"size:10000"`
    SortIndex       int `json:"sort_index" spacedock:"lock"`
    FileSize        int64 `json:"file_size" spacedock:"lock"`
}

func (s *ModVersion) AfterFind() {
    app.DBRecursionLock.Lock()
    if _, ok := app.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        app.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if app.DBRecursion[utils.CurrentGoroutineID()] >= app.DBRecursionMax {
        app.DBRecursionLock.Unlock()
        return
    }
    isRoot := app.DBRecursion[utils.CurrentGoroutineID()] == 0
    app.DBRecursion[utils.CurrentGoroutineID()] += 1
    app.DBRecursionLock.Unlock()

    app.Database.Model(s).Related(&(s.GameVersion), "GameVersion")

    app.DBRecursionLock.Lock()
    app.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(app.DBRecursion, utils.CurrentGoroutineID())
    }
    app.DBRecursionLock.Unlock()
}

func NewModVersion(mod Mod, friendly_version string, gameversion GameVersion, download_path string, beta bool) *ModVersion {
    mv := &ModVersion{
        ModID: mod.ID,
        FriendlyVersion: friendly_version,
        Beta: beta,
        GameVersion: gameversion,
        GameVersionID: gameversion.ID,
        DownloadPath: download_path,
        Changelog: "",
        SortIndex: 0,
        FileSize: 0,
    }
    mv.Meta = "{}"
    if mv.DownloadPath != "" {
        f, err := os.Open(filepath.Join(app.Settings.Storage, mv.DownloadPath))
        defer f.Close()
        if err == nil {
            info, err := f.Stat()
            if err == nil {
                mv.FileSize = info.Size()
            }
        }
    }
    return mv
}