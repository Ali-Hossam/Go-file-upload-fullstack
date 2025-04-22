package model

import (
	"github.com/google/uuid"
)

type Student struct {
	Student_id   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Student_name string    `gorm:"index;not null"`
	Subject      string    `gorm:"index;not null"`
	Grade        uint      `gorm:"not null"`
}

type StudentTest struct {
	Student_id   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Student_name string    `gorm:"index;not null"`
	Subject      string    `gorm:"index;not null"`
	Grade        uint      `gorm:"not null"`
}
