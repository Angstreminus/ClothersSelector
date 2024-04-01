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
	middleware "github.com/Angstreminus/ClothersSelector/middleware/auth"
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
	cfg, err := config.NewConfig()
	if err != nil {
		s.Logger.ZapLogger.Error("Error to load config")
	}
	dbHandler, err := postgres.NewDatabaseHandler(s.Config)
	if err != nil {
		s.Logger.ZapLogger.Error("Error to init postgres")
	}
	redis, err := chache.NewChache(s.Config)
	if err != nil {
		s.Logger.ZapLogger.Error("Error to init redis")
	}
	userRepo := repository.NewUserRepository(redis, dbHandler, s.Logger)
	presetRepo := repository.NewPresetRepository(dbHandler, s.Logger)
	s.Logger.ZapLogger.Info("User repository initialized")
	userService := service.NewUserService(userRepo, logger.Log)
	prestService := service.NewPresetService(presetRepo, s.Logger)
	s.Logger.ZapLogger.Info("User service initialized")
	userHandler := handler.NewUserHandler(userService, logger.Log, cfg)
	presetHandler := handler.NewPresetHandler(cfg, prestService, s.Logger)
	s.Logger.ZapLogger.Info("User handler initialized")
	authmiddleware := middleware.NewAuthMiddleware(cfg, userRepo)
	router := http.NewServeMux()
	s.Router = router
	router.HandleFunc("POST /register", userHandler.RegisterUser)
	router.HandleFunc("POST /login", userHandler.LoginUser)
	router.HandleFunc("POST /users/:id/presets", authmiddleware.ValidateToken(presetHandler.CreatePreset))
	// !TODO:
	// !!! FIX code below
	router.HandleFunc("GET /users/:id/presets", authmiddleware.ValidateToken(presetHandler.CreatePreset))
	router.HandleFunc("GET /users/:id/presets/:id", authmiddleware.ValidateToken(presetHandler.CreatePreset))
	router.HandleFunc("DELETE /users/:id/presets/:id", authmiddleware.ValidateToken(presetHandler.CreatePreset))
	router.HandleFunc("POST /users/:id/presets/:id/clothes", authmiddleware.ValidateToken(presetHandler.CreatePreset))
	router.HandleFunc("POST /users/:id/presets/:id/clothes", authmiddleware.ValidateToken(presetHandler.CreatePreset))
	router.HandleFunc("POST /users/:id/presets/:id/clothes", authmiddleware.ValidateToken(presetHandler.CreatePreset))
	if err := http.ListenAndServe(":8080", s.Router); err != nil {
		s.Logger.ZapLogger.Fatal("Error to run server")
	}
	s.Logger.ZapLogger.Info("Server is running")
}
