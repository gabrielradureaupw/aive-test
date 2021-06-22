package main

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/gabrielradureaupw/aivetest/pkg/email"
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

func (t TimeSlot) String() string {
	if t.IsZero() {
		return "- busy -" // for display purpose
	}
	return t.Format("2006-01-02T15:04")
}
func (t *TimeSlot) UnmarshalJSON(b []byte) (err error) {
	inputStr := strings.Trim(string(b), "\"")
	t.Time, err = time.Parse("2006-01-02T15:04", inputStr)
	return
}
func (t TimeSlot) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.String())), nil
}
func (t *TimeSlot) Scan(src interface{}) error {
	switch tp := src.(type) {
	case time.Time:
		t.Time = src.(time.Time)
		return nil
	default:
		fmt.Printf("timeslot is used with type %+v\n", tp)
	}
	return fmt.Errorf("Failed to scan timeslot from the database")
}
func (t TimeSlot) Value() (driver.Value, error) {
	return t.Time, nil
}

func (t TimeSlot) Valid() bool {
	isSunday := t.Weekday() == 0
	isTimeOff := t.Hour() < 10 || t.Hour() > 15 || t.Hour() > 12 && t.Hour() < 14
	return !(isSunday || isTimeOff)
}

// ComputeTimeSlots set available time slots MON-SAT 10-12AM to 2PM-4PM
func (c *VaccinationCenter) ComputeTimeSlots(nDays int) {
	set := make(map[string]struct{})
	for _, rdv := range c.Appointments {
		set[rdv.TimeSlot.String()] = struct{}{}
	}
	from := time.Now().Local().Truncate(time.Hour)
	for slot := (TimeSlot{from}); slot.Sub(from) < time.Duration(nDays)*24*time.Hour; slot = (TimeSlot{slot.Add(time.Hour)}) {
		if !slot.Valid() {
			continue
		}
		if _, ok := set[slot.String()]; !ok {
			c.Slots = append(c.Slots, slot)
		} else {
			c.Slots = append(c.Slots, TimeSlot{}) // time slot with appointment for display purpose
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
