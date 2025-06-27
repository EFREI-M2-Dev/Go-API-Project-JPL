package config

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

// DONE Créer Config qui est la structure principale qui mappe l'intégralité de la configuration de l'application.
// Les tags `mapstructure` sont utilisés par Viper pour mapper les clés du fichier de config
// (ou des variables d'environnement) aux champs de la structure Go.
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Analytics AnalyticsConfig `mapstructure:"analytics"`
	Monitor   MonitorConfig   `mapstructure:"monitor"`
	Workers   WorkersConfig   `mapstructure:"workers"`
}

type ServerConfig struct {
	Port    int    `mapstructure:"port"`
	BaseURL string `mapstructure:"base_url"`
}

type DatabaseConfig struct {
	Name string `mapstructure:"name"`
}

type AnalyticsConfig struct {
	BufferSize  int `mapstructure:"buffer_size"`
	WorkerCount int `mapstructure:"worker_count"`
}

type MonitorConfig struct {
	IntervalMinutes int `mapstructure:"interval_minutes"`
}

type WorkersConfig struct {
	ClickEventsBufferSize int `mapstructure:"click_events_buffer_size"`
}

func LoadConfig() (*Config, error) {
	// Load config from 'configs' directory
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")

	// DONE : Définir les valeurs par défaut pour toutes les options de configuration.
	// DONE : Lire le fichier de configuration.

	if err := viper.ReadInConfig(); err != nil {
		// If error in config, switch to default values
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("Fichier de configuration non trouvé, utilisation des valeurs par défaut.")
			// Load default values
			viper.SetDefault("server.port", 8080)
			viper.SetDefault("server.base_url", "http://localhost")
			viper.SetDefault("database.name", "url_shortener.db")
			viper.SetDefault("analytics.buffer_size", 1000)
			viper.SetDefault("analytics.worker_count", 5)
			viper.SetDefault("monitor.interval_minutes", 5)
		} else {
			log.Printf("Erreur lors de la lecture du fichier de configuration: %v", err)
		}
	} else {
		log.Println("Fichier de configuration chargé avec succès.")
	}

	// DONE 4: Démapper (unmarshal) la configuration lue (ou les valeurs par défaut) dans la structure Config.
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Printf("Erreur lors du démappage de la configuration: %v", err)
		return nil, err
	}

	log.Printf("Configuration loaded: Server Port=%d, DB Name=%s, Analytics Buffer=%d, Monitor Interval=%dmin",
		cfg.Server.Port, cfg.Database.Name, cfg.Analytics.BufferSize, cfg.Monitor.IntervalMinutes)

	return &cfg, nil
}
