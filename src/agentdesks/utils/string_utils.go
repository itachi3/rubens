package utils

func Search(format []string, list []string) bool {
	for _, a := range format {
		var success = false
		for _, b := range list {
			if a == b {
				success = true
			}
		}

		if !success {
			return false
		}
	}
	return true
}
