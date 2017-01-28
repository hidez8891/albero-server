package array

import "strings"

func Search(ary []string, t string) int {
	for i, s := range ary {
		if s == t {
			return i
		}
	}
	return -1
}

func ToJson(ary []string) string {
	store := make([]string, len(ary))
	for i, s := range ary {
		store[i] = "\"" + s + "\""
	}
	return "[" + strings.Join(store, ",") + "]"
}

func IsInclude(s string, ary []string) bool {
	for _, ss := range ary {
		if ss == s {
			return true
		}
	}
	return false
}
