---
title: "pls get my issues"
permalink: pls_get_my_issues
url: /pls/pls_get_my_issues/
summary: "interact with your issues"
layout: default
---
# pls get my issues 

---
**Aliases**: i,issue

**TL;DR:** interact with your issues

## Usage:

### Examples

```

[issues in current directory's repository]: pls get my issues
[issues in a repository you own]: pls get my issues in myrepo
[issues in another's repository]: pls get my issues in organization/repo
[issues from all of github]: pls get --all my issues
```

### Local Flags

```
  -h, --help   help for issues
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
