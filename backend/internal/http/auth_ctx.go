package http

import "net/http"

type ctxKey string

const authCtxKey ctxKey = "auth"

type AuthCtx struct {
	UserID int    `json:"user_id"`
	OrgID  int    `json:"org_id"`
	Role   string `json:"role"`
}

func mustAuth(r *http.Request) (*AuthCtx, bool) {
	v := r.Context().Value(authCtxKey)
	if v == nil {
		return nil, false
	}
	a, ok := v.(*AuthCtx)
	return a, ok
}
