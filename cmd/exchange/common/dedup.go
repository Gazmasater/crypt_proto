package common

import "sort"

func CanonicalKey(a, b, c string) string {
	arr := []string{a, b, c}
	sort.Strings(arr)
	return arr[0] + "|" + arr[1] + "|" + arr[2]
}
