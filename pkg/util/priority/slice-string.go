package priority

func PrioritySliceString(s ...[]string) []string {
	var result []string
	for _, str := range s {
		if len(str) != 0 {
			return str
		}
	}
	return result
}
