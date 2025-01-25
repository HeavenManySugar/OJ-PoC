package sandbox

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

const execTimeoutDuration = time.Second * 60

func (s *Sandbox) RunShellCommand(shellCommand []byte) string {
	ctx, cancel := context.WithTimeout(context.Background(), execTimeoutDuration)
	defer cancel()

	boxID := s.Reserve()
	defer s.Release(boxID)

	// saving code as file
	codeID, err := WriteToTempFile([]byte(shellCommand))
	if err != nil {
		log.Println("error saving code as file:", err)
		return "Failed to save code as file"
	}
	defer os.Remove(shellFilename(codeID))

	// running the code
	cmd := exec.CommandContext(ctx,
		"isolate",
		fmt.Sprintf("--box-id=%v", boxID), 
		// max size (in KB) of files that can be created per execution = 5MB
		"--fsize=5120",
		// makes directory visible in the sandbox
		fmt.Sprintf("--dir=%v", CodeStorageFolder),
		// if sandbox is busy, wait instead of returning error right away
		// instead of serving 25/100 requests in 10 sandbox, it's gonna serve all
		"--wait",
		// to keep the child process in parent’s network namespace and communicate with the outside world
		// "--share-net",
		"--processes=100",
		// unlimited open files
		"--open-files=0",
		// makes commands visible in the sandbox e.g. 'ls', 'echo' or other installed command
		"--env=PATH",
		// log package writes to stderr instead of stdout, so we need to redirect this to stdout.
		// only exit code determines if the program ran successfully or not
		"--stderr-to-stdout",
		"--run", 
		"--", 
		"/usr/bin/sh", 
		shellFilename(codeID))

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to run command: %v", err)
		return fmt.Sprintf("Failed to run command: %v", err)
	}

	log.Printf("Command output: %s", string(out))
	return string(out)
}