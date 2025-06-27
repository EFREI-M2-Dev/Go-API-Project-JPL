package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	"gorm.io/gorm" // Nécessaire pour la gestion spécifique de gorm.ErrRecordNotFound

	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/repository" // Importe le package repository
)

// Définition du jeu de caractères pour la génération des codes courts.
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type LinkService struct {
	linkRepo repository.LinkRepository // Référence vers le repository de liens
}

// NewLinkService crée et retourne une nouvelle instance de LinkService.
func NewLinkService(linkRepo repository.LinkRepository) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

// GenerateShortCode est une méthode rattachée à LinkService
func GenerateShortCode(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("[GenerateShortCode] la longueur doit être supérieure à 0")
	}

	shortCode := make([]byte, length)
	for i := range shortCode {
		// Génère un index aléatoire dans le charset
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("[GenerateShortCode] Erreur de génération du code: %w", err)
		}
		shortCode[i] = charset[index.Int64()]
	}

	return string(shortCode), nil
}

// CreateLink crée un nouveau lien raccourci.
// Il génère un code court unique, puis persiste le lien dans la base de données.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	const maxRetries = 5
	var shortCode string
	var err error

	for i := 0; i < maxRetries; i++ {
		shortCode, err = GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("erreur lors de la génération du code court: %w", err)
		}

		_, err = s.linkRepo.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Le code est unique
				break
			}
			return nil, fmt.Errorf("erreur lors de la vérification d'unicité du code court: %w", err)
		}

		// Collision détectée, on log et on retente
		log.Printf("Short code '%s' déjà existant, nouvelle tentative (%d/%d)...", shortCode, i+1, maxRetries)
	}

	if err == nil {
		// Si on sort de la boucle sans erreur, c'est qu'on a trouvé un code existant à chaque fois
		return nil, errors.New("impossible de générer un code court unique après plusieurs tentatives")
	}

	link := &models.Link{
		LongURL:   longURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
	}

	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("erreur lors de la création du lien: %w", err)
	}

	return link, nil
}

// GetLinkByShortCode récupère un lien via son code court.
// Il délègue l'opération de recherche au repository.
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	// TODO : Récupérer un lien par son code court en utilisant s.linkRepo.GetLinkByShortCode.
	// Retourner le lien trouvé ou une erreur si non trouvé/problème DB.

}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// TODO : Récupérer le lien par son shortCode

	// TODO 4: Compter le nombre de clics pour ce LinkID

	// TODO : on retourne les 3 valeurs
	return
}
