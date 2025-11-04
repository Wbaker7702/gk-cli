package patch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// CreatePatch creates a git patch from current changes
func CreatePatch(name, description string) ([]byte, error) {
	// Check if we're in a git repo
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return nil, fmt.Errorf("not in a git repository")
	}

	// Create patch from current changes
	cmd := exec.Command("git", "diff", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to create patch: %w", err)
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("no changes to create patch from")
	}

	return output, nil
}

// CreatePatchFromCommit creates a patch from a commit
func CreatePatchFromCommit(commitHash string) ([]byte, error) {
	cmd := exec.Command("git", "format-patch", "-1", commitHash, "--stdout")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to create patch from commit: %w", err)
	}

	return output, nil
}

// ApplyPatch applies a patch file
func ApplyPatch(patchData []byte, targetDir string) error {
	if targetDir == "" {
		targetDir = "."
	}

	// Write patch to temp file
	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("patch-%d.patch", time.Now().Unix()))
	if err := os.WriteFile(tmpFile, patchData, 0644); err != nil {
		return fmt.Errorf("failed to write patch file: %w", err)
	}
	defer os.Remove(tmpFile)

	// Apply patch
	cmd := exec.Command("git", "apply", tmpFile)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply patch: %w", err)
	}

	return nil
}

// ViewPatch displays patch contents
func ViewPatch(patchData []byte) {
	fmt.Println(string(patchData))
}
