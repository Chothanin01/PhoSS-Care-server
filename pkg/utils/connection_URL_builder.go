package utils

import (
	"fmt"

	"github.com/chothanin01/PhoSS-Care-server/configs"
)

func ConnectionUrlBuilder(target string, cfg *configs.Config) (string, error) {
	switch target {
	case "fiber":
		return fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port), nil

	case "gorm":
		return fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s TimeZone=Asia/Bangkok",
			cfg.DB.Host,
			cfg.DB.Username,
			cfg.DB.Password,
			cfg.DB.DBName,
			cfg.DB.Port,
		), nil

	default:
		return "", fmt.Errorf("unsupported connection target: %s", target)
	}
}
