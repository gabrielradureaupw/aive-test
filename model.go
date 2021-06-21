package main

import "time"

type TimeSlot struct {
	time.Time `time_format:"2006-01-02"`
}

type VaccinationCenter struct {
	City         string        `json:"city,omitempty"`
	Name         string        `json:"name,omitempty"`
	Slots        []TimeSlot    `json:"slots,omitempty"`
	Appointments []Appointment `json:"appointments,omitempty"`
}

type Appointment struct {
	VaccinationCenter `json:"vaccination_center,omitempty"`
	TimeSlot          `json:"time_slot,omitempty"`
	Email             string `json:"email,omitempty"`
	Confirmed         *bool  `json:"confirmed,omitempty"`
}
