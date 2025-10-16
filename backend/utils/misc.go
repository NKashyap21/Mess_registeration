package utils


func Contains(slice []int8, value int8) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}