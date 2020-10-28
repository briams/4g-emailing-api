package utils

import (
	"time"
)

// Now returns date information object in UTC
func Now() map[string]interface{} {
	dt := time.Now().UTC()
	datetime := make(map[string]interface{})
	datetime["date"] = dt.Format("2006-01-02")
	datetime["time"] = dt.Format("2006-01-02 15:04:05")
	datetime["unix"] = dt.Unix()
	return datetime
}
