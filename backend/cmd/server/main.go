package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"

	dbpkg "tmplatform-backend/internal/db"
)

type App struct {
	DB        *pgxpool.Pool
	JWTSecret []byte
	JWTTTL    time.Duration
}

type ctxKey string

const authCtxKey ctxKey = "auth"

type AuthCtx struct {
	UserID int    `json:"user_id"`
	OrgID  int    `json:"org_id"`
	Role   string `json:"role"`
}

func main() {
	_ = godotenv.Load()

	appPort := getenvDefault("APP_PORT", "3001")

	dbURL := buildDBURL()
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	ttlMin, _ := strconv.Atoi(getenvDefault("JWT_TTL_MIN", "60"))
	secret := []byte(getenvDefault("JWT_SECRET", "supersecret_change_me"))

	app := &App{
		DB:        pool,
		JWTSecret: secret,
		JWTTTL:    time.Duration(ttlMin) * time.Minute,
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(corsMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// public auth
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.handleRegister)
		r.Post("/login", app.handleLogin)
	})

	// projects & stages (public with optional auth)
	r.Get("/projects", app.handleProjectsList)
	r.Post("/projects", app.handleProjectsCreate)
	r.Get("/projects/{id}", app.handleProjectGet)
	r.Get("/projects/{id}/stages", app.handleProjectStagesList)
	r.Put("/stages/{id}", app.handleStageUpdate)

	// protected
	r.Group(func(r chi.Router) {
		r.Use(app.authMiddleware)

		r.Get("/me", app.handleMe)

		// notifications (protected)
		r.Get("/notifications", app.handleNotificationsList)
		r.Put("/notifications/{id}/read", app.handleNotificationMarkRead)

		// tasks & files
		r.Get("/stages/{id}/tasks", app.handleStageTasksList)
		r.Get("/tasks", app.handleTasksList)
		r.Post("/tasks", app.handleTaskCreate)
		r.Get("/tasks/{id}", app.handleTaskGet)
		r.Post("/tasks/{id}/files", app.handleTaskFileUpload)

		// stages
		r.Post("/projects/{id}/stages", app.handleProjectStageCreate)
	})

	log.Printf("server started on port %s", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, r))
}

