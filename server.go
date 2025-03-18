package main

import (
	"api-backend-go/config"
	"api-backend-go/handler"
	"api-backend-go/middleware"
	"api-backend-go/model"
	"api-backend-go/repository"
	"api-backend-go/service"
	"api-backend-go/usecase"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Server struct {
	userUC     usecase.UserUseCase
	authUC     usecase.AuthenticationUseCase
	jwtService service.JwtService
	engine     *gin.Engine
	host       string
}

func (s *Server) initMigration() {
	if DB == nil {
		log.Fatal("Database connection is nil, migration failed!")
	}

	err := DB.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("Failed to migrate database: " + err.Error())
	}

	fmt.Println("Migrate database successfully")
}

func (s *Server) initRoute() {
	rgAuth := s.engine.Group("/api")
	handler.NewAuthHandler(s.authUC, rgAuth).Route()

	authMiddleware := middleware.NewAuthMiddleware(s.jwtService)

	rgV1 := s.engine.Group("/api/v1")
	handler.NewUserHandler(s.userUC, rgV1, authMiddleware).Route()
}

func (s *Server) Run() {
	s.initMigration()
	s.initRoute()

	if err := s.engine.Run(s.host); err != nil {
		panic(fmt.Errorf("server not running on host %s, because error %v", s.host, err))
	}
}
func NewServer() *Server {
	var err error

	cfg, err := config.NewConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBConfig.Username, cfg.DBConfig.Password, cfg.DBConfig.Host, cfg.DBConfig.Port, cfg.DBConfig.Database)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connection db", err)
	} else {
		log.Println("Successfully connected to database")
	}

	jwtService := service.NewJwtService(cfg.TokenConfig)

	userRepo := repository.NewUserRepository(DB)

	userUseCase := usecase.NewUserUseCase(userRepo)
	authUseCase := usecase.NewAuthenticationImpl(userUseCase, jwtService)

	userUseCase.StartWorkerPool(5)

	engine := gin.Default()
	host := fmt.Sprintf(":%s", cfg.ApiPort)
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           time.Duration(cfg.MaxAge),
	}))
	return &Server{
		engine:     engine,
		userUC:     userUseCase,
		authUC:     authUseCase,
		jwtService: jwtService,
		host:       host,
	}

}
