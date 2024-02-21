package database

import (
	"TrafficPolice/internal/database"
	"TrafficPolice/internal/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
)

type DBPostgres struct {
	conn *pgx.Conn
}

func NewCameraDBPostgres(conn *pgx.Conn) database.CameraDB {
	return &DBPostgres{conn: conn}
}

func (db *DBPostgres) AddCameraType(cameraType models.CameraType) error {
	query := "INSERT INTO camera_types (camera_type_id, camera_type_name, camera_type_desc) VALUES ($1, $2, $3)"

	_, err := db.conn.Exec(context.Background(), query, cameraType.ID, cameraType.Name, cameraType.Desc)
	if err != nil {
		return fmt.Errorf("unable to insert camera type %v", cameraType)
	}

	return nil
}

func (db *DBPostgres) RegisterCamera(camera models.Camera) error {
	query := `INSERT INTO cameras (camera_id, camera_type_id, camera_latitude, camera_longitude, short_desc) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := db.conn.Exec(context.Background(), query,
		camera.ID, camera.CameraTypeID, camera.Latitude, camera.Longitude, camera.ShortDesc)

	if err != nil {
		return fmt.Errorf("unable to insert camera: %v", camera)
	}

	return nil
}