/* =========================
   Handlers: Auth
========================= */

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgID    int    `json:"org_id"`
		Role     string `json:"role"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	body.Role = strings.TrimSpace(strings.ToLower(body.Role))

	if body.Email == "" || len(body.Password) < 8 || body.OrgID <= 0 || !isValidRole(body.Role) {
		writeErr(w, http.StatusBadRequest, "validation failed (email, password>=8, org_id>0, role)")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "hash error")
		return
	}

	var userID int
	err = a.DB.QueryRow(r.Context(),
		`INSERT INTO users (org_id, email, password_hash, role) VALUES ($1,$2,$3,$4) RETURNING id`,
		body.OrgID, body.Email, string(hash), body.Role,
	).Scan(&userID)
	if err != nil {
		writeErr(w, http.StatusConflict, "user already exists or db error")
		return
	}

	token, err := a.issueJWT(AuthCtx{UserID: userID, OrgID: body.OrgID, Role: body.Role})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "token error")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"id":          userID,
		"email":       body.Email,
		"org_id":      body.OrgID,
		"role":        body.Role,
		"accessToken": token,
	})
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if body.Email == "" || body.Password == "" {
		writeErr(w, http.StatusBadRequest, "validation failed")
		return
	}

	var (
		id           int
		orgID        int
		role         string
		passwordHash string
	)

	err := a.DB.QueryRow(r.Context(),
		`SELECT id, org_id, role, password_hash FROM users WHERE email=$1`,
		body.Email,
	).Scan(&id, &orgID, &role, &passwordHash)

	if err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(body.Password)); err != nil {
		writeErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	token, err := a.issueJWT(AuthCtx{UserID: id, OrgID: orgID, Role: role})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "token error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"accessToken": token,
	})
}

func (a *App) handleMe(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	writeJSON(w, http.StatusOK, auth)
}

/* =========================
   Handlers: Notifications
========================= */

func (a *App) handleNotificationsList(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}

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

	items, err := dbpkg.ListNotifications(r.Context(), a.DB, auth.UserID, limit, offset)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"meta":  map[string]any{"limit": limit, "offset": offset},
	})
}

func (a *App) handleNotificationMarkRead(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 {
		writeErr(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	nid, err := strconv.Atoi(idStr)
	if err != nil || nid <= 0 {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	updated, err := dbpkg.MarkNotificationRead(r.Context(), a.DB, auth.UserID, nid)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	if !updated {
		writeErr(w, http.StatusNotFound, "not found")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

/* =========================
   Handlers: Projects & Stages
========================= */

func (a *App) handleProjectsList(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 || auth.OrgID <= 0 {
		devAuth, err := a.ensureDevUser(r.Context())
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "auth error")
			return
		}
		auth = devAuth
	}

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

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	items, err := dbpkg.ListProjects(r.Context(), a.DB, auth.OrgID, limit, offset, sort, order)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	for i := range items {
		assignees, err := dbpkg.ListProjectAssignees(r.Context(), a.DB, items[i].ID)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "db error")
			return
		}
		items[i].Assignees = assignees
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"meta":  map[string]any{"limit": limit, "offset": offset, "sort": sort, "order": order},
	})
}

func (a *App) handleProjectGet(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 || auth.OrgID <= 0 {
		devAuth, err := a.ensureDevUser(r.Context())
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "auth error")
			return
		}
		auth = devAuth
	}

	idStr := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(idStr)
	if err != nil || pid <= 0 {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	item, err := dbpkg.GetProject(r.Context(), a.DB, auth.OrgID, pid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, http.StatusOK, item)
}

func (a *App) handleProjectStagesList(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 || auth.OrgID <= 0 {
		devAuth, err := a.ensureDevUser(r.Context())
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "auth error")
			return
		}
		auth = devAuth
	}

	idStr := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(idStr)
	if err != nil || pid <= 0 {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	_, err = dbpkg.GetProject(r.Context(), a.DB, auth.OrgID, pid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	items, err := dbpkg.ListStages(r.Context(), a.DB, pid, sort, order)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"items": items,
		"meta":  map[string]any{"sort": sort, "order": order},
	})
}

func (a *App) handleProjectStageCreate(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 || auth.OrgID <= 0 {
		devAuth, err := a.ensureDevUser(r.Context())
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "auth error")
			return
		}
		auth = devAuth
	}

	idStr := chi.URLParam(r, "id")
	pid, err := strconv.Atoi(idStr)
	if err != nil || pid <= 0 {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	_, err = dbpkg.GetProject(r.Context(), a.DB, auth.OrgID, pid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	var input struct {
		Title       string  `json:"title"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}

	input.Title = strings.TrimSpace(input.Title)
	if input.Title == "" {
		writeErr(w, http.StatusBadRequest, "title required")
		return
	}

	stage, err := dbpkg.CreateStage(r.Context(), a.DB, pid, input.Title, input.Description)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, http.StatusCreated, stage)
}

