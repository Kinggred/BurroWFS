package utils

func Contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
