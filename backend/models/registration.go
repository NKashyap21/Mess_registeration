package models

import "time"

type MessRegistrationDetails struct {
	VegRegistrationStart    time.Time `json:"veg_registration_start,omitempty"`
	VegRegistrationEnd      time.Time `json:"veg_registration_end,omitempty"`
	NormalRegistrationStart time.Time `json:"normal_registration_start,omitempty"`
	NormalRegistrationEnd   time.Time `json:"normal_registration_end,omitempty"`
	MessALDHCapacity        int       `json:"mess_a_ldh_capacity,omitempty"`
	MessAUDHCapacity        int       `json:"mess_a_udh_capacity,omitempty"`
	MessBLDHCapacity        int       `json:"mess_b_ldh_capacity,omitempty"`
	MessBUDHCapacity        int       `json:"mess_b_udh_capacity,omitempty"`
	VegMessCapacity         int       `json:"veg_mess_capacity,omitempty"`
}
