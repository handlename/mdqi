package mdqi

import (
	"reflect"
	"regexp"
	"sort"
	"strings"
)

func compareAfterTrim(a, b string) bool {
	re := regexp.MustCompile(" +\n")

	a = re.ReplaceAllString(a, "\n")
	a = strings.Trim(a, " \n")

	b = re.ReplaceAllString(b, "\n")
	b = strings.Trim(b, " \n")

	return a == b
}

func sortEqual(a, b []string) bool {
	sort.Strings(a)
	sort.Strings(b)

	return reflect.DeepEqual(a, b)
}
