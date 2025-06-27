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
		return "", errors.New("[Service::GenerateShortCode] la longueur doit être supérieure à 0")
	}

	shortCode := make([]byte, length)
	for i := range shortCode {
		// Génère un index aléatoire dans le charset
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("[Service::GenerateShortCode] Erreur de génération du code: %w", err)
		}
		shortCode[i] = charset[index.Int64()]
	}

	return string(shortCode), nil
}

// CreateLink crée un nouveau lien raccourci.
func (s *LinkService) CreateLink(longURL string) (*models.Link, error) {
	const maxRetries = 5
	var shortCode string
	var err error

	for i := 0; i < maxRetries; i++ {
		shortCode, err = GenerateShortCode(6)
		if err != nil {
			return nil, fmt.Errorf("[Service::CreateLink] Erreur lors de la génération du code court: %w", err)
		}

		_, err = s.linkRepo.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Le code est unique on peut sortir de la boucle
				break
			}
			return nil, fmt.Errorf("[Service::CreateLink] Erreur lors de la vérification d'unicité du code court: %w", err)
		}

		// Collision détectée, on log et on retente
		log.Printf("[Service::CreateLink] Short code '%s' déjà existant, nouvelle tentative (%d/%d)...", shortCode, i+1, maxRetries)
	}

	if err == nil {
		// Si on sort de la boucle sans erreurs, c'est qu'on a trouvé un code existant à chaque fois
		return nil, errors.New("[Service::CreateLink] Impossible de gén un code court unique après plusieurs tentatives")
	}

	link := &models.Link{
		LongURL:   longURL,
		ShortCode: shortCode,
		CreatedAt: time.Now(),
	}

	if err := s.linkRepo.CreateLink(link); err != nil {
		return nil, fmt.Errorf("[Service::CreateLink] erreur lors de la création du lien: %w", err)
	}

	return link, nil
}

// GetLinkByShortCode récupère un lien via son court code
func (s *LinkService) GetLinkByShortCode(shortCode string) (*models.Link, error) {
	link, err := s.linkRepo.GetLinkByShortCode(shortCode)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("[Service::GetLinkByShortCode] Lien non trouvé pour le code court '%s': %w", shortCode, err)
		}
		return nil, fmt.Errorf("[Service::GetLinkByShortCode] Erreur lors de la récupération du lien: %w", err)
	}
	return link, nil
}

// GetLinkStats récupère les statistiques pour un lien donné (nombre total de clics).
// Il interagit avec le LinkRepository pour obtenir le lien, puis avec le ClickRepository
func (s *LinkService) GetLinkStats(shortCode string) (*models.Link, int, error) {
	// TODO : Récupérer le lien par son shortCode
	link, err := s.GetLinkByShortCode(shortCode)
	if err != nil {
		return nil, 0, fmt.Errorf("[Service::GetLinkStats] Erreur lors de la récupération du lien: %w", err)
	}
	if link == nil {
		return nil, 0, fmt.Errorf("[Service::GetLinkStats] Lien non trouvé pour le code court '%s'", shortCode)
	}

	clicksCount, err := s.linkRepo.CountClicksByLinkID(link.ID)
	if err != nil {
		return nil, 0, fmt.Errorf("[Service::GetLinkStats] Erreur lors du comptage des clics pour le lien ID %d: %w", link.ID, err)
	}

	// On retourne les 3 valeurs demandées
	return link, clicksCount, nil
}
