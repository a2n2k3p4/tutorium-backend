package models

import (
	"gorm.io/gorm"
)

type Review struct {
	gorm.Model
	LearnerID uint `gorm:"not null;uniqueIndex:idx_learner_class"`
	ClassID   uint `gorm:"not null;uniqueIndex:idx_learner_class"`
	Rating    int  `gorm:"check:rating >= 1 AND rating <= 5"`
	Comment   string

	Learner Learner `gorm:"foreignKey:LearnerID;references:ID;constraint:OnDelete:CASCADE"`
	Class   Class   `gorm:"foreignKey:ClassID;references:ID;constraint:OnDelete:CASCADE"`
}
