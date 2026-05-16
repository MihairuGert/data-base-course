package postgres

import (
	"context"

	"autopark/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Store) FindUserByUsername(ctx context.Context, username string) (domain.User, error) {
	row := s.pool.QueryRow(ctx, `
		SELECT id, username, password_hash, role, employee_id
		FROM app_users
		WHERE username = $1
	`, username)

	var user domain.User
	var employeeID pgtype.Int8
	if err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Role, &employeeID); err != nil {
		return domain.User{}, err
	}
	if employeeID.Valid {
		user.EmployeeID = &employeeID.Int64
	}
	return user, nil
}
