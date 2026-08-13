package repository

import (
	"database/sql"
	"encoding/csv"
	"io"
	"strings"
	"time"

	"karyawan-app/internal/domain"
)

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) domain.EmployeeRepository {
	return &employeeRepository{db: db}
}

// escapeLike removes wildcard characters for exact match in LIKE
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}

func (r *employeeRepository) FindAll(page, limit int, search, status string) ([]domain.Employee, int, error) {
	// First get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM employees WHERE deleted_at IS NULL`
	var countArgs []interface{}
	
	if search != "" {
		countQuery += ` AND (name LIKE ? OR email LIKE ?)`
		searchParam := "%" + escapeLike(search) + "%"
		countArgs = append(countArgs, searchParam, searchParam)
	}
	
	err := r.db.QueryRow(countQuery, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * limit

	query := `SELECT id, name, email, position, role, phone, alamat, created_at, updated_at FROM employees WHERE deleted_at IS NULL`
	var args []interface{}
	
	if search != "" {
		query += ` AND (name LIKE ? OR email LIKE ?)`
		searchParam := "%" + escapeLike(search) + "%"
		args = append(args, searchParam, searchParam)
	}
	
	query += ` ORDER BY created_at DESC LIMIT ? OFFSET ?`
	args = append(args, limit, offset)
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []domain.Employee
	for rows.Next() {
		var e domain.Employee
		var createdAtStr, updatedAtStr sql.NullString
		if err := rows.Scan(&e.ID, &e.Name, &e.Email, &e.Position, &e.Role, &e.Phone, &e.Alamat, &createdAtStr, &updatedAtStr); err != nil {
			return nil, 0, err
		}
		// Convert string time to time.Time if needed
		e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr.String)
		if updatedAtStr.Valid {
			e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr.String)
		}
		employees = append(employees, e)
	}
	
	if employees == nil {
		employees = []domain.Employee{}
	}
	
	return employees, total, nil
}

func (r *employeeRepository) FindByID(id int) (*domain.Employee, error) {
	query := `SELECT id, name, email, position, role, phone, alamat, created_at, updated_at FROM employees WHERE id = ? AND deleted_at IS NULL`
	row := r.db.QueryRow(query, id)

	var e domain.Employee
	var createdAtStr, updatedAtStr sql.NullString
	err := row.Scan(&e.ID, &e.Name, &e.Email, &e.Position, &e.Role, &e.Phone, &e.Alamat, &createdAtStr, &updatedAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	e.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr.String)
	if updatedAtStr.Valid {
		e.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr.String)
	}
	return &e, nil
}

func (r *employeeRepository) Create(employee *domain.Employee) error {
	query := `INSERT INTO employees (name, email, position, role, phone, alamat) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := r.db.Exec(query, employee.Name, employee.Email, employee.Position, employee.Role, employee.Phone, employee.Alamat)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	employee.ID = int(id)
	return nil
}

func (r *employeeRepository) Update(employee *domain.Employee) error {
	query := `UPDATE employees SET name=?, email=?, position=?, role=?, phone=?, alamat=?, updated_at=NOW() WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, employee.Name, employee.Email, employee.Position, employee.Role, employee.Phone, employee.Alamat, employee.ID)
	return err
}

func (r *employeeRepository) Delete(id int) error {
	// Soft delete - set deleted_at timestamp instead of actually deleting
	query := `UPDATE employees SET deleted_at=NOW() WHERE id=? AND deleted_at IS NULL`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *employeeRepository) ExportCSV(writer io.Writer) error {
	query := "SELECT name, email, position, role, phone, alamat FROM employees WHERE deleted_at IS NULL ORDER BY id ASC"
	rows, err := r.db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	csvWriter := csv.NewWriter(writer)
	// Don't defer csvWriter.Flush(), just call it before returning
	
	if err := csvWriter.Write([]string{"Name", "Email", "Position", "Role", "Phone", "Alamat"}); err != nil {
		return err
	}

	for rows.Next() {
		var name, email, position, role, phone, alamat string
		if err := rows.Scan(&name, &email, &position, &role, &phone, &alamat); err != nil {
			return err
		}
		if err := csvWriter.Write([]string{name, email, position, role, phone, alamat}); err != nil {
			return err
		}
	}
	
	csvWriter.Flush()
	return csvWriter.Error()
}
