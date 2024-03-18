package server

import (
	"net/http"

	"github.com/Angstreminus/ClothersSelector/config"
	"github.com/Angstreminus/ClothersSelector/internal/chache"
	"github.com/Angstreminus/ClothersSelector/internal/handler"
	"github.com/Angstreminus/ClothersSelector/internal/postgres"
	"github.com/Angstreminus/ClothersSelector/internal/repository"
	"github.com/Angstreminus/ClothersSelector/internal/service"
	"github.com/Angstreminus/ClothersSelector/logger"
)

type Server struct {
	Config *config.Config
	Router *http.ServeMux
	Logger *logger.Logger
}

func NewServer(cfg *config.Config, logger *logger.Logger) *Server {
	return &Server{
		Config: cfg,
		Logger: logger,
	}
}

func (s *Server) MustRun() {
	dbHandler, err := postgres.NewDatabaseHandler(s.Config)
	if err != nil {
		s.Logger.ZapLogger.Error("Error to init postgres")
	}
	redis, err := chache.NewChache(s.Config)
	if err != nil {
		s.Logger.ZapLogger.Error("Error to init redis")
	}
	repo := repository.NewUserRepository(redis, dbHandler, s.Logger)
	s.Logger.ZapLogger.Info("User repository initialized")
	service := service.NewUserService(repo)
	s.Logger.ZapLogger.Info("User service initialized")
	handler := handler.NewUserHandler(service)
	s.Logger.ZapLogger.Info("User handler initialized")

	router := http.NewServeMux()
	s.Logger.ZapLogger.Info("Router initialized")
	s.Router = router
	router.HandleFunc("POST /register", handler.RegisterUser)
	router.HandleFunc("POST /login", handler.LoginUser)
	router.HandleFunc("POST /actors")
	router.HandleFunc("PUT /actors/{id}/")
	router.HandleFunc("DELETE /actors/{id}")
	router.HandleFunc("GET /movies")
	router.HandleFunc("GET /movies/{id}/actors")
	router.HandleFunc("POST /movies")
	router.HandleFunc("PUT /movies/{id}")
	router.HandleFunc("DELETE /movies/{id}")

	if err := http.ListenAndServe(s.Config.ServerAddr, s.Router); err != nil {
		s.Logger.ZapLogger.Fatal("Error to run server")
	}
	s.Logger.ZapLogger.Info("Server is running")
}
