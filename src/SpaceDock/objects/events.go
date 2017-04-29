/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (StollD), RockyTV
 */

package objects

import (
    "SpaceDock"
    "SpaceDock/utils"
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
    SpaceDock.DBRecursionLock.Lock()
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] >= SpaceDock.DBRecursionMax {
        SpaceDock.DBRecursionLock.Unlock()
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.DBRecursionLock.Unlock()

    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")
    SpaceDock.Database.Related(&(s.Version), "Version")

    SpaceDock.DBRecursionLock.Lock()
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
    SpaceDock.DBRecursionLock.Unlock()
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
    SpaceDock.DBRecursionLock.Lock()
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] >= SpaceDock.DBRecursionMax {
        SpaceDock.DBRecursionLock.Unlock()
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.DBRecursionLock.Unlock()

    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")

    SpaceDock.DBRecursionLock.Lock()
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
    SpaceDock.DBRecursionLock.Unlock()
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
    SpaceDock.DBRecursionLock.Lock()
    if _, ok := SpaceDock.DBRecursion[utils.CurrentGoroutineID()]; !ok {
        SpaceDock.DBRecursion[utils.CurrentGoroutineID()] = 0
    }
    if SpaceDock.DBRecursion[utils.CurrentGoroutineID()] >= SpaceDock.DBRecursionMax {
        SpaceDock.DBRecursionLock.Unlock()
        return
    }
    isRoot := SpaceDock.DBRecursion[utils.CurrentGoroutineID()] == 0
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] += 1
    SpaceDock.DBRecursionLock.Unlock()

    SpaceDock.Database.Model(s).Related(&(s.Mod), "Mod")

    SpaceDock.DBRecursionLock.Lock()
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
    SpaceDock.DBRecursionLock.Unlock()
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
