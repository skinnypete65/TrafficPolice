package repository

import (
	"TrafficPolice/internal/domain"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthRepoPostgres struct {
	conn *pgx.Conn
}

func NewAuthRepoPostgres(conn *pgx.Conn) *AuthRepoPostgres {
	return &AuthRepoPostgres{conn: conn}
}

const checkUserExistsQuery = "SELECT username FROM users WHERE username = $1"

func (r *AuthRepoPostgres) CheckUserExists(username string) error {
	row := r.conn.QueryRow(context.Background(), checkUserExistsQuery, username)

	var userName string
	err := row.Scan(&userName)

	return err
}

const insertUserQuery = `INSERT INTO users (user_id, username, hash_pass, role) 
	VALUES ($1, $2, $3, $4)`

func (r *AuthRepoPostgres) InsertUser(user domain.UserInfo) error {
	_, err := r.conn.Exec(context.Background(), insertUserQuery,
		user.ID.String(),
		user.Username,
		user.Password,
		user.UserRole,
	)

	return err
}

const insertExpertQuery = `INSERT INTO experts (expert_id, is_confirmed, user_id, competence_skill) 
	VALUES ($1, false, $2, 1)`

func (r *AuthRepoPostgres) InsertExpert(expert domain.Expert) error {
	_, err := r.conn.Exec(context.Background(), insertExpertQuery, expert.ID, expert.UserInfo.ID.String())
	return err
}

const insertDirectorQuery = "INSERT INTO directors (director_id, user_id) VALUES ($1, $2)"

func (r *AuthRepoPostgres) InsertDirector(director domain.Director) error {
	_, err := r.conn.Exec(context.Background(), insertDirectorQuery, director.ID.String(), director.User.ID.String())
	return err
}

const signInQuery = `SELECT user_id, hash_pass, role FROM users WHERE username = $1`

func (r *AuthRepoPostgres) SignIn(username string) (domain.UserInfo, error) {
	row := r.conn.QueryRow(context.Background(), signInQuery, username)

	var user domain.UserInfo
	var userID string
	err := row.Scan(&userID, &user.Password, &user.UserRole)
	if err != nil {
		return domain.UserInfo{}, err
	}
	user.ID = uuid.MustParse(userID)

	return user, nil
}

const confirmExpertQuery = "UPDATE experts SET is_confirmed = $1 WHERE expert_id = $2"

func (r *AuthRepoPostgres) ConfirmExpert(data domain.ConfirmExpert) error {
	_, err := r.conn.Exec(context.Background(), confirmExpertQuery, data.IsConfirmed, data.ExpertID)
	return err
}
