package completion

// BashHelp prints for users who elect to install their bash completion scripts manually
const BashHelp = `
Bash:

$ source <(pls add complete bash)

# To load completions for each session, execute once:
Linux:
  $ pls add complete bash > /etc/bash_completion.d/pls
MacOS:
  $ pls add complete bash > /usr/local/etc/bash_completion.d/pls
`
