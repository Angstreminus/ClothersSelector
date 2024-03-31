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
	service := service.NewUserService(repo, logger.Log)
	s.Logger.ZapLogger.Info("User service initialized")
	handler := handler.NewUserHandler(service, logger.Log)
	s.Logger.ZapLogger.Info("User handler initialized")
	router := http.NewServeMux()
	s.Router = router
	router.HandleFunc("/register", handler.RegisterUser)
	router.HandleFunc("/login", handler.LoginUser)
	if err := http.ListenAndServe(":8080", s.Router); err != nil {
		s.Logger.ZapLogger.Fatal("Error to run server")
	}
	s.Logger.ZapLogger.Info("Server is running")
}
