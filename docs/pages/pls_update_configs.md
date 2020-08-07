---
title: "pls update configs"
permalink: pls_update_configs
url: /pls/pls_update_configs/
summary: "update your pls configs"
layout: default
---
# pls update configs 

---
**Aliases**: config,cnf,cnfs

**TL;DR:** update your pls configs

## Description

`pls` was written to make life easier, and flexibility to change configuration valus is a key component of that. you can either change a value through the interactive GUI (via `pls update configs`), or if you already know the key and value to set, you can invoke it directly (via `pls update configs --raw <key> <value>`)

## Usage:

### Local Flags

```
  -h, --help   help for configs
  -r, --raw    input as key value, skip the dropdown GUI
```

### Inherited Flags

```
  -V, --verbose   verbose output
      --viper     use viper for configuration (default true)
```
