package models

import "gorm.io/gorm"

type Attendance struct {
	ID        uint    `gorm:"primary key;autoIncrement" json:"id"`
	FullName  *string `json:"fullname"`
	Address   *string `json:"address"`
	Degree    *string `json:"degree"`
	Year      *string `json:"year"`
	Block     *string `json:"block"`
	Subject   *string `json:"subject"`
	Date      *string `json:"date"`
	StartTime *string `json:"startTime"`
	EndTime   *string `json:"endTime"`
}

func MigrateAttendance(db *gorm.DB) error {
	err := db.AutoMigrate(&Attendance{})
	return err
}
