package db

import (
	"context"
	"database/sql"
)

type HierarchyUser struct {
	ID        string
	Email     string
	Name      string
	ManagerID sql.NullString
}

// GetManager returns manager of given user.
// If user has no manager (CEO), returns (nil, nil).
func GetManager(ctx context.Context, conn *sql.DB, userID string) (*HierarchyUser, error) {
	// 1) Check user exists
	var one int
	err := conn.QueryRowContext(ctx, `SELECT 1 FROM users WHERE id = $1 LIMIT 1`, userID).Scan(&one)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows // user not found
	}
	if err != nil {
		return nil, err
	}

	// 2) Load manager (can be NULL)
	const q = `
SELECT m.id, m.email, m.name, m.manager_id
FROM users u
JOIN users m ON m.id = u.manager_id
WHERE u.id = $1
`
	var u HierarchyUser
	err = conn.QueryRowContext(ctx, q, userID).Scan(&u.ID, &u.Email, &u.Name, &u.ManagerID)
	if err == sql.ErrNoRows {
		// manager_id is NULL => CEO
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// ListDirectSubordinates returns direct subordinates of managerID.
func ListDirectSubordinates(ctx context.Context, conn *sql.DB, managerID string) ([]HierarchyUser, error) {
	// Check manager exists
	var one int
	err := conn.QueryRowContext(ctx, `SELECT 1 FROM users WHERE id = $1 LIMIT 1`, managerID).Scan(&one)
	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows // user not found
	}
	if err != nil {
		return nil, err
	}

	const q = `
SELECT id, email, name, manager_id
FROM users
WHERE manager_id = $1
ORDER BY name ASC
`
	rows, err := conn.QueryContext(ctx, q, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]HierarchyUser, 0)
	for rows.Next() {
		var u HierarchyUser
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.ManagerID); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// IsInManagerChain checks whether possibleManagerID is in the manager chain of targetUserID.
// True means: possibleManagerID is targetUserID itself OR is above it in hierarchy.
func IsInManagerChain(ctx context.Context, conn *sql.DB, targetUserID, possibleManagerID string) (bool, error) {
	const q = `
WITH RECURSIVE chain AS (
  SELECT id, manager_id
  FROM users
  WHERE id = $1
  UNION ALL
  SELECT u.id, u.manager_id
  FROM users u
  JOIN chain c ON c.manager_id = u.id
)
SELECT 1
FROM chain
WHERE id = $2
LIMIT 1;
`
	var one int
	err := conn.QueryRowContext(ctx, q, targetUserID, possibleManagerID).Scan(&one)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
