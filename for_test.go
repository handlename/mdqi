package mdqi

import (
	"reflect"
	"sort"
	"strings"
)

func compareAfterTrim(a, b string, cutset string) bool {
	return strings.Trim(a, cutset) == strings.Trim(b, cutset)
}

func sortEqual(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)

	return reflect.DeepEqual(a, b)
}
