package completion

import "os/exec"

// FishHelp prints to users who elect to manually install fish
const FishHelp = `
Fish:

$ pls add complete fish | source

# To load completions for each session, execute once:
$ pls add complete fish > ~/.config/fish/completions/pls.fish
`

var fishInit = exec.Command("pls", "add", "complete", "fish", "|", "source")
