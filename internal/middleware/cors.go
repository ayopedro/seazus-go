package middleware

import "net/http"

func CORS(trustedOrigins []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Origin")

			origin := r.Header.Get("Origin")

			if origin != "" {
				for _, trustedOrigin := range trustedOrigins {
					if origin == trustedOrigin {
						w.Header().Set("Access-Control-Allow-Origin", origin)

						if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
							w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
							w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

							w.WriteHeader(http.StatusOK)
							return
						}
						break
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
