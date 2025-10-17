package db

import "time"

type MessRegistrationDetails struct {
	VegRegistrationStart    time.Time `json:"veg_registration_start" time_format:"2006-01-02T15:04:05" time_utc:"IST"`
	VegRegistrationEnd      time.Time `json:"veg_registration_end" time_format:"2006-01-02T15:04:05" time_utc:"IST"`
	NormalRegistrationStart time.Time `json:"normal_registration_start" time_format:"2006-01-02T15:04:05" time_utc:"IST"`
	NormalRegistrationEnd   time.Time `json:"normal_registration_end" time_format:"2006-01-02T15:04:05" time_utc:"IST"`
	MessALDHCapacity        int       `json:"mess_a_ldh_capacity"`
	MessAUDHCapacity        int       `json:"mess_a_udh_capacity"`
	MessBLDHCapacity        int       `json:"mess_b_ldh_capacity"`
	MessBUDHCapacity        int       `json:"mess_b_udh_capacity"`
	VegMessCapacity         int       `json:"veg_mess_capacity"`
}