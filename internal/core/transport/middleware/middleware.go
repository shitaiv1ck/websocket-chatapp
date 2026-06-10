package core_middleware

import (
	"net/http"

	core_logger "github.com/shitaiv1ck/realtime-chat/internal/core/logger"
)

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(h http.Handler, m ...Middleware) http.Handler {
	if len(m) == 0 {
		return h
	}

	for i := len(m) - 1; i >= 0; i-- {
		h = m[i](h)
	}

	return h
}

func CommonMiddleware(h http.Handler, log *core_logger.Logger) http.Handler {
	return ChainMiddleware(
		h,
		RequestID(),
		Logger(log),
		Trace(),
	)
}

func ProtectedMiddleware(h http.Handler, service SessionsService) http.Handler {
	return ChainMiddleware(
		h,
		Authorization(service),
	)
}
