package servers

import (
	"log"

	"github.com/chothanin01/PhoSS-Care-server/configs"
	"github.com/chothanin01/PhoSS-Care-server/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Server struct {
	App *fiber.App
	Cfg *configs.Config
	Db  *gorm.DB
}

func NewServer(cfg *configs.Config, db *gorm.DB) *Server {
	return &Server{
		App: fiber.New(),
		Cfg: cfg,
		Db:  db,
	}
}

func (s *Server) Start() {
	if err := s.MapHandlers(); err != nil {
		log.Fatalf("Failed to map handlers: %v", err)
	}

	fiberConnURL, err := utils.ConnectionUrlBuilder("fiber", s.Cfg)
	if err != nil {
		log.Fatalf("Failed to build fiber connection URL: %v", err)
	}

	host := s.Cfg.App.Host
	port := s.Cfg.App.Port
	log.Printf("Server starting on %s:%s", host, port)

	if err := s.App.Listen(fiberConnURL); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
