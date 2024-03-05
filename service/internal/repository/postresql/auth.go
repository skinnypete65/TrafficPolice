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

	err := row.Scan()
	return err
}

const insertUserQuery = `INSERT INTO users (user_id, username, hash_pass) 
	VALUES ($1, $2, $3)`

func (r *AuthRepoPostgres) InsertUser(user domain.User) error {
	_, err := r.conn.Exec(context.Background(), insertUserQuery,
		user.ID.String(),
		user.Username,
		user.Password,
	)

	return err
}

const insertExpertQuery = "INSERT INTO experts (expert_id, is_confirmed, user_id) VALUES ($1, false, $2)"

func (r *AuthRepoPostgres) InsertExpert(expert domain.Expert) error {
	_, err := r.conn.Exec(context.Background(), insertExpertQuery, expert.ID.String(), expert.User.ID.String())
	return err
}

const insertDirectorQuery = "INSERT INTO directors (director_id, user_id) VALUES ($1, $2)"

func (r *AuthRepoPostgres) InsertDirector(director domain.Director) error {
	_, err := r.conn.Exec(context.Background(), insertDirectorQuery, director.ID.String(), director.User.ID.String())
	return err
}

const signInQuery = `SELECT user_id, hash_pass FROM users WHERE username = $1`

func (r *AuthRepoPostgres) SignIn(username string) (domain.User, error) {
	row := r.conn.QueryRow(context.Background(), signInQuery, username)

	var userID, hashPass string
	err := row.Scan(&userID, &hashPass)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:       uuid.MustParse(userID),
		Password: hashPass,
	}, nil
}
