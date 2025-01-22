package handlers

import (
	"log"
	"net/http"

	"code.gitea.io/sdk/gitea"
	"github.com/gofiber/fiber/v2"
)


type WebhookPayload struct {
	Ref        string      `json:"ref"`
	Before     string      `json:"before"`
	After      string      `json:"after"`
	CompareURL string      `json:"compare_url"`
	Commits    []gitea.Commit    `json:"commits"`
	Repository gitea.Repository  `json:"repository"`
	Pusher     gitea.User        `json:"pusher"`
	Sender     gitea.User        `json:"sender"`
}

// PostGiteaHook is a function to receive Gitea hook
//	@Summary		Receive Gitea hook
//	@Description	Receive Gitea hook
//	@Tags			Gitea
//	@Accept			json
//	@Produce		json
//	@Param			hook	body		WebhookPayload	true	"Gitea Hook"
//	@Success		200		{object}	ResponseHTTP{type=WebhookPayload}
//	@Failure		503		{object}	ResponseHTTP{}
//	@Router			/api/gitea [post]
func PostGiteaHook(c *fiber.Ctx) error {
	var payload WebhookPayload
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to parse hook",
		})
	}
	log.Printf("Received hook: %+v", payload)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully received hook",
		Data:    payload,
	})
}
