package model

type Student struct {
	Student_id   uint   `gorm:"primaryKey;autoIncrement"`
	Student_name string `gorm:"index"`
	Subject      string `gorm:"index"`
	Grade        uint
}
