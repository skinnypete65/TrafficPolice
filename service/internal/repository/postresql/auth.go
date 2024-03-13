package repository

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/repository"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type authRepoPostgres struct {
	conn *pgx.Conn
}

func NewAuthRepoPostgres(conn *pgx.Conn) repository.AuthRepo {
	return &authRepoPostgres{conn: conn}
}

const checkUserExistsQuery = "SELECT username FROM users WHERE username = $1"

func (r *authRepoPostgres) CheckUserExists(username string) bool {
	row := r.conn.QueryRow(context.Background(), checkUserExistsQuery, username)

	var userName string
	err := row.Scan(&userName)

	return err == nil
}

const insertUserQuery = `INSERT INTO users (user_id, username, hash_pass, role) 
	VALUES ($1, $2, $3, $4)`

func (r *authRepoPostgres) InsertUser(user domain.UserInfo) error {
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

func (r *authRepoPostgres) InsertExpert(expert domain.Expert) error {
	_, err := r.conn.Exec(context.Background(), insertExpertQuery, expert.ID, expert.UserInfo.ID.String())
	return err
}

const insertDirectorQuery = "INSERT INTO directors (director_id, user_id) VALUES ($1, $2)"

func (r *authRepoPostgres) InsertDirector(director domain.Director) error {
	_, err := r.conn.Exec(context.Background(), insertDirectorQuery, director.ID, director.User.ID)
	return err
}

const signInQuery = `SELECT user_id, hash_pass, role FROM users WHERE username = $1`

func (r *authRepoPostgres) SignIn(username string) (domain.UserInfo, error) {
	row := r.conn.QueryRow(context.Background(), signInQuery, username)

	var user domain.UserInfo
	var userID string
	err := row.Scan(&userID, &user.Password, &user.UserRole)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.UserInfo{}, errs.ErrNoRows
		}
		return domain.UserInfo{}, err
	}
	user.ID = uuid.MustParse(userID)

	return user, nil
}

const confirmExpertQuery = "UPDATE experts SET is_confirmed = $1 WHERE expert_id = $2"

func (r *authRepoPostgres) ConfirmExpert(data domain.ConfirmExpert) error {
	n, err := r.conn.Exec(context.Background(), confirmExpertQuery, data.IsConfirmed, data.ExpertID)

	if n.RowsAffected() == 0 {
		return errs.ErrNoRows
	}
	return err
}

const insertCameraQuery = `INSERT INTO cameras (
                     camera_id, camera_type_id, camera_latitude, camera_longitude, short_desc, user_id) 
		VALUES ($1, $2, $3, $4, $5, $6)`

func (r *authRepoPostgres) InsertCamera(camera domain.Camera, userID uuid.UUID) error {

	_, err := r.conn.Exec(context.Background(), insertCameraQuery,
		camera.ID, camera.CameraType.ID, camera.Latitude, camera.Longitude, camera.ShortDesc, userID)

	if err != nil {
		return err
	}

	return nil
}
