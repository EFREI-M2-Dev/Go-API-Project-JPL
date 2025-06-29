package cli

import (
	"fmt"
	"github.com/axellelanca/urlshortener/internal/config"
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// MigrateCmd représente la commande 'migrate'
var MigrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Exécute les migrations de la base de données pour créer ou mettre à jour les tables.",
	Long: `Cette commande se connecte à la base de données configurée (SQLite)
et exécute les migrations automatiques de GORM pour créer les tables 'links' et 'clicks'
basées sur les modèles Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : Charger la configuration chargée globalement via cmd.cfg
		configs, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Erreur lors du chargement de la configuration : %v\n", err)
			os.Exit(1)
		}

		// TODO 2: Initialiser la connexion à la base de données SQLite avec GORM.
		db, err := gorm.Open(sqlite.Open(configs.Database.Name), &gorm.Config{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de l'ouverture de la base SQLite : %v", err)
			os.Exit(1)
		}

		sqlDB, err := db.DB()
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
			os.Exit(1)
		}
		// TODO Assurez-vous que la connexion est fermée après la migration.
		defer sqlDB.Close()

		// TODO 3: Exécuter les migrations automatiques de GORM.
		err = db.AutoMigrate(&models.Link{}, &models.Click{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de l'exécution des migrations : %v", err)
			os.Exit(1)
		}

		// Pas touche au log
		fmt.Println("Migrations de la base de données exécutées avec succès.")
	},
}

func init() {
	// TODO : Ajouter la commande à RootCmd
	cmd2.RootCmd.AddCommand(MigrateCmd)
}
