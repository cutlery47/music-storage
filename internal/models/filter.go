package models

import "time"

type Filter struct {
	Group          *string
	Song           *string
	ReleasedBefore *time.Time
	ReleasedAfter  *time.Time
}
