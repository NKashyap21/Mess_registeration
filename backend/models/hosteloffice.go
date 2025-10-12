package models

type UserInfo struct {
	Name   string
	Email  string
	RollNo string
	Mess   int8
}

type EditUserInfo struct {
	RollNo      string `json:"roll_no"`      //To identify the student
	Mess        int8   `json:"mess"`         //The hostel office can change this value
	CanRegister bool   `json:"can_register"` //fasle -> The user has been deactivated.
}
