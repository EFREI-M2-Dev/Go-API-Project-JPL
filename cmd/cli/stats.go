package cli

import (
	"fmt"
	"github.com/axellelanca/urlshortener/internal/config"
	"log"
	"os"
	"sync"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"

	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : variable shortCodeFlag qui stockera la valeur du flag --code
var (
	inputShortenedURL string
)

// StatsCmd représente la commande 'stats'
var StatsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Affiche les statistiques (nombre de clics) pour un lien court.",
	Long: `Cette commande permet de récupérer et d'afficher le nombre total de clics
pour une URL courte spécifique en utilisant son code.

Exemple:
  url-shortener stats --code="xyz123"`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO : Valider que le flag --code a été fourni.
		if inputShortenedURL == "" {
			fmt.Println("Aucun code d'URL raccourcie n'a été fournie.")
			os.Exit(1)
		}

		// TODO : Charger la configuration chargée globalement via cmd.cfg
		configs, err := config.LoadConfig()
		if err != nil {
			fmt.Printf("Erreur lors du chargement de la configuration : %v\n", err)
			os.Exit(1)
		}

		// TODO 3: Initialiser la connexion à la base de données SQLite avec GORM.
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("Erreur lors de l'ouverture de la base SQLite : %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		// TODO S'assurer que la connexion est fermée à la fin de l'exécution de la commande
		defer sqlDB.Close()

		// TODO : Initialiser les repositories et services nécessaires NewLinkRepository & NewLinkService
		repo := repository.NewLinkRepository(db)
		service := services.NewLinkService(repo)

		// TODO 5: Appeler GetLinkStats pour récupérer le lien et ses statistiques.
		link, totalClicks, err := service.GetLinkStats(inputShortenedURL)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				fmt.Printf("Aucun lien trouvé pour le code: %s\n", inputShortenedURL)
			} else {
				fmt.Printf("Erreur lors de la récupération des statistiques : %v\n", err)
			}
			os.Exit(1)
		}

		fmt.Printf("Statistiques pour le code court: %s\n", link.ShortCode)
		fmt.Printf("URL longue: %s\n", link.LongURL)
		fmt.Printf("Total de clics: %d\n", totalClicks)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	StatsCmd.Flags().StringVarP(&inputShortenedURL, "stats", "s", "", "shortened URL")
	// TODO 7: Définir le flag --code pour la commande stats.

	StatsCmd.Flags().StringVar(&inputShortenedURL, "code", "", "Code court de l'URL raccourcie")
	// TODO Marquer le flag comme requis

	StatsCmd.MarkFlagRequired("code")
	cmd2.RootCmd.AddCommand(StatsCmd)
}
