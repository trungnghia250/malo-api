package utils

func IsStringContains(s []string, value string) bool {
	for _, v := range s {
		if v == value {
			return true
		}
	}

	return false
}
