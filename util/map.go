package util

import "strconv"


func GetIntFromInterface(value interface{}) (int, bool) {
	switch v := value.(type) {
	case int:
		return v, true
	case float64:
		return int(v), true // Converts float64 to int
	case string:
		// Optional: attempt to parse the string as an integer
		if intValue, err := strconv.Atoi(v); err == nil {
			return intValue, true
		}
	default:
		return 0, false
	}
	return 0, false // Return zero if conversion is not possible
}


