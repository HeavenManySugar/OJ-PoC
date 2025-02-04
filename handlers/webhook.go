package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/gofiber/fiber/v2"

	"github.com/HeavenManySugar/OJ-PoC/database"
	"github.com/HeavenManySugar/OJ-PoC/models"
	"github.com/HeavenManySugar/OJ-PoC/sandbox"
)

const GitServer = "http://server.gitea.orb.local/"
const RepoFolder = "/sandbox/repo"

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
//	@Tags			WebHook
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
	r, err := git.PlainClone(fmt.Sprintf("%s/%s", RepoFolder, payload.Repository.FullName), false, &git.CloneOptions{
		URL:      GitServer + payload.Repository.FullName,
		Progress: os.Stdout,
	})
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to clone repository",
		})
	}
	os.Chmod(fmt.Sprintf("%s/%s", RepoFolder, payload.Repository.FullName), 0777)
	defer os.RemoveAll(fmt.Sprintf("%s/%s", RepoFolder, payload.Repository.FullName))
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
	
	sandbox.SandboxPtr.RunShellCommandByRepo(payload.Repository.Parent.FullName, []byte(fmt.Sprintf("%s/%s", RepoFolder, payload.Repository.FullName)))

	// read score from file
	score, err := os.ReadFile(fmt.Sprintf("%s/%s/score.txt", RepoFolder, payload.Repository.FullName))
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to read score",
		})
	}
	log.Printf("Score: %s", score)
	// save score to database
	db := database.DBConn
	scoreFloat, err := strconv.ParseFloat(strings.TrimSpace(string(score)), 64)
	if err != nil {
		log.Printf("Failed to convert score to int: %v", err)
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to convert score to int",
		})
	}
	
	// read message from file
	message, err := os.ReadFile(fmt.Sprintf("%s/%s/message.txt", RepoFolder, payload.Repository.FullName))
	if err != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to read message",
		})
	}
	log.Printf("Message: %s", message)

	// Create a new score entry in the database
	newScore := models.Score{
		GitRepo: payload.Repository.FullName,
		Score:   scoreFloat,
		Message: strings.TrimSpace(string(message)),
	}
	
	if err := db.Create(&newScore).Error; err != nil {
		log.Printf("Failed to create new score entry: %v", err)
		return c.Status(http.StatusServiceUnavailable).JSON(ResponseHTTP{
			Success: false,
			Message: "Failed to create new score entry",
		})
	}
	return c.JSON(ResponseHTTP{
		Success: true,
		Message: "Successfully received hook",
		Data:    payload,
	})
}