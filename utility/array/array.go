package array

import "sort"

func Filter(ary []string, pred func(string) bool) []string {
	res := make([]string, 0)
	for _, s := range ary {
		if pred(s) {
			res = append(res, s)
		}
	}
	return res
}

func Map(ary []string, conv func(string) string) []string {
	res := make([]string, len(ary))
	for i, s := range ary {
		res[i] = conv(s)
	}
	return res
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
