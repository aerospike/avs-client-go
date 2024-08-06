package avs

import (
	"strconv"
	"strings"
)

var minimumSupportedAVSVersion = newVersion("0.9.0")

type version []any

func newVersion(s string) version {
	split := strings.Split(s, ".")
	v := version{}

	for _, token := range split {
		if intVal, err := strconv.ParseUint(token, 10, 64); err == nil {
			v = append(v, intVal)
		} else {
			v = append(v, token)
		}
	}

	return v
}

func (v version) String() string {
	s := ""

	for i, token := range v {
		if i > 0 {
			s += "."
		}

		switch val := token.(type) {
		case uint64:
			s += strconv.FormatUint(val, 10)
		case string:
			s += val
		}
	}

	return s
}

func (v version) lt(b version) bool {
	strFunc := func(x, y string) bool {
		return x < y
	}
	intFunc := func(x, y int) bool {
		return x < y
	}

	return compare(v, b, strFunc, intFunc)
}

func (v version) gt(b version) bool {
	strFunc := func(x, y string) bool {
		return x > y
	}
	intFunc := func(x, y int) bool {
		return x > y
	}

	return compare(v, b, strFunc, intFunc)
}

type compareFunc[T comparable] func(x, y T) bool

func compare(a, b version, strFunc compareFunc[string], intFunc compareFunc[int]) bool {
	sharedLen := min(len(a), len(b))

	for i := 0; i < sharedLen; i++ {
		switch aVal := a[i].(type) {
		case uint64:
			switch bVal := b[i].(type) {
			case uint64:
				if intFunc(int(aVal), int(bVal)) {
					return true
				}
			default:
				return false
			}
		case string:
			switch bVal := b[i].(type) {
			case string:
				if strFunc(aVal, bVal) {
					return true
				}
			default:
				return false
			}
		default:
			return false
		}
	}

	return intFunc(len(a), len(b))
}
