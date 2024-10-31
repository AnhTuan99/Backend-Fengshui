package mid

import (
	"context"
	"net/http"
	"runtime/debug"

	"fengshui.com/back-fengshui/app/sdk/errs"
	"fengshui.com/back-fengshui/app/sdk/metrics"
	"fengshui.com/back-fengshui/foundation/web"
)

// Panics recovers from panics and converts the panic to an error so it is
// reported in Metrics and handled in Errors.
func Panics() web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) (resp web.Encoder) {

			// Defer a function to recover from a panic and set the err return
			// variable after the fact.
			defer func() {
				if rec := recover(); rec != nil {
					trace := debug.Stack()
					resp = errs.Newf(errs.InternalOnlyLog, "PANIC [%v] TRACE[%s]", rec, string(trace))

					metrics.AddPanics(ctx)
				}
			}()

			return next(ctx, r)
		}

		return h
	}

	return m
}