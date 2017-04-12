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

func IsIncludeFunc(s string, ary []string, pred func(string, string) bool) bool {
	for _, ss := range ary {
		if pred(ss, s) {
			return true
		}
	}
	return false
}
