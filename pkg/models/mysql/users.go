package mysql

import (
	"database/sql"
	"strings"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"watchess.org/watchess/pkg/models"
)

type UserModel struct {
	DB *sql.DB
}

// Insert user (name, email, password, role) to database
func (m *UserModel) Insert(name, email, password string, role models.UserRole) (int, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return 0, err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created, role)
	VALUES(?, ?, ?, UTC_TIMESTAMP(), ?)`

	result, err := m.DB.Exec(stmt, name, email, string(hashedPassword), role.String())

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok {
			if mysqlErr.Number == 1062 {
				if strings.Contains(mysqlErr.Message, "name") {
					return 0, models.ErrDuplicateUsername
				} else {
					return 0, models.ErrDuplicateEmail
				}
			}
		}
	}

	// Check last inserted row
	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), err
}

// Authenticate user (email, password)
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	stmt := "SELECT id, hashed_password FROM users WHERE email = ?"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)
	if err == sql.ErrNoRows {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, models.ErrInvalidCredentials
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

// Get user details from database
func (m *UserModel) Get(id int) (*models.User, error) {
	s := &models.User{}

	stmt := `SELECT id, name, email, created, role FROM users WHERE id = ?`
	var roleStr string
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created, &roleStr)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}

	role, err := models.GetUserRole(roleStr)
	if err != nil {
		return nil, err
	}

	s.Role = *role
	return s, nil
}
