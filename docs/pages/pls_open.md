---
title: "pls open"
permalink: pls_open
url: /pls/pls_open/
summary: "open any url in your default browser from the command line, or select from a set of common favorites"
layout: default
---
# pls open 

---
**Aliases**: o

**TL;DR:** open any url in your default browser from the command line, or select from a set of common favorites

## Description

`pls` comes with a few pre-baked url shortcuts, but the rest is up to you. view your configs (`pls show configs`) to see what shortcuts have already been set to the `webshort` property. if you ever want to update these values - whether that be changing an existing url or adding a new shortcut - simply run `pls update configs` and follow the onscreen prompts.

## Usage:

### Examples

```
pls open [opens dropdown GUI of your url shortcuts]
pls open https://google.com
pls open google.com
```

### Local Flags

```
  -h, --help   help for open
```

### Inherited Flags

```
  -V, --verbose   verbose output
      --viper     use viper for configuration (default true)
```
### See Also

* [pls](/pls/pls)	 - a helpful little CLI that does things for you when you ask nice...
