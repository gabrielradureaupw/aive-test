package timeslot

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type TimeSlot struct {
	time.Time `time_format:"2006-01-02T15:04"`
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
