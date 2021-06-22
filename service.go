package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gabrielradureaupw/aivetest/pkg/email"
	"github.com/gabrielradureaupw/aivetest/pkg/timeslot"
)

func ListAvailableSlots(ctx context.Context, st AppointmentStore) (centers []*VaccinationCenter, err error) {
	centers, err = st.ListVaccinationCenters()
	if err != nil {
		return
	}
	const days = 3
	for _, c := range centers {
		c.ComputeTimeSlots(days)
		c.Appointments = nil
	}
	return
}

// ComputeTimeSlots set available time slots MON-SAT 10-12AM to 2PM-4PM
func (c *VaccinationCenter) ComputeTimeSlots(nDays int) {
	set := make(map[string]struct{})
	for _, rdv := range c.Appointments {
		set[rdv.TimeSlot.String()] = struct{}{}
	}
	from := time.Now().Local().Truncate(time.Hour)
	for slot := (timeslot.TimeSlot{Time: from}); slot.Sub(from) < time.Duration(nDays)*24*time.Hour; slot = (timeslot.TimeSlot{Time: slot.Add(time.Hour)}) {
		if !slot.Valid() {
			continue
		}
		if _, ok := set[slot.String()]; !ok {
			c.Slots = append(c.Slots, slot)
		} else {
			c.Slots = append(c.Slots, timeslot.TimeSlot{}) // time slot with appointment for display purpose
		}
	}
}

func ListDailyAppointments(ctx context.Context, st AppointmentStore) (resp ListDailyAppointmentsResponse, err error) {
	centers, err := st.ListVaccinationCenters()
	if err != nil {
		return
	}
	m := make(map[string][]*Appointment)
	for _, c := range centers {
		ref := *c              // vaccination center copy
		ref.Appointments = nil // strip other appointements info
		for _, rdv := range c.Appointments {
			rdv.VaccinationCenter = &ref
			dayStr := rdv.TimeSlot.String()
			m[dayStr] = append(m[dayStr], rdv)
		}
	}
	resp.DailyAppointments = m
	return
}

func MakeAppointment(ctx context.Context, st AppointmentStore, req MakeAppointmentRequest) error {
	if center, err := st.FindVaccinationCenter(VaccinationCenter{
		Name: req.Name,
		City: req.City,
	}); err != nil {
		return err
	} else {
		req.VaccinationCenterID = center.ID
	}
	confirmed := false
	rdv := &Appointment{
		Email:               req.Email,
		VaccinationCenterID: req.VaccinationCenterID,
		TimeSlot:            req.TimeSlot,
		Confirmed:           &confirmed,
	}
	if err := st.CreateAppointment(rdv); err != nil {
		return err
	}
	return email.Send(req.Email, "Confirmez votre rendez-vous", fmt.Sprintf(`
	Blablabla
	Bla blabla
	<a href="http://localhost:8080/appointments/confirm?id=%d&email=%s">Confirmez</a>
	Formule de politesse
	`, rdv.ID, rdv.Email))
}

func ConfirmAppointment(ctx context.Context, st AppointmentStore, req ConfirmAppointmentRequest) error {
	confirmed := false
	rdv, err := st.FindAppointment(Appointment{
		ID:        req.ID,
		Email:     req.Email,
		Confirmed: &confirmed,
	})
	if err != nil {
		return err
	}
	confirmed = true
	rdv.Confirmed = &confirmed
	if err := st.UpdateAppointment(rdv); err != nil {
		return err
	}
	return nil
}
