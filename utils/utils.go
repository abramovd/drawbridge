package utils

import "strings"

// SingleJoiningSlash returns the concatenation of the
// parameters with a single "/" in between
func SingleJoiningSlash(a, b string) string {
	aSlash := strings.HasSuffix(a, "/")
	bSlash := strings.HasPrefix(b, "/")
	switch {
	case aSlash && bSlash:
		return a + b[1:]
	case !aSlash && !bSlash:
		return a + "/" + b
	}
	return a + b
}

// AddSlashes returns the string with a single slash
// on each side
func AddSlashes(s string) string {
	preSlash := strings.HasPrefix(s, "/")
	postSlash := strings.HasSuffix(s, "/")
	switch {
	case preSlash && postSlash:
		return s
	case !preSlash && postSlash:
		return "/" + s
	case preSlash && !postSlash:
		return s + "/"
	}
	return "/" + s + "/"
}

// StringInSlice checks if value is a valid choice
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
