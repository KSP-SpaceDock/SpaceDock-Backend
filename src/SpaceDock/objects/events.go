/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "SpaceDock"
)

type DownloadEvent struct {
    Model

    Mod   Mod `json:"-" spacedock:"lock"`
    ModID uint `json:"mod" spacedock:"lock"`
    Version ModVersion `json:"-" spacedock:"lock"`
    VersionID uint `json:"version" spacedock:"lock"`
    Downloads int `json:"downloads" spacedock:"lock"`
}

func (s *DownloadEvent) AfterFind() {
    if SpaceDock.DBRecursion == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.Database.Related(&(s.Version), "Version")
    SpaceDock.DBRecursion -= 1
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
}

func NewDownloadEvent(mod Mod, version ModVersion) *DownloadEvent {
    d := &DownloadEvent{
        Mod: mod,
        ModID: mod.ID,
        Version: version,
        VersionID: version.ID,
        Downloads: 0,
    }
    d.Meta = "{}"
    return d
}

/* ========================================= */

type FollowEvent struct {
    Model

    Mod   Mod `json:"-" spacedock:"lock"`
    ModID uint `json:"mod" spacedock:"lock"`
    Events int `json:"events" spacedock:"lock"`
    Delta int `json:"delta" spacedock:"lock"`
}

func (s *FollowEvent) AfterFind() {
    if SpaceDock.DBRecursion == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.DBRecursion -= 1
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
}

func NewFollowEvent(mod Mod) *FollowEvent {
    f := &FollowEvent{
        Mod: mod,
        ModID: mod.ID,
        Events: 0,
        Delta: 0,
    }
    f.Meta = "{}"
    return f
}

/* ========================================= */

type ReferralEvent struct {
    Model

    Mod   Mod `json:"-" spacedock:"lock"`
    ModID uint `json:"mod" spacedock:"lock"`
    Events int `json:"events" spacedock:"lock"`
    Host string `json:"host" gorm:"size:128" spacedock:"lock"`
}

func (s *ReferralEvent) AfterFind() {
    if SpaceDock.DBRecursion == SpaceDock.DBRecursionMax {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.DBRecursion -= 1
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
}

func NewReferralEvent(mod Mod, host string) *ReferralEvent {
    r := &ReferralEvent{
        Mod: mod,
        ModID: mod.ID,
        Events: 0,
        Host: host,
    }
    r.Meta = "{}"
    return r
}
