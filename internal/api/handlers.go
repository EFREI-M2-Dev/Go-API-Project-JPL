package api

import (
	"errors"
	"github.com/axellelanca/urlshortener/cmd"
	"github.com/axellelanca/urlshortener/internal/models"
	"github.com/axellelanca/urlshortener/internal/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// ClickEventsChannel est le channel global (ou injecté) utilisé pour envoyer les événements de clic
var ClickEventsChannel chan *models.ClickEvent

// SetupRoutes configure toutes les routes de l'API Gin et injecte les dépendances nécessaires
func SetupRoutes(router *gin.Engine, linkService *services.LinkService) {
	if ClickEventsChannel == nil {
		ClickEventsChannel = make(chan *models.ClickEvent, cmd.Cfg.Analytics.BufferSize)
	}

	v1 := router.Group("/api/v1")
	v1.GET("/health", HealthCheckHandler)
	v1.POST("/links", CreateShortLinkHandler(linkService))
	v1.GET("/links/:shortCode/stats", GetLinkStatsHandler(linkService))

	router.GET("/:shortCode", RedirectHandler(linkService))
}

// HealthCheckHandler gère la route /health pour vérifier l'état du service.
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// CreateLinkRequest représente le corps de la requête JSON pour la création d'un lien.
type CreateLinkRequest struct {
	LongURL string `json:"long_url" binding:"required,url"` // 'binding:required' pour validation, 'url' pour format URL
}

// CreateShortLinkHandler gère la création d'une URL courte.
func CreateShortLinkHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateLinkRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide ou URL incorrecte"})
			return
		}

		link, err := linkService.CreateLink(req.LongURL)
		if err != nil {
			log.Printf("[Handlers::CreateLink] Erreur lors de la création du lien : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur Serveur"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"short_code":     link.ShortCode,
			"long_url":       link.LongURL,
			"full_short_url": cmd.Cfg.Server.BaseURL + "/" + link.ShortCode,
		})
	}
}

// RedirectHandler gère la redirection d'une URL courte vers l'URL longue et l'enregistrement asynchrone des clics.
func RedirectHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		// TODO 2: Récupérer l'URL longue associée au shortCode depuis le linkService (GetLinkByShortCode)
		link, err := linkService.GetLinkByShortCode(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Si le lien n'est pas trouvé, retourner HTTP 404 Not Found.
				// Utiliser errors.Is et l'erreur Gorm
				c.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
				return
			}
			// Gérer d'autres erreurs potentielles de la base de données ou du service
			log.Printf("Error retrieving link for %s: %v", shortCode, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		// TODO 3: Créer un ClickEvent avec les informations pertinentes.
		clickEvent := &models.ClickEvent{
			LinkID:    link.ID,
			ShortCode: link.ShortCode,
			Timestamp: time.Now(),
			UserAgent: c.Request.UserAgent(),
			IPAddress: c.ClientIP(),
		}

		// TODO 4: Envoyer le ClickEvent dans le ClickEventsChannel avec le Multiplexage.
		// Utilise un `select` avec un `default` pour éviter de bloquer si le channel est plein.
		// Pour le default, juste un message à afficher :
		// log.Printf("Warning: ClickEventsChannel is full, dropping click event for %s.", shortCode)
		select {
		case ClickEventsChannel <- clickEvent:
			log.Printf("Click event for %s sent to channel.", shortCode)
		default:
			log.Printf("Warning: ClickEventsChannel is full, dropping click event for %s.", shortCode)
		}

		c.Redirect(http.StatusFound, link.LongURL)
	}
}

// GetLinkStatsHandler gère la récupération des statistiques pour un lien spécifique.
func GetLinkStatsHandler(linkService *services.LinkService) gin.HandlerFunc {
	return func(c *gin.Context) {
		shortCode := c.Param("shortCode")

		link, totalClicks, err := linkService.GetLinkStats(shortCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Lien non trouvé"})
				return
			}
			log.Printf("[Handlers::GetLinkStatsHandler] Erreur lors de la récupération des stats : %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur serveur"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"short_code":   link.ShortCode,
			"long_url":     link.LongURL,
			"total_clicks": totalClicks,
		})
	}
}
