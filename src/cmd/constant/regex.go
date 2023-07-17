package constant

import "regexp"

var NumberRegex = regexp.MustCompile("(\\+[0-9]{0,2}|0|62)[0-9]+")
