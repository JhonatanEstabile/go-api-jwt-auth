package util

import (
	"net/http"
	"strings"
)

//ExtractToken returns the request authentication token
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}
