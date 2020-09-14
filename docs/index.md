# Overview

The "Pathfinder Society Chronicle Filler" (or short: pfscf) is a small command line application available for different platforms (Windows, macOS, Linux) that helps filling out chronicle sheets for the Pathfinder and Starfinder Roleplaying Games from [Paizo Inc](https://paizo.com).

You can download the program in the [Releases section](https://github.com/Blesmol/pfscf/releases). The downloaded archives should be extracted into a new directory. For details, please have a look at the [installation instructions](usage.md#installation).

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

### Writing your own templates

Please see the [Documentation on writing your own templates](docs/templates.md). **This documentation is currently outdated and needs some serious rework! In case you are still interested, have a look at the existing templates in the `templates/` subdir**

## Legal disclaimer

Pathfinder, the Pathfinder logo, the Pathfinder Society, Starfinder, the Starfinder Society and the Starfinder logo are registered trademarks of Paizo Inc. Their games use the [Open Gaming License](https://paizo.com/pathfinder/compatibility/ogl) and are pretty cool. Support their games!

This program is being developed as private hobby. Although this program is intended to be used with chronicles for organized play from Paizo, I am in no way associated with Paizo Inc. Also, you're using this program at your own risk. I won't take any responsibility or liability for any direct or indirect damage, data loss, data corruption and the like done by using this program. I cannot guarantee that the program is free of bugs (on the contrary, I am pretty sure that there is a sufficient number of bugs still included), and if I am made aware of a problem, I won't guarantee on a timeline until when that bug will be fixed.
