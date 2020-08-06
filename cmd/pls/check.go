package pls

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/google/go-github/v32/github"
	"github.com/kathleenfrench/pls/pkg/utils"
	"github.com/kathleenfrench/pls/pkg/web/git"
	"github.com/spf13/cobra"
)

// CURRENTLY NOT IN USE

var checkCmd = &cobra.Command{
	Use:     "check",
	Short:   "run a check on a given resource",
	Aliases: []string{"ck"},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		gc := git.NewClient(ctx, plsCfg.GitToken)
		user, _, err := gc.Users.Get(ctx, "")
		if err != nil {
			utils.ExitWithError(err)
		}

		color.HiYellow(fmt.Sprintf("git user: %s", github.Stringify(user)))
	},
	Hidden: true,
}
