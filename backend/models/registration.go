package models

import "time"

type MessRegistrationDetails struct {
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	MessALDHCapacity int       `json:"mess_a_ldh_capacity"`
	MessAUDHCapacity int       `json:"mess_a_udh_capacity"`
	MessBLDHCapacity int       `json:"mess_b_ldh_capacity"`
	MessBUDHCapacity int       `json:"mess_b_udh_capacity"`
}
