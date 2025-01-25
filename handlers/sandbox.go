package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/HeavenManySugar/OJ-PoC/database"
	"github.com/HeavenManySugar/OJ-PoC/models"
)

// Specify the shell command for the corresponding repo
// @Summary		Specify the shell command for the corresponding repo
// @Description	Specify the shell command for the corresponding repo
// @Tags			sandbox
// @Accept			json
// @Produce		json
// @Param			cmd	body		models.Sandbox	true	"Shell command"
// @Success		200		{object}	ResponseHTTP{data=models.Sandbox}
// @Failure		503		{object}	ResponseHTTP{}
// @Router			/api/sandbox [post]
func PostSandboxCmd(c *fiber.Ctx) error {
	db := database.DBConn

	cmd := new(models.Sandbox)
	if err := c.BodyParser(cmd); err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to parse shell command",
		})
	}

	db.Create(cmd)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: fmt.Sprintf("Success set shell command for %v.", cmd.SourceGitRepo),
		Data:    *cmd,
	})
}