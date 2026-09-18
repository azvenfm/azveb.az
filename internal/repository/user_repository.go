package repository
import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"autosmm/internal/models"
	_ "github.com/lib/pq"
)
type UserRepository struct { db *sql.DB }
func NewUserRepository(db *sql.DB) *UserRepository { return &UserRepository{db: db} }
func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (id, username, email, password_hash, full_name, role, plan) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash, user.FullName, user.Role, user.Plan)
	return err
}
func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, email, password_hash, full_name, role, plan FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &user.Role, &user.Plan)
	if err != nil {
		if err == sql.ErrNoRows { return nil, errors.New("user not found") }
		return nil, err
	}
	return user, nil
}
