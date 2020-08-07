package clean

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kathleenfrench/pls/pkg/gui"
	"github.com/kyokomi/emoji"
)

// SystemPrune prunes the docker system
func SystemPrune() error {
	var allImagesIcon string
	var volumesIcon string

	allUnusedImages := gui.ConfirmPrompt("do you want to prune all unused images, not just dangling ones?", "", false, false)
	pruneVolumes := gui.ConfirmPrompt("do you want to prune all local volumes?", "this will delete local data - make sure you know what you're doing!", false, false)
	cmd := "docker system prune --force"

	if allUnusedImages {
		cmd += " --all"
		allImagesIcon = ":white_check_mark:"
	} else {
		allImagesIcon = ":x:"
	}

	if pruneVolumes {
		cmd += " --volumes"
		volumesIcon = ":white_check_mark:"
	} else {
		volumesIcon = ":x:"
	}

	gui.Log(":broom:", "pruning local docker resources", fmt.Sprintf("include all images: %s | include volumes: %s", emoji.Sprint(allImagesIcon), emoji.Sprint(volumesIcon)))

	cleanup := exec.Command("bash", "-c", cmd)
	cleanup.Stdout = os.Stdout
	cleanup.Stdin = os.Stdin
	cleanup.Stderr = os.Stderr

	if err := cleanup.Run(); err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout)

	// TODO: find out how much space was saved?

	gui.Log(":thumbs_up:", "success!", nil)
	return nil
}
