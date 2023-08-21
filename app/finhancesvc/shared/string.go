package shared

import "time"

func StringToRFC3339(stringTime string) (parsedTime time.Time, err error) {
	return time.Parse(time.RFC3339, stringTime)
}
