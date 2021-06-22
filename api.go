package main

type MakeAppointmentRequest struct {
	Appointment
}

type ConfirmAppointmentRequest struct {
	ID    uint   `form:"id" valid:"required"`
	Email string `form:"email" valid:"email,required"`
}

type ListDailyAppointmentsResponse struct { // Authenticated
	DailyAppointments map[string][]*Appointment `json:"dailyAppointments"`
}
