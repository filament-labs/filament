package interceptors

import (
	"context"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog/log"
)

func LoggingUnaryHandler() connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			start := time.Now()
			resp, err := next(ctx, req)
			log.Info().
				Str("method", req.Spec().Procedure).
				Dur("Duration", time.Since(start)).
				Err(err).
				Msg("request")

			return resp, err
		}
	}
}
