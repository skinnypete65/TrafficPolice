package app

import (
	"TrafficPolice/internal/domain"
	"TrafficPolice/internal/transport/rest/middlewares"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

type ServeMuxInit struct {
	h              *handlers
	authMiddleware *middlewares.AuthMiddleware
	mux            *http.ServeMux
}

func newServeMuxInit(h *handlers, authMiddleware *middlewares.AuthMiddleware) *ServeMuxInit {
	return &ServeMuxInit{
		h:              h,
		authMiddleware: authMiddleware,
		mux:            http.NewServeMux(),
	}
}

func (s *ServeMuxInit) Init() *http.ServeMux {
	// Setup Routes
	s.mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	s.initCameraHandlers()
	s.initCaseHandlers()
	s.initContactInfoHandlers()
	s.initViolationHandlers()
	s.initAuthHandlers()
	s.initExpertHandlers()
	s.initRatingHandlers()
	s.initDirectorHandlers()

	return s.mux
}

func (s *ServeMuxInit) initCameraHandlers() {
	s.mux.Handle("POST /camera/type",
		s.authMiddleware.IdentifyRole(http.HandlerFunc(s.h.camera.AddCameraType), domain.DirectorRole),
	)

	s.mux.Handle("POST /camera",
		s.authMiddleware.IdentifyRole(http.HandlerFunc(s.h.camera.RegisterCamera), domain.DirectorRole),
	)
}
func (s *ServeMuxInit) initCaseHandlers() {
	s.mux.Handle("POST /case",
		s.authMiddleware.IdentifyRole(http.HandlerFunc(s.h.caseHandler.AddCase), domain.CameraRole),
	)

	s.mux.Handle("POST /case/{id}/img",
		s.authMiddleware.IdentifyRole(http.HandlerFunc(s.h.caseHandler.UploadCaseImg), domain.CameraRole),
	)

	s.mux.Handle("GET /case/{id}/img",
		s.authMiddleware.IdentifyRole(
			s.authMiddleware.IsExpertConfirmed(http.HandlerFunc(s.h.caseHandler.GetCaseImg)),
			domain.DirectorRole, domain.ExpertRole,
		),
	)
}

func (s *ServeMuxInit) initContactInfoHandlers() {
	s.mux.Handle("POST /contact_info",
		s.authMiddleware.IdentifyRole(http.HandlerFunc(s.h.contactInfo.InsertContactInfo), domain.DirectorRole),
	)
}

func (s *ServeMuxInit) initViolationHandlers() {
	s.mux.Handle("POST /violations",
		s.authMiddleware.IdentifyRole(http.HandlerFunc(s.h.violation.InsertViolations), domain.DirectorRole),
	)
}

func (s *ServeMuxInit) initAuthHandlers() {
	s.mux.HandleFunc("POST /auth/sign_up", s.h.auth.SignUp)

	s.mux.HandleFunc("POST /auth/sign_in", s.h.auth.SignIn)

	s.mux.Handle("POST /auth/confirm/expert",
		s.authMiddleware.IdentifyRole(http.HandlerFunc(s.h.auth.ConfirmExpert), domain.DirectorRole),
	)
}

func (s *ServeMuxInit) initExpertHandlers() {
	s.mux.Handle("POST /expert/{id}/img", s.authMiddleware.IdentifyRole(
		s.authMiddleware.IsExpertConfirmed(http.HandlerFunc(s.h.expert.UploadExpertImg)),
		domain.DirectorRole, domain.ExpertRole),
	)

	s.mux.Handle("GET /expert/{id}/img",
		s.authMiddleware.IdentifyRole(
			s.authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(s.h.expert.GetExpertImg),
			),
			domain.DirectorRole, domain.ExpertRole,
		),
	)

	s.mux.Handle("GET /expert/case",
		s.authMiddleware.IdentifyRole(
			s.authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(s.h.expert.GetCaseForExpert),
			),
			domain.ExpertRole,
		),
	)

	s.mux.Handle("POST /expert/decision",
		s.authMiddleware.IdentifyRole(
			s.authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(s.h.expert.SetCaseDecision),
			),
			domain.ExpertRole,
		),
	)

	s.mux.Handle("POST /expert/training",
		s.authMiddleware.IdentifyRole(
			s.authMiddleware.IsExpertConfirmed(
				http.HandlerFunc(s.h.training.GetSolvedCasesByParams),
			),
			domain.ExpertRole,
		),
	)
}

func (s *ServeMuxInit) initRatingHandlers() {
	s.mux.Handle("GET /rating",
		s.authMiddleware.IdentifyRole(
			http.HandlerFunc(s.h.rating.GetRating),
			domain.ExpertRole, domain.DirectorRole,
		),
	)
}

func (s *ServeMuxInit) initDirectorHandlers() {
	s.mux.Handle("GET /director/case",
		s.authMiddleware.IdentifyRole(
			http.HandlerFunc(s.h.director.GetCase),
			domain.DirectorRole,
		),
	)

	s.mux.Handle("GET /director/analytics/expert",
		s.authMiddleware.IdentifyRole(
			http.HandlerFunc(s.h.director.ExpertAnalytics),
			domain.DirectorRole,
		),
	)

	s.mux.Handle("PATCH /director/expert_skill",
		s.authMiddleware.IdentifyRole(
			http.HandlerFunc(s.h.director.UpdateExpertSkill),
			domain.DirectorRole,
		),
	)
}
