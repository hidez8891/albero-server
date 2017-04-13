package array

import "sort"

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

func Uniq(ary []string) []string {
	tmp := make([]string, len(ary))
	copy(tmp, ary)
	sort.Strings(tmp)

	if len(tmp) == 0 {
		return tmp
	}

	pre := ""
	res := make([]string, 0)
	for _, s := range tmp[:len(tmp)-1] {
		if s != pre {
			res = append(res, s)
			pre = s
		}
	}

	if tmp[len(tmp)-1] != pre {
		res = append(res, tmp[len(tmp)-1])
	}
	return res
}
