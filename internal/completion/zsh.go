package completion

// ZshHelp is printed for users who elect to manually install the completion scripts
const ZshHelp = `
Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ pls add complete zsh > "${fpath[1]}/_pls"

# You will need to start a new shell for this setup to take effect.
`
