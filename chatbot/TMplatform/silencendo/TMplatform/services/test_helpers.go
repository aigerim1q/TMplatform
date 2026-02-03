package services

import (
	"testing"

	"silencendo/db"
)

func newTestService(t *testing.T) *ProjectService {
	t.Helper()
	cfg := db.Config{Driver: "sqlite", DSN: "file:memdb1?mode=memory&cache=shared&_foreign_keys=1"}
	conn, err := db.Connect(cfg)
	if err != nil {
		t.Fatalf("connect test db: %v", err)
	}
	if err := db.Migrate(conn); err != nil {
		t.Fatalf("migrate test db: %v", err)
	}
	return NewProjectService(conn)
}
