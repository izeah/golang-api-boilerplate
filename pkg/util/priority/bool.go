package priority

func PriorityBool(i ...bool) bool {
	for _, str := range i {
		if str {
			return str
		}
	}
	return false
}
