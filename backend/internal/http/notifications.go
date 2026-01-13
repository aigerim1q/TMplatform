package http

import (
	"net/http"
	"strconv"

	"tmplatform-backend/internal/db"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NotificationsHandler struct {
	DB *pgxpool.Pool
}

func NewNotificationsHandler(conn *pgxpool.Pool) *NotificationsHandler {
	return &NotificationsHandler{DB: conn}
}

func (h *NotificationsHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	auth, ok := mustAuth(r)
	if !ok {
		writeJSON(w, 401, map[string]any{"error": "unauthorized"})
		return
	}
	userID := auth.UserID

	limit := 50
	offset := 0
	if v := r.URL.Query().Get("limit"); v != "" {
		if x, err := strconv.Atoi(v); err == nil {
			limit = x
		}
	}
	if v := r.URL.Query().Get("offset"); v != "" {
		if x, err := strconv.Atoi(v); err == nil {
			offset = x
		}
	}

	items, err := db.ListNotifications(ctx, h.DB, userID, limit, offset)
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}

	writeJSON(w, 200, map[string]any{
		"items": items,
		"meta":  map[string]any{"limit": limit, "offset": offset},
	})
}

func (h *NotificationsHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	auth, ok := mustAuth(r)
	if !ok {
		writeJSON(w, 401, map[string]any{"error": "unauthorized"})
		return
	}
	userID := auth.UserID

	idStr := chi.URLParam(r, "id")
	nid, err := strconv.Atoi(idStr)
	if err != nil || nid <= 0 {
		writeJSON(w, 400, map[string]any{"error": "invalid id"})
		return
	}

	updated, err := db.MarkNotificationRead(ctx, h.DB, userID, nid)
	if err != nil {
		writeJSON(w, 500, map[string]any{"error": "db error"})
		return
	}
	if !updated {
		writeJSON(w, 404, map[string]any{"error": "not found"})
		return
	}

	writeJSON(w, 200, map[string]any{"ok": true})
}
