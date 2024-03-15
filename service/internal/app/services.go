package app

import (
	"TrafficPolice/internal/config"
	"TrafficPolice/internal/service"
	"TrafficPolice/internal/tokens"
	"TrafficPolice/pkg/hash"
)

type services struct {
	img         service.ImgService
	rating      service.RatingService
	auth        service.AuthService
	pagination  service.PaginationService
	camera      service.CameraService
	caseService service.CaseService
	contactInfo service.ContactInfoService
	violation   service.ViolationService
	expert      service.ExpertService
	training    service.TrainingService
	director    service.DirectorService
}

func newServices(r *repos, manager tokens.TokenManager, cfg *config.Config) *services {
	hasher := hash.NewSHA1Hasher(cfg.PassSalt)

	return &services{
		img:         service.NewImgService(),
		rating:      service.NewRatingService(r.rating, cfg.Rating),
		auth:        service.NewAuthService(r.auth, r.rating, hasher, manager),
		pagination:  service.NewPaginationService(r.pagination),
		camera:      service.NewCameraService(r.camera),
		caseService: service.NewCaseService(r.caseRepo, r.transport),
		contactInfo: service.NewContactInfoService(r.contactInfo),
		violation:   service.NewViolationService(r.violation),
		expert:      service.NewExpertService(r.expert, r.caseRepo, cfg.Consensus),
		training:    service.NewTrainingService(r.training),
		director:    service.NewDirectorService(r.director, r.checker),
	}
}
