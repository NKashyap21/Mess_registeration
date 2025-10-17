package state

var (
	RegRegistrationOpen = true
	VegRegistrationOpen = true
)

func GetRegistrationStatusReg() bool {
	return RegRegistrationOpen
}

func GetRegistrationStatusVeg() bool {
	return VegRegistrationOpen
}

func ChangeRegistrationStatusReg(value bool) {
	RegRegistrationOpen = value
}

func ChangeRegistrationStatusVeg(value bool) {
	VegRegistrationOpen = value
}
