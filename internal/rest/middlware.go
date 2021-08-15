package rest

import "net/http"

func (a *RBACHandler) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("token")
		if err != nil {
			renderErrorResponse(r.Context(), w, "invalid token", err)
			return
		}
		payload, err := a.svc.VerifyToken(c.Value)
		if err != nil {
			renderErrorResponse(r.Context(), w, "invalid token", err)
			return
		}
		r.Header.Set("username", payload.Username)
		// fmt.Println(payload)
		next.ServeHTTP(w, r)
	})
}
