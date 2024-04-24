package needforheat

import (
	"encoding/json"
	"time"
)

// A Time is a wrapper for time.Time that can be unmarshalled from ISO8601 or unix seconds in JSON.
type Time time.Time

func (t Time) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

// Unmarshal JSON to Time. Can be ISO8601 or unix seconds.
func (t *Time) UnmarshalJSON(b []byte) error {
	isoTime, err := UnmarshalISO8601(b)
	if err == nil {
		*t = Time(isoTime)
		return nil
	}

	unixTime, err := UnmarshalUnix(b)

	*t = Time(unixTime)
	return err
}

func UnmarshalISO8601(b []byte) (time.Time, error) {
	var time time.Time
	err := json.Unmarshal(b, &time)
	return time, err
}

func UnmarshalUnix(b []byte) (time.Time, error) {
	var unixSeconds int64
	err := json.Unmarshal(b, &unixSeconds)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(unixSeconds, 0), nil
}
