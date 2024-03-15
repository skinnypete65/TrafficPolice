package app

import (
	"TrafficPolice/internal/repository"
	postgres "TrafficPolice/internal/repository/postresql"
	"github.com/jackc/pgx/v5"
)

type repos struct {
	rating      repository.RatingRepo
	auth        repository.AuthRepo
	pagination  repository.PaginationRepo
	camera      repository.CameraRepo
	transport   repository.TransportRepo
	caseRepo    repository.CaseRepo
	contactInfo repository.ContactInfoRepo
	violation   repository.ViolationRepo
	training    repository.TrainingRepo
	checker     repository.CheckerRepo
	expert      repository.ExpertRepo
	director    repository.DirectorRepo
}

func newRepos(dbConn *pgx.Conn) *repos {
	return &repos{
		rating:      postgres.NewRatingRepoPostgres(dbConn),
		auth:        postgres.NewAuthRepoPostgres(dbConn),
		pagination:  postgres.NewPaginationRepoPostgres(dbConn),
		camera:      postgres.NewCameraRepoPostgres(dbConn),
		transport:   postgres.NewTransportRepoPostgres(dbConn),
		caseRepo:    postgres.NewCaseRepoPostgres(dbConn),
		contactInfo: postgres.NewContactInfoRepoPostgres(dbConn),
		violation:   postgres.NewViolationDBPostgres(dbConn),
		training:    postgres.NewTrainingRepoPostgres(dbConn),
		checker:     postgres.NewCheckerRepoPostgres(dbConn),
		expert:      postgres.NewExpertRepoPostgres(dbConn),
		director:    postgres.NewDirectorRepoPostgres(dbConn),
	}
}
