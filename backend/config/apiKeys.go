package config

import "os"

type APIKeys struct {
	MessA string
	MessB string
}

var apiKeys = APIKeys{
	MessA: os.Getenv("MESS_A_API_KEY"),
	MessB: os.Getenv("MESS_B_API_KEY"),
}

func GetAPIKeys() *APIKeys {
	return &apiKeys
}
