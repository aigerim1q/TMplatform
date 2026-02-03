package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type File struct {
	ID          int       `json:"id"`
	TaskID      int       `json:"task_id"`
	UploaderID  int       `json:"uploader_id"`
	Filename    string    `json:"filename"`
	StoragePath string    `json:"-"` // Internal path
	MimeType    string    `json:"mime_type"`
	SizeBytes   int64     `json:"size_bytes"`
	CreatedAt   time.Time `json:"created_at"`

	UploaderName string `json:"uploader_name,omitempty"`
}

func CreateFile(ctx context.Context, db *pgxpool.Pool, f *File) error {
	q := `
INSERT INTO files (task_id, uploader_id, filename, storage_path, mime_type, size_bytes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, created_at
`
	return db.QueryRow(ctx, q,
		f.TaskID, f.UploaderID, f.Filename, f.StoragePath, f.MimeType, f.SizeBytes,
	).Scan(&f.ID, &f.CreatedAt)
}

func ListFilesByTask(ctx context.Context, db *pgxpool.Pool, taskID int) ([]File, error) {
	q := `
SELECT f.id, f.task_id, f.uploader_id, f.filename, f.mime_type, f.size_bytes, f.created_at, u.name
FROM files f
JOIN users u ON f.uploader_id = u.id
WHERE f.task_id = $1
ORDER BY f.created_at DESC
`
	rows, err := db.Query(ctx, q, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	files := []File{}
	for rows.Next() {
		var f File
		// NOTE: if using pgx v5 with time.Time, it should work directly
		err := rows.Scan(
			&f.ID, &f.TaskID, &f.UploaderID, &f.Filename, &f.MimeType, &f.SizeBytes, &f.CreatedAt,
			&f.UploaderName,
		)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	return files, nil
}

func GetFileByID(ctx context.Context, db *pgxpool.Pool, fileID int) (*File, error) {
	q := `SELECT id, task_id, uploader_id, filename, storage_path, mime_type, size_bytes, created_at FROM files WHERE id = $1`
	var f File
	err := db.QueryRow(ctx, q, fileID).Scan(
		&f.ID, &f.TaskID, &f.UploaderID, &f.Filename, &f.StoragePath, &f.MimeType, &f.SizeBytes, &f.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func DeleteFile(ctx context.Context, db *pgxpool.Pool, fileID int) error {
	_, err := db.Exec(ctx, "DELETE FROM files WHERE id = $1", fileID)
	return err
}
