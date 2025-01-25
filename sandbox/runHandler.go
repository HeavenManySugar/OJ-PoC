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

func (s *Sandbox) RunShellCommand(shellCommand []byte, codePath []byte) string {
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
	cmdArgs := []string{
		fmt.Sprintf("--box-id=%v", boxID),
		"--fsize=5120",
		fmt.Sprintf("--dir=%v", CodeStorageFolder),
		"--wait",
		"--processes=100",
		"--open-files=0",
		"--env=PATH",
		"--stderr-to-stdout",
	}

	if len(codePath) > 0 {
		cmdArgs = append(cmdArgs, fmt.Sprintf("--dir=%v", string(codePath)), fmt.Sprintf("--env=CODE_PATH=%v", string(codePath)))
	}

	cmdArgs = append(cmdArgs, "--run", "--", "/usr/bin/sh", shellFilename(codeID))


	cmd := exec.CommandContext(ctx, "isolate", cmdArgs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to run command: %v", err)
		return fmt.Sprintf("Failed to run command: %v", err)
	}

	log.Printf("Command output: %s", string(out))
	return string(out)
}