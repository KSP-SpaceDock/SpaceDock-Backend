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
)

type Publisher struct {
    Model

    Name             string `gorm:"size:1024;unique_index;not null" json:"name"`
    Description      string `gorm:"size:100000" json:"description"`
    ShortDescription string `gorm:"size:1000" json:"short_description"`
    Games            []Game `json:"-" spacedock:"lock"`
}

func (s *Publisher) AfterFind() {
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

    SpaceDock.Database.Model(s).Related(&(s.Games), "Games")

    SpaceDock.DBRecursionLock.Lock()
    SpaceDock.DBRecursion[utils.CurrentGoroutineID()] -= 1
    if isRoot {
        delete(SpaceDock.DBRecursion, utils.CurrentGoroutineID())
    }
    SpaceDock.DBRecursionLock.Unlock()
}

func NewPublisher(name string) *Publisher {
    pub := &Publisher{ Name: name }
    pub.Meta = "{}"
    return pub
}