package main

import "time"

type ListAvailableSlotsRequest struct {
	VaccinationCenter `json:"vaccination_center,omitempty"`
	Day               time.Time `json:"day,omitempty"`
}

type ListAvailableSlotsResponse struct {
}

type MakeAppointmentRequest struct {
	Appointment
}

type ConfirmAppointmentRequest struct {
	Appointment
}

type ListVaccinationCenterDailyAppointmentsRequest struct { // Authenticated
	VaccinationCenter
}
