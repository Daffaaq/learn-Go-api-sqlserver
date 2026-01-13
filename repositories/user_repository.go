// ---------------- repositories/user_repository.go ----------------
package repositories

import (
	"context"
	"database/sql"
	"go-sqlserver-api/models"
)

type UserStore interface {
	GetUsersWithPagination(ctx context.Context, limit, offset int) ([]models.User, error)
	CountUsers(ctx context.Context) (int, error)
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	CreateUser(ctx context.Context, u models.User) error
	UpdateUser(ctx context.Context, u models.User) error
	DeleteUser(ctx context.Context, id int) error
	IsEmailExists(ctx context.Context, email string) (bool, error)
}

type UserRepository struct {
	DB *sql.DB
}

// ---------- PAGINATION ----------
func (r *UserRepository) GetUsersWithPagination(ctx context.Context, limit, offset int) ([]models.User, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT id, name, email FROM users ORDER BY id OFFSET @Offset ROWS FETCH NEXT @Limit ROWS ONLY`,
		sql.Named("Offset", offset),
		sql.Named("Limit", limit),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) CountUsers(ctx context.Context) (int, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM users")
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// ---------- EMAIL CHECK ----------
func (r *UserRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT 1 FROM users WHERE email=@Email", sql.Named("Email", email))
	var tmp int
	err := row.Scan(&tmp)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// ---------- CRUD ----------
func (r *UserRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	row := r.DB.QueryRowContext(ctx, "SELECT id, name, email FROM users WHERE id=@ID", sql.Named("ID", id))
	var u models.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, u models.User) error {
	_, err := r.DB.ExecContext(ctx, "INSERT INTO users (name, email) VALUES (@Name, @Email)",
		sql.Named("Name", u.Name),
		sql.Named("Email", u.Email),
	)
	return err
}

func (r *UserRepository) UpdateUser(ctx context.Context, u models.User) error {
	res, err := r.DB.ExecContext(ctx, "UPDATE users SET name=@Name, email=@Email WHERE id=@ID",
		sql.Named("Name", u.Name),
		sql.Named("Email", u.Email),
		sql.Named("ID", u.ID),
	)
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	res, err := r.DB.ExecContext(ctx, "DELETE FROM users WHERE id=@ID", sql.Named("ID", id))
	if err != nil {
		return err
	}
	affected, _ := res.RowsAffected()
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
