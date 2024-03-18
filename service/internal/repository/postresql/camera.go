package repository

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/errs"
	"TrafficPolice/internal/repository"
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
)

type cameraRepoPostgres struct {
	conn *pgx.Conn
}

func NewCameraRepoPostgres(conn *pgx.Conn) repository.CameraRepo {
	return &cameraRepoPostgres{conn: conn}
}

const addCameraTypeQuery = `INSERT INTO camera_types (camera_type_id, camera_type_name) 
VALUES ($1, $2) RETURNING camera_type_id`

func (r *cameraRepoPostgres) AddCameraType(cameraType domain.CameraType) (string, error) {
	var cameraTypeID string

	err := r.conn.QueryRow(context.Background(), addCameraTypeQuery, cameraType.ID, cameraType.Name).
		Scan(&cameraTypeID)
	if err != nil {
		return "", errs.ErrAlreadyExists
	}

	return cameraTypeID, nil
}

const getCameraTypeByCameraIDQuery = `SELECT camera_type_name 
FROM camera_types as type 
JOIN cameras as c ON type.camera_type_id = c.camera_type_id
WHERE c.camera_id = $1`

func (r *cameraRepoPostgres) GetCameraTypeByCameraID(cameraID string) (string, error) {
	row := r.conn.QueryRow(context.Background(), getCameraTypeByCameraIDQuery, cameraID)

	var cameraType string
	err := row.Scan(&cameraType)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errs.ErrCameraNotExists
		}
		return "", err
	}

	return cameraType, nil
}
