package handlers

import (
	"log"
	"net/http"
	"net/url"

	"github.com/HeavenManySugar/OJ-PoC/database"
	"github.com/HeavenManySugar/OJ-PoC/models"
	"github.com/gofiber/fiber/v2"
)

// GetScores is a function to get all scores
//	@Summary		Get all scores
//	@Description	Get all scores
//	@Tags			Score
//	@Accept			json
//	@Produce			json
//	@Success		200		{object}	ResponseHTTP{type=[]models.Score}
//	@Failure		503		{object}	ResponseHTTP{}
//	@Router			/api/scores [get]
func GetScores(c *fiber.Ctx) error {
	db := database.DBConn
	var scores []models.Score
	if err := db.Raw(`
		SELECT * FROM scores
		WHERE id IN (
			SELECT id FROM (
				SELECT id, ROW_NUMBER() OVER (PARTITION BY git_repo ORDER BY updated_at DESC) AS rn
				FROM scores
			) AS subquery
			WHERE rn = 1
		)
		ORDER BY updated_at DESC
	`).Scan(&scores).Error; err != nil {
		log.Printf("Failed to get scores: %v", err)
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to get scores",
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully get scores",
		Data:    scores,
	})
}

// GetScoreByRepo is a function to get a score by repo
//	@Summary		Get a score by repo
//	@Description	Get a score by repo
//	@Tags			Score
//	@Accept			json
//	@Produce			json
//	@Param			repo	path	string	true	"Repo name"
//	@Success		200		{object}	ResponseHTTP{type=models.Score}
//	@Failure		404		{object}	ResponseHTTP{}
//	@Failure		503		{object}	ResponseHTTP{}
//	@Router			/api/score/{repo} [get]
func GetScoreByRepo(c *fiber.Ctx) error {
	db := database.DBConn
	repo, err := url.PathUnescape(c.Params("repo"))
	if err != nil {
		log.Printf("Failed to unescape repo name: %v", err)
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to unescape repo name",
		})
	}
	var score models.Score
	if err := db.Where("git_repo = ?", repo).Order("updated_at DESC").First(&score).Error; err != nil {
		log.Printf("Failed to get score by repo: %v", err)
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to get score by repo",
		})
	}
	if score.ID == 0 {
		return c.Status(http.StatusNotFound).JSON(ResponseHTTP{
			Success: false,
			Message: "Score not found",
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully get score by repo",
		Data:    score,
	})
}