package repository

import (
	"TrafficPolice/errs"
	"TrafficPolice/internal/domain"
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

func (r *cameraRepoPostgres) AddCameraType(cameraType domain.CameraType) error {
	query := "INSERT INTO camera_types (camera_type_id, camera_type_name) VALUES ($1, $2)"

	_, err := r.conn.Exec(context.Background(), query, cameraType.ID, cameraType.Name)
	if err != nil {
		return err
	}

	return nil
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
			return "", errs.ErrNoRows
		}
		return "", err
	}

	return cameraType, nil
}
