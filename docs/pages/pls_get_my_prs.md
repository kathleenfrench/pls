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

**Purpose:**

interact with your pull requests

## Usage:

### Examples

```

[PRs in current directory's repository]: pls get my prs
[PRs in a repository you own]: pls get my prs in myrepo
[PRs in another's repository]: pls get my prs in organization/repo
[PRs from all of github]: pls get --all my prs
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

* [pls get my prs where](/pls/pls_get_my_prs_where/)	 - 

