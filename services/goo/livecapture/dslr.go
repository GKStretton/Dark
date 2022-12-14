package livecapture

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/gkstretton/dark/services/goo/filesystem"
)

// After testing it seems like this isn't necessary
func setDslrState(b bool) {
	state := "off"
	if b {
		state = "on"
	}
	cmd := exec.Command("./scripts/set-dslr-state.sh", state)
	err := cmd.Run()
	if err != nil {
		fmt.Printf("failed to run set-dslr-state: %v\n", err)
	}
}

func captureSessionImage(sessionId uint64) {
	p := filesystem.GetIncrementalFileName(filesystem.GetRawDslrDir(sessionId), "jpg")
	err := captureImage(p)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = saveCropConfig(CC_DSLR, p)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = processImage(p, filesystem.GetPostDslrDir(sessionId))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func captureImage(p string) error {
	captureCmd := exec.Command("./scripts/capture-dslr.sh", p)
	// cmd.Stdout = os.Stdout
	captureCmd.Stderr = os.Stderr

	err := captureCmd.Run()
	if err != nil {
		return fmt.Errorf("failed to run capture-dslr: %v", err)
	}

	return nil
}

func processImage(imgPath, outDir string) error {
	postCmd := exec.Command("python", "../../user-tools/auto_image_post.py",
		"-i", imgPath, "-o", outDir,
	)
	postCmd.Stderr = os.Stdout
	// postCmd.Stdout = os.Stdout
	err := postCmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start imagePost cmd: %v", err)
	}
	return nil
}
