package core_request

import "net/http"

func IsMethodSafe(method string) bool {
	if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
		return true
	}

	return false
}
