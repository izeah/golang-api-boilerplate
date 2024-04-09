package priority

func PriorityInt(i ...int) int {
	for _, str := range i {
		if str != 0 {
			return str
		}
	}
	return 0
}
