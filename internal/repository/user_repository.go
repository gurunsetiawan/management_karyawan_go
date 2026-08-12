package repository

import (
	"database/sql"
	"errors"

	"karyawan-app/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// CreateUser membuat user baru
func (r *UserRepository) CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, role_id, is_active)
		VALUES (?, ?, ?, ?, ?)
	`

	result, err := r.DB.Exec(query, user.Username, user.Email, user.PasswordHash, user.RoleID, user.IsActive)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)

	// Fetch the created user to get timestamps
	return r.DB.QueryRow("SELECT created_at, updated_at FROM users WHERE id = ?", user.ID).Scan(
		&user.CreatedAt, &user.UpdatedAt,
	)
}

// GetUserByUsername mengambil user berdasarkan username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.role_id, u.is_active, 
		       u.last_login, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.created_at, r.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.username = ?
	`

	user := &models.User{Role: &models.Role{}}
	err := r.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.IsActive,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Description,
		&user.Role.CreatedAt,
		&user.Role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID mengambil user berdasarkan ID
func (r *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.role_id, u.is_active, 
		       u.last_login, u.created_at, u.updated_at,
		       r.id, r.name, r.description, r.created_at, r.updated_at
		FROM users u
		JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`

	user := &models.User{Role: &models.Role{}}
	err := r.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.RoleID,
		&user.IsActive,
		&user.LastLogin,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Role.ID,
		&user.Role.Name,
		&user.Role.Description,
		&user.Role.CreatedAt,
		&user.Role.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// UpdateLastLogin memperbarui waktu terakhir login user
func (r *UserRepository) UpdateLastLogin(userID int) error {
	query := `UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.DB.Exec(query, userID)
	return err
}
