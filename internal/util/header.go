package util

import "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

var allowedHeaders = map[string]struct{}{
	"x-request-id": {},
}

func IsHeaderAllowed(s string) (string, bool) {
	if _, ok := allowedHeaders[s]; ok {
		return s, true
	}

	return runtime.DefaultHeaderMatcher(s)
}
