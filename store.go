package main

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AppointmentStore interface {
	ListVaccinationCenters() ([]*VaccinationCenter, error)
	FindVaccinationCenter(req VaccinationCenter) (*VaccinationCenter, error)
	CreateAppointment(rdv *Appointment) error
	FindAppointment(rdv Appointment) (*Appointment, error)
	UpdateAppointment(rdv *Appointment) error
}

type Store struct {
	db *gorm.DB
}

func NewStore() AppointmentStore {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	checkErr(err)
	err = db.AutoMigrate(
		Appointment{},
		VaccinationCenter{},
	)
	checkErr(err)
	for i := 0; i < 3; i++ {
		checkErr(db.Save(&VaccinationCenter{
			City: randomdata.City(),
			Name: randomdata.SillyName(),
		}).Error)
	}
	return Store{db}
}

func (s Store) ListVaccinationCenters() ([]*VaccinationCenter, error) {
	centers := []*VaccinationCenter{}
	return centers, s.db.Model(VaccinationCenter{}).
		Preload("Appointments").
		Find(&centers).Error
}

func (s Store) FindVaccinationCenter(req VaccinationCenter) (found *VaccinationCenter, err error) {
	return found, s.db.Model(&found).Where(req).First(&found).Error
}

func (s Store) CreateAppointment(rdv *Appointment) error {
	return s.db.Debug().Create(rdv).Error
}

func (s Store) FindAppointment(rdv Appointment) (found *Appointment, err error) {
	found = &Appointment{}
	return found, s.db.Model(found).Where(rdv).First(found).Error
}

func (s Store) UpdateAppointment(rdv *Appointment) error {
	if rdv.ID == 0 {
		return fmt.Errorf("can't update appointment without its ID")
	}
	return s.db.Model(rdv).Save(rdv).Error
}