func (a *App) handleProjectsCreate(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 || auth.OrgID <= 0 {
		devAuth, err := a.ensureDevUser(r.Context())
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "auth error")
			return
		}
		auth = devAuth
	}

	type payload struct {
		Name            string
		Description     *string
		Status          string
		StartDate       string
		EndDate         string
		Priority        int
		ImageURL        *string
		BudgetAllocated float64
		BudgetSpent     float64
		BudgetCurrency  string
		AssigneeIDs     []int
	}

	var body payload

	if strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid multipart")
			return
		}

		body.Name = strings.TrimSpace(r.FormValue("name"))
		desc := strings.TrimSpace(r.FormValue("description"))
		if desc != "" {
			body.Description = &desc
		}
		body.Status = strings.TrimSpace(strings.ToLower(r.FormValue("status")))
		body.StartDate = strings.TrimSpace(r.FormValue("start_date"))
		body.EndDate = strings.TrimSpace(r.FormValue("end_date"))
		body.Priority = parseIntDefault(r.FormValue("priority"), 0)
		body.BudgetAllocated = parseFloatDefault(r.FormValue("budget_allocated"), 0)
		body.BudgetSpent = parseFloatDefault(r.FormValue("budget_spent"), 0)
		body.BudgetCurrency = strings.TrimSpace(r.FormValue("budget_currency"))

		imageURL := strings.TrimSpace(r.FormValue("image_url"))
		if imageURL != "" {
			body.ImageURL = &imageURL
		}

		file, header, err := r.FormFile("image")
		if err == nil && file != nil {
			defer file.Close()
			data, err := io.ReadAll(file)
			if err != nil {
				writeErr(w, http.StatusBadRequest, "invalid image")
				return
			}
			mimeType := header.Header.Get("Content-Type")
			if mimeType == "" {
				mimeType = http.DetectContentType(data)
			}
			encoded := base64.StdEncoding.EncodeToString(data)
			dataURL := "data:" + mimeType + ";base64," + encoded
			body.ImageURL = &dataURL
		}
	} else {
		type budgetReq struct {
			Allocated float64 `json:"allocated"`
			Spent     float64 `json:"spent"`
			Currency  string  `json:"currency"`
		}

		type assigneeRef struct {
			ID int `json:"id"`
		}

		type req struct {
			Name        string        `json:"name"`
			Description *string       `json:"description"`
			Status      string        `json:"status"`
			StartDate   string        `json:"start_date"`
			EndDate     string        `json:"end_date"`
			Priority    int           `json:"priority"`
			ImageURL    *string       `json:"image_url"`
			Budget      budgetReq     `json:"budget"`
			AssigneeIDs []int         `json:"assignee_ids"`
			Assignees   []assigneeRef `json:"assignees"`
		}

		var input req
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeErr(w, http.StatusBadRequest, "invalid json")
			return
		}
		body.Name = strings.TrimSpace(input.Name)
		body.Description = input.Description
		body.Status = strings.TrimSpace(strings.ToLower(input.Status))
		body.StartDate = input.StartDate
		body.EndDate = input.EndDate
		body.Priority = input.Priority
		body.ImageURL = input.ImageURL
		body.BudgetAllocated = input.Budget.Allocated
		body.BudgetSpent = input.Budget.Spent
		body.BudgetCurrency = input.Budget.Currency
		body.AssigneeIDs = append(body.AssigneeIDs, input.AssigneeIDs...)
		for _, a := range input.Assignees {
			if a.ID > 0 {
				body.AssigneeIDs = append(body.AssigneeIDs, a.ID)
			}
		}
	}

	body.Name = strings.TrimSpace(body.Name)
	body.Status = strings.TrimSpace(strings.ToLower(body.Status))

	if body.Name == "" {
		writeErr(w, http.StatusBadRequest, "name is required")
		return
	}
	if body.Status == "" {
		body.Status = "draft"
	}
	if !isValidProjectStatus(body.Status) {
		writeErr(w, http.StatusBadRequest, "invalid status")
		return
	}

	startDate, err := parseDate(body.StartDate)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid start_date")
		return
	}
	endDate, err := parseDate(body.EndDate)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "invalid end_date")
		return
	}

	item, err := dbpkg.CreateProject(r.Context(), a.DB, dbpkg.CreateProjectInput{
		OrgID:           auth.OrgID,
		OwnerID:         auth.UserID,
		Name:            body.Name,
		Description:     body.Description,
		Status:          body.Status,
		StartDate:       startDate,
		EndDate:         endDate,
		Priority:        body.Priority,
		ImageURL:        body.ImageURL,
		BudgetAllocated: body.BudgetAllocated,
		BudgetSpent:     body.BudgetSpent,
		BudgetCurrency:  body.BudgetCurrency,
		AssigneeIDs:     body.AssigneeIDs,
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, http.StatusCreated, item)
}

func (a *App) handleStageUpdate(w http.ResponseWriter, r *http.Request) {
	auth := mustAuth(r)
	if auth.UserID <= 0 || auth.OrgID <= 0 {
		devAuth, err := a.ensureDevUser(r.Context())
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "auth error")
			return
		}
		auth = devAuth
	}

	idStr := chi.URLParam(r, "id")
	sid, err := strconv.Atoi(idStr)
	if err != nil || sid <= 0 {
		writeErr(w, http.StatusBadRequest, "invalid id")
		return
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	body.Status = strings.TrimSpace(strings.ToLower(body.Status))
	if !isValidStageStatus(body.Status) {
		writeErr(w, http.StatusBadRequest, "invalid status")
		return
	}

	var orgID int
	if err := a.DB.QueryRow(r.Context(), `SELECT p.org_id FROM stages s JOIN projects p ON p.id = s.project_id WHERE s.id = $1`, sid).Scan(&orgID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			writeErr(w, http.StatusNotFound, "not found")
			return
		}
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}
	if orgID != auth.OrgID {
		writeErr(w, http.StatusForbidden, "forbidden")
		return
	}

	item, err := dbpkg.UpdateStageStatus(r.Context(), a.DB, sid, body.Status)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "db error")
		return
	}

	writeJSON(w, http.StatusOK, item)
}

