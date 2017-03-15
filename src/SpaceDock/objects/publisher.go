/*
 SpaceDock Backend
 API Backend for the SpaceDock infrastructure to host modfiles for various games

 SpaceDock Backend is licensed under the Terms of the MIT License.
 Copyright (c) 2017 Dorian Stoll (ThomasKerman/TMSP), RockyTV
 */

package objects

import "SpaceDock"

type Publisher struct {
    Model

    Name             string `gorm:"size:1024;unique_index;not null"`
    Description      string `gorm:"size:100000"`
    ShortDescription string `gorm:"size:1000"`
    Games            []Game `json:"-" spacedock:"lock"`
}

func (s *Publisher) AfterFind() {
    if SpaceDock.DBRecursion == 2 {
        return
    }
    isRoot := SpaceDock.DBRecursion == 0
    SpaceDock.DBRecursion += 1
    SpaceDock.Database.Model(s).Related(&(s.Games), "Games")
    if isRoot {
        SpaceDock.DBRecursion = 0
    }
}

func NewPublisher(name string) *Publisher {
    pub := &Publisher{ Name: name }
    pub.Meta = "{}"
    return pub
}