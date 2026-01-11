package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	ttlMin, _ := strconv.Atoi(getenvDefault("JWT_TTL_MIN", "60"))
	secret := []byte(getenvDefault("JWT_SECRET", "supersecret_change_me"))

	app := &App{
		DB:        db,
		JWTSecret: secret,
		JWTTTL:    time.Duration(ttlMin) * time.Minute,
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", app.handleRegister)
		r.Post("/login", app.handleLogin)

		// Protected:
		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware)
			r.Get("/me", app.handleMe)
			r.Post("/logout", app.handleLogout)
			r.Post("/refresh", app.handleRefresh)
		})
	})

	log.Printf("server started on port %s", appPort)
	log.Fatal(http.ListenAndServe(":"+appPort, r))
}

func (a *App) handleRegister(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		OrgID    int    `json:"org_id"`
		Role     string `json:"role"` // allowed at register; later role must come only from JWT
	}
	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid json")
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
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
		// likely unique violation
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

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	// stateless logout: client deletes token
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) handleRefresh(w http.ResponseWriter, r *http.Request) {
	// simplest refresh: if current token valid, issue new one
	auth := mustAuth(r)
	token, err := a.issueJWT(*auth)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "token error")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"accessToken": token})
}

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
	return v.(*AuthCtx)
}

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
	switch role {
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
