# Overview

The "Pathfinder Society Chronicle Filler" (or short: pfscf) is a small command line application available for different platforms (Windows, macOS, Linux) that helps filling out chronicle sheets for the Pathfinder and Starfinder Roleplaying Games from [Paizo Inc](https://paizo.com).

You can download the program in the [Releases section](https://github.com/Blesmol/pfscf/releases). The downloaded archives should be extracted into a new directory. For details, please have a look at the [installation instructions](usage.md#installation).

## Quickstart

If you want to dive right in, have a look at the [quickstart](quickstart.md)

## Short Usage Overview

For detailed instructions on how to use this, have a look at the [Usage documentation](usage.md).

```
$ pfscf
The Pathfinder Society Chronicle Filler

Usage:
  pfscf [command]

Available Commands:
  batch       Fill out multiple chronicles in one go
  fill        Fill out a single chronicle sheet
  help        Help about any command
  template    Various actions on templates: list, describe, etc

Flags:
  -h, --help      help for pfscf
  -v, --verbose   verbose output

Use "pfscf [command] --help" for more information about a command.
```

### Fill out a single chronicle

Filling out a single Pathfinder/Starfinder Society chronicle sheet can be done with the `fill` subcommand.

#### Example Call
```
$ pfscf fill pfs2 s103_blank_chronicle.pdf s103_for_bob.pdf player=Bob char="The Bobbynator" xp=4
```
