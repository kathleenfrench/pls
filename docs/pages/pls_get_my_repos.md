---
title: "pls get my repos"
permalink: pls_get_my_repos
url: /pls/pls_get_my_repos/
summary: "interact with your github repositories"
layout: default
---
# pls get my repos 

---
**Aliases**: r,repositories,repo,repository

**TL;DR:** interact with your github repositories

## Description

`pls` makes it easy to interact with your github repositories. after fetching your repos and sorting them into a dropdown GUI for you to select from, `pls` currently supports:
- opening your default browser to the repository page in github
- cloning the repository to your chosen default codebase directory (as set in your config file), cloning it into the current directory, or choosing a custom directory
- prints a table with relevant metadata about a repository (like description, default branch, when it was created, when it was last updated, language)

## Usage:

### Examples

```
pls get my repos
```

### Local Flags

```
  -h, --help   help for repos
```

### Inherited Flags

```
      --all             search all of github
  -a, --approved        [PR] fetch only PRs that have been approved
      --assigned        [PR|ISSUE] fetch only PRs or issues assigned to you
  -x, --changesneeded   [PR] fetch only PRs where changes have been requested
  -c, --closed          [PR|ISSUE] fetch only closed prs or issues
  -b, --current         [PR] fetch the PR, if one exists, for your current working branch
  -d, --draft           [PR] fetch only draft PRs
  -l, --locked          [PR|ISSUE] fetch only locked PRs or issues
      --mention         [PR|ISSUE] fetch only PRs or issues where i've been mentioned
  -m, --merged          [PR] fetch only PRs that have been merged
  -p, --pending         [PR] get only PRs that are pending approval
  -V, --verbose         verbose output
      --viper           use viper for configuration (default true)
  -w, --work            [ALL] fetch resources via your github enterprise account
```
