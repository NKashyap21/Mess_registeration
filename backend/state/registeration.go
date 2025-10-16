package state

var (
	RegRegisterationOpen = false
	VegRegistrationOpen  = false
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
