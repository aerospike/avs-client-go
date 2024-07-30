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

func (v version) LT(b version) bool {
	strFunc := func(x, y string) bool {
		return x < y
	}
	intFunc := func(x, y int) bool {
		return x < y
	}

	return compare(v, b, strFunc, intFunc)
}

func (v version) GT(b version) bool {
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
		switch a[i].(type) {
		case uint64:
			switch b[i].(type) {
			case uint64:
				if intFunc(int(a[i].(uint64)), int(b[i].(uint64))) {
					return true
				}
			default:
				return false
			}
		case string:
			switch b[i].(type) {
			case string:
				if strFunc(a[i].(string), b[i].(string)) {
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