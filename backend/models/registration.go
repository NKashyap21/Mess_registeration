package models

import "time"

type MessRegistrationDetails struct {
	VegRegistrationStart time.Time `json:"veg_registration_start"`
	VegRegistrationEnd   time.Time `json:"veg_registration_end"`
	NormalRegistrationStart time.Time `json:"normal_registration_start"`
	NormalRegistrationEnd   time.Time `json:"normal_registration_end"`
	MessALDHCapacity int       `json:"mess_a_ldh_capacity"`
	MessAUDHCapacity int       `json:"mess_a_udh_capacity"`
	MessBLDHCapacity int       `json:"mess_b_ldh_capacity"`
	MessBUDHCapacity int       `json:"mess_b_udh_capacity"`
	VegMessCapacity  int       `json:"veg_mess_capacity",omitempty`
}