package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"code.gitea.io/sdk/gitea"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gofiber/fiber/v2"

	"github.com/HeavenManySugar/OJ-PoC/sandbox"
)

const GitServer = "http://server.gitea.orb.local/"

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

	// Clone the given repository to the given directory
	log.Printf("git clone %s", GitServer+payload.Repository.FullName)
	r, err := git.PlainClone("/tmp/" + payload.Repository.FullName, false, &git.CloneOptions{
		URL:      GitServer + payload.Repository.FullName,
		Progress: os.Stdout,
	})
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to clone repository",
		})
	}
	log.Printf("git show-ref --head HEAD")
	ref, err := r.Head()
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to get HEAD",
		})
	}
	fmt.Println(ref.Hash())

	w, err := r.Worktree()
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to get worktree",
		})
	}
	err = w.Checkout(&git.CheckoutOptions{
		Hash: plumbing.NewHash(payload.After),
	})
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to checkout",
		})
	}
	ref, err = r.Head()
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to get HEAD",
		})
	}
	fmt.Println(ref.Hash())
	
	sandbox.SandboxPtr.RunShellCommandByRepo(payload.Repository.Parent.FullName, nil)

	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully received hook",
		Data:    payload,
	})
}
