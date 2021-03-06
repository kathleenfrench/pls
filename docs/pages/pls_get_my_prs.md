---
title: "pls get my prs"
permalink: pls_get_my_prs
url: /pls/pls_get_my_prs/
summary: "interact with your pull requests"
layout: default
---
# pls get my prs 

---
**Aliases**: pulls,pull,pr

**TL;DR:** interact with your pull requests

## Description

when `pls` fetches your PRs, you will be greeted with a straightforward dropdown to select the one you want to do something with. currently, `pls` supports viewing the PR description in the terminal with rendered markdown, editing the PR title, body, and/or state, merging a PR, and opening it in your default browser.

**on merging:** `pls` makes merging a breeze, precisely because you can trust you're not forgetting anything. `pls` makes sure you don't have any unstage, uncommitted, and/or unpushed code to your remote branch before initiating a merge. after your code is merged successfully, `pls` checks you back into `master`, pulls down the latest code, and removes already-merged branches from your local machine. easy!

## Usage:

### Examples

```
[PRs in current directory's repository]: pls get my prs
[PRs in a repository you own]: pls get my prs in myrepo
[PRs in another's repository]: pls get my prs in organization/repo
[PRs from all of github]: pls get --all my prs
[PRs on my work account]: pls get my --work prs
```

### Local Flags

```
  -h, --help   help for prs
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
### Sub Commands

* [pls get my prs where](/pls/pls_get_my_prs_where)	 - 

