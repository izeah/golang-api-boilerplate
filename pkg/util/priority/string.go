package priority

func PriorityString(s ...string) string {
	for _, str := range s {
		if str != "" {
			return str
		}
	}
	return ""
}
