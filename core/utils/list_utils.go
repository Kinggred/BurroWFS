package utils

// Contains checks if a string is present in a list of strings.
// Returns true if found, false otherwise.
func Contains(list []string, str string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
