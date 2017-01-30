package array

func Search(ary []string, t string) int {
	for i, s := range ary {
		if s == t {
			return i
		}
	}
	return -1
}

func IsInclude(s string, ary []string) bool {
	for _, ss := range ary {
		if ss == s {
			return true
		}
	}
	return false
}