/* =========================
   Middleware + JWT
========================= */

func (a *App) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" || !strings.HasPrefix(h, "Bearer ") {
			writeErr(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		raw := strings.TrimPrefix(h, "Bearer ")

		claims := jwt.MapClaims{}
		tok, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return a.JWTSecret, nil
		})
		if err != nil || tok == nil || !tok.Valid {
			writeErr(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		userID, ok1 := toInt(claims["user_id"])
		orgID, ok2 := toInt(claims["org_id"])
		role, ok3 := claims["role"].(string)
		role = strings.TrimSpace(strings.ToLower(role))

		if !ok1 || !ok2 || !ok3 || userID <= 0 || orgID <= 0 || !isValidRole(role) {
			writeErr(w, http.StatusUnauthorized, "invalid token claims")
			return
		}

		ctx := context.WithValue(r.Context(), authCtxKey, &AuthCtx{
			UserID: userID,
			OrgID:  orgID,
			Role:   role,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *App) issueJWT(auth AuthCtx) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": auth.UserID,
		"org_id":  auth.OrgID,
		"role":    auth.Role,
		"iat":     now.Unix(),
		"exp":     now.Add(a.JWTTTL).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(a.JWTSecret)
}

func mustAuth(r *http.Request) *AuthCtx {
	v := r.Context().Value(authCtxKey)
	if v == nil {
		return &AuthCtx{}
	}
	auth, ok := v.(*AuthCtx)
	if !ok || auth == nil {
		return &AuthCtx{}
	}
	return auth
}

/* =========================
   Helpers
========================= */

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func getenvDefault(k, def string) string {
	v := strings.TrimSpace(os.Getenv(k))
	if v == "" {
		return def
	}
	return v
}

func buildDBURL() string {
	host := getenvDefault("DB_HOST", "localhost")
	port := getenvDefault("DB_PORT", "5432")
	name := getenvDefault("DB_NAME", "tmplatform")
	user := getenvDefault("DB_USER", "tmplatform")
	pass := getenvDefault("DB_PASSWORD", "tmplatform")
	// pgx format
	return "postgres://" + user + ":" + pass + "@" + host + ":" + port + "/" + name
}

func isValidRole(role string) bool {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin", "manager", "employee":
		return true
	default:
		return false
	}
}

func toInt(v any) (int, bool) {
	switch t := v.(type) {
	case float64:
		return int(t), true
	case float32:
		return int(t), true
	case int:
		return t, true
	case int32:
		return int(t), true
	case int64:
		return int(t), true
	case string:
		i, err := strconv.Atoi(t)
		return i, err == nil
	default:
		return 0, false
	}
}

func isValidProjectStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "draft", "active", "paused", "done":
		return true
	default:
		return false
	}
}

func isValidStageStatus(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "todo", "in_progress", "done":
		return true
	default:
		return false
	}
}

func parseDate(v string) (*time.Time, error) {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

func parseIntDefault(v string, def int) int {
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	if i, err := strconv.Atoi(v); err == nil {
		return i
	}
	return def
}

func parseFloatDefault(v string, def float64) float64 {
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	if f, err := strconv.ParseFloat(v, 64); err == nil {
		return f
	}
	return def
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *App) ensureDevUser(ctx context.Context) (*AuthCtx, error) {
	const email = "dev@tmplatform.local"
	var (
		id    int
		orgID int
		role  string
	)

	err := a.DB.QueryRow(ctx, `SELECT id, org_id, role FROM users WHERE email = $1`, email).Scan(&id, &orgID, &role)
	if err == nil {
		return &AuthCtx{UserID: id, OrgID: orgID, Role: role}, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("devpassword123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	role = "admin"
	orgID = 1
	if err := a.DB.QueryRow(ctx,
		`INSERT INTO users (org_id, email, password_hash, role) VALUES ($1,$2,$3,$4) RETURNING id`,
		orgID, email, string(passwordHash), role,
	).Scan(&id); err != nil {
		return nil, err
	}

	return &AuthCtx{UserID: id, OrgID: orgID, Role: role}, nil
}
