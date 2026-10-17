package arrays

func In(value string, list []string) bool {
	for _, v := range list {
		if value == v {
			return true
		}
	}
	return false
}
