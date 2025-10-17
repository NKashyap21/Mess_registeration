package models

type MessRegistrationDetails struct {
	VegRegistrationOpen    bool `json:"veg_registration_open,omitempty"`
	NormalRegistrationOpen bool `json:"normal_registration_open,omitempty"`
	MessALDHCapacity       int  `json:"mess_a_ldh_capacity,omitempty"`
	MessAUDHCapacity       int  `json:"mess_a_udh_capacity,omitempty"`
	MessBLDHCapacity       int  `json:"mess_b_ldh_capacity,omitempty"`
	MessBUDHCapacity       int  `json:"mess_b_udh_capacity,omitempty"`
	VegMessCapacity        int  `json:"veg_mess_capacity,omitempty"`
}
