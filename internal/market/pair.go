package market

import "strings"

type Pair struct {
	Base  string
	Quote string
}

func ParsePair(s string) Pair {
	if !strings.Contains(s, "/") {
		return Pair{}
	}

	parts := strings.Split(s, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return Pair{}
	}

	return Pair{
		Base:  parts[0],
		Quote: parts[1],
	}
}
