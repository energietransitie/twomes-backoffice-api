package needforheat

import (
	"encoding/json"
	"time"
)

// A Time is a wrapper for time.Time that can be unmarshalled into unix seconds in JSON.
type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

// Unmarshal JSON to Time. Unix seconds.
func (t *Time) UnmarshalJSON(b []byte) error {
	unixTime, err := UnmarshalUnix(b)

	*t = Time(unixTime)
	return err
}

func UnmarshalUnix(b []byte) (time.Time, error) {
	var unixSeconds int64
	err := json.Unmarshal(b, &unixSeconds)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(unixSeconds, 0), nil
}
