package state

var (
	RegRegisterationOpen = true
	VegRegistrationOpen  = true
)

func GetRegistrationStatusReg() bool {
	return RegRegisterationOpen
}

func GetRegistrationStatusVeg() bool {
	return VegRegistrationOpen
}

func ChangeRegistrationStatusReg(value bool) {
	RegRegisterationOpen = value
}

func ChangeRegistrationStatusVeg(value bool) {
	VegRegistrationOpen = value
}
