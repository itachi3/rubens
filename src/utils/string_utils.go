package utils

func Search(format string, list []string) bool {
	for _, b := range list {
		if format == b {
			return true
		}
	}
	return false
}
