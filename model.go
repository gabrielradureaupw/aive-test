package main

import "github.com/gabrielradureaupw/aivetest/pkg/timeslot"

type VaccinationCenter struct {
	ID   uint   `json:"-"`
	City string `json:"city" valid:"required"`
	Name string `json:"name" valid:"required"`

	Slots        []timeslot.TimeSlot `json:"slots,omitempty" gorm:"-"`
	Appointments []*Appointment      `json:"appointments,omitempty"`
}

type Appointment struct {
	ID                  uint `json:"-"`
	VaccinationCenterID uint `json:"-" gorm:"unique:uniq_rdv"`
	*VaccinationCenter  `json:"vaccinationCenter" gorm:"-" valid:"required"`
	TimeSlot            timeslot.TimeSlot `json:"timeSlot" gorm:"unique:uniq_rdv" valid:"required"`
	Email               string            `json:"email,omitempty" valid:"email,required"`
	Confirmed           *bool             `json:"confirmed"`
}
