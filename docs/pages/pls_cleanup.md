---
title: "pls cleanup"
permalink: pls_cleanup
url: /pages/pls_cleanup/
summary: "cleanup subcommands"
layout: default
---
# pls cleanup 

---
**Aliases**: c,clean

**Purpose:**

cleanup subcommands

### Local Flags

```
  -h, --help   help for cleanup
```

### Inherited Flags

```
  -V, --verbose   verbose output
      --viper     use viper for configuration (default true)
```
### Sub Commands

* [pls cleanup docker](/pages/pls_cleanup_docker/)	 - prune local docker resources to free up space
* [pls cleanup git](/pages/pls_cleanup_git/)	 - remove git branches that have already been merged into master - defaults to just within the current working directory

### See Also

* [pls](/pages/pls/)	 - a helpful little CLI that does things for you when you ask nice...