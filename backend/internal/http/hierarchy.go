package http

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"tmplatform-backend/internal/db"

	"github.com/go-chi/chi/v5"
)

// HierarchyHandler holds db connection
type HierarchyHandler struct {
	DB *sql.DB
}

func NewHierarchyHandler(conn *sql.DB) *HierarchyHandler {
	return &HierarchyHandler{DB: conn}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// TODO: Replace these with your real auth extraction logic.
// For now it returns empty string, so RBAC will not work until you wire it.
func getAuthUserID(r *http.Request) string {
	// пример: return r.Context().Value("user_id").(string)
	return ""
}

// TODO: if you have roles, check admin/hr here.
func isAdminOrHR(r *http.Request) bool {
	return false
}

// GetManager: GET /users/{id}/manager
func (h *HierarchyHandler) GetManager(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// IMPORTANT: replace ParamUserID with your router-specific param getter (see Step 4)
	userID := ParamUserID(r)

	authID := getAuthUserID(r)

	// RBAC: allow if self OR admin/hr OR auth in chain of user
	if authID != "" && authID != userID && !isAdminOrHR(r) {
		ok, err := db.IsInManagerChain(ctx, h.DB, userID, authID)
		if err != nil {
			writeJSON(w, 500, map[string]any{"error": "db error"})
			return
		}
		if !ok {
			writeJSON(w, 403, map[string]any{"error": "forbidden"})
			return
		}
	}

	manager, err := db.GetManager(ctx, h.DB, userID)
	if err == sql.ErrNoRows {
		writeJSON(w, 404, map[string]any{"error": "user not found"})
		return
	}
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	if manager == nil {
		// CEO
		writeJSON(w, 200, map[string]any{"manager": nil})
		return
	}

	writeJSON(w, 200, map[string]any{
		"manager": map[string]any{
			"id":    manager.ID,
			"email": manager.Email,
			"name":  manager.Name,
		},
	})
}

// ListSubordinates: GET /users/{id}/subordinates?direct_only=true
func (h *HierarchyHandler) ListSubordinates(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := ParamUserID(r)

	authID := getAuthUserID(r)

	// RBAC: allow if self OR admin/hr OR auth in chain of user
	if authID != "" && authID != userID && !isAdminOrHR(r) {
		ok, err := db.IsInManagerChain(ctx, h.DB, userID, authID)
		if err != nil {
			writeJSON(w, 500, map[string]any{"error": "db error"})
			return
		}
		if !ok {
			writeJSON(w, 403, map[string]any{"error": "forbidden"})
			return
		}
	}

	// query param: direct_only (default true)
	directOnly := true
	if v := r.URL.Query().Get("direct_only"); v != "" {
		b, _ := strconv.ParseBool(v)
		directOnly = b
	}

	// MVP: реализуем только direct_only=true (прямые подчиненные).
	// Если direct_only=false — пока вернем 400, чтобы было честно.
	if !directOnly {
		writeJSON(w, 400, map[string]any{"error": "direct_only=false not implemented yet"})
		return
	}

	items, err := db.ListDirectSubordinates(ctx, h.DB, userID)
	if err == sql.ErrNoRows {
		writeJSON(w, 404, map[string]any{"error": "user not found"})
		return
	}
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}

	resp := make([]map[string]any, 0, len(items))
	for _, u := range items {
		var mid any = nil
		if u.ManagerID.Valid {
			mid = u.ManagerID.String
		}
		resp = append(resp, map[string]any{
			"id":         u.ID,
			"email":      u.Email,
			"name":       u.Name,
			"manager_id": mid,
		})
	}

	writeJSON(w, 200, map[string]any{"items": resp})
}

/*
ParamUserID is router-dependent.
In Step 4 you will edit this function depending on whether you use chi/gin/echo/etc.
*/

func ParamUserID(r *http.Request) string {
	return chi.URLParam(r, "id")
}
