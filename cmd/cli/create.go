package cli

import (
	"fmt"
	"log"
	"net/url" // Pour valider le format de l'URL
	"os"

	cmd2 "github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/repository"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite" // Driver SQLite pour GORM
	"gorm.io/gorm"
)

// TODO : Faire une variable longURLFlag qui stockera la valeur du flag --url
var longURLFlag string

// CreateCmd représente la commande 'create'
var CreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Crée une URL courte à partir d'une URL longue.",
	Long: `Cette commande raccourcit une URL longue fournie et affiche le code court généré.

Exemple:
  url-shortener create --url="https://www.google.com/search?q=go+lang"`,
	Run: func(cmd *cobra.Command, args []string) {
		if longURLFlag == "" {
			fmt.Fprintln(os.Stderr, "Erreur: le flag --url est requis.")
			os.Exit(1)
		}
		if _, err := url.ParseRequestURI(longURLFlag); err != nil {
			fmt.Fprintf(os.Stderr, "Erreur: format d'URL invalide: %v\n", err)
			os.Exit(1)
		}
		cfg := cmd2.Cfg
		if cfg == nil {
			fmt.Fprintln(os.Stderr, "Erreur: configuration non chargée.")
			os.Exit(1)
		}
		db, err := gorm.Open(sqlite.Open(cfg.Database.Name), &gorm.Config{})
		if err != nil {
			log.Fatalf("FATAL: Échec de la connexion à la base de données: %v", err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("FATAL: Échec de l'obtention de la base de données SQL sous-jacente: %v", err)
		}
		defer sqlDB.Close()
		linkRepo := repository.NewLinkRepository(db)
		linkService := services.NewLinkService(linkRepo)
		link, err := linkService.CreateLink(longURLFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erreur lors de la création du lien: %v\n", err)
			os.Exit(1)
		}
		fullShortURL := fmt.Sprintf("%s/%s", cfg.Server.BaseURL, link.ShortCode)
		fmt.Printf("URL courte créée avec succès:\n")
		fmt.Printf("Code: %s\n", link.ShortCode)
		fmt.Printf("URL complète: %s\n", fullShortURL)
	},
}

// init() s'exécute automatiquement lors de l'importation du package.
// Il est utilisé pour définir les flags que cette commande accepte.
func init() {
	CreateCmd.Flags().StringVar(&longURLFlag, "url", "", "URL longue à raccourcir")
	CreateCmd.MarkFlagRequired("url")
	cmd2.RootCmd.AddCommand(CreateCmd)
}
