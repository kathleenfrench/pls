package git

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/kathleenfrench/pls/pkg/utils"
)

// CleanupCurrentBranches culls branches that have already been merged from the local working environment
func CleanupCurrentBranches() error {
	cmd := fmt.Sprintf("git branch --merged | grep -v '%s' | grep -v '%s' | grep -v '%s' | xargs -n 1 git branch -d", "\\master", "\\development", "\\*")

	eligible, err := utils.BashExec(cmd)
	if err != nil {
		alreadyMergedList, _ := utils.BashExec("git branch --merged | grep -v master")
		color.HiYellow("There was an error deleting one of your branches, likely because it has not been fully merged - refer to the list below for your already merged branches to find the offending branch and force delete (via git branch -D <branchname> if desired.")
		color.HiBlue(fmt.Sprintf("\n%s", alreadyMergedList))
		return fmt.Errorf("%s: could not cleanup branches", err)
	}

	if len(string(eligible)) == 0 {
		color.HiBlue("No unmerged local branches found!")

	} else {
		color.HiGreen(string(eligible))
	}

	return nil
}
