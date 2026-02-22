package metadata

import "time"

type Table struct {
	Name            string
	RetentionPeriod *time.Duration
}
