---
title: "pls make a pullrequest"
permalink: pls_make_a_pullrequest
url: /pls/pls_make_a_pullrequest/
summary: "let `pls` open a pull request for you based off of the branch in your current working directory"
layout: default
---
# pls make a pullrequest 

---
**Aliases**: pr,pull

**TL;DR:** let `pls` open a pull request for you based off of the branch in your current working directory

## Description

`pls` will open a pull request for you, but only after verifying that your current local branch is synced with its remote ref. if it isn't synced, then `pls` will confirm whether or not you want to let `pls` handle adding, committing, and/or pushing the branch for you. once that's done, `pls` will prompt you for various values:
- title
- PR description
- whether to link it to an existing issue
- whether to create it as a draft (if using an enterprise account)

`pls` will even spawn a temporary file in your preferred text editor for composing the body of the pull request. once you finish adding values, `pls` handles creating the PR!

## Usage:

### Examples

```
pls make a pr
```

### Local Flags

```
  -h, --help   help for pullrequest
```

### Inherited Flags

```
  -V, --verbose   verbose output
      --viper     use viper for configuration (default true)
```
