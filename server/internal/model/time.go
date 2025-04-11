package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// JSONTime is a custom time type that handles JSON marshaling/unmarshaling
type JSONTime time.Time

// MarshalJSON implements the json.Marshaler interface
func (t JSONTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format(time.RFC3339))
	return []byte(stamp), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (t *JSONTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	// Remove quotes
	s = s[1 : len(s)-1]

	// Parse time
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = JSONTime(parsedTime)
	return nil
}

// Value implements the driver.Valuer interface
func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan implements the sql.Scanner interface
func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		*t = JSONTime(v)
		return nil
	default:
		return fmt.Errorf("cannot scan %T into JSONTime", value)
	}
}
