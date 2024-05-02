package needforheat

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// Time is a wrapper for time.Time that can be unmarshalled into unix seconds in JSON.
type Time time.Time

// MarshalJSON marshals the Time to Unix seconds.
func (t Time) MarshalJSON() ([]byte, error) {
	unixSeconds := time.Time(t).Unix()
	return json.Marshal(unixSeconds)
}

// UnmarshalJSON unmarshals Unix seconds to Time.
func (t *Time) UnmarshalJSON(b []byte) error {
	var unixSeconds int64
	err := json.Unmarshal(b, &unixSeconds)
	if err != nil {
		return err
	}

	*t = Time(time.Unix(unixSeconds, 0))
	return nil
}

// Value implements the driver.Valuer interface.
func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan implements the sql.Scanner interface.
func (t *Time) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	if v, ok := value.(time.Time); ok {
		*t = Time(v)
		return nil
	}
	return fmt.Errorf("failed to scan Time field")
}
