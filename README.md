# The Pathfinder Society Chronicle Tagger (pfsct)

The "Pathfinder Society Chronicle Tagger" (or short: pfsct) is a small
command line application available for different platforms (Windows,
macOS, Linux) that helps filling out chronicle sheets for the
Pathfinder Roleplaying Game.

## Usage

```
$ pfsct
The Pathfinder Society Chronicle Tagger

Usage:
  pfsct [command]

Available Commands:
  fill        Fill a single chronicle sheet
  help        Help about any command

Flags:
  -h, --help   help for pfsct

Use "pfsct [command] --help" for more information about a command.
```

### Fill out a single chronicle

Fill me

### Fill out multiple chronicles (e.g. for a complete group)

Fill me

## FAQ

* **Q: What types of chronicles are supported at the moment?**

  **A:** At the moment only the chronicles for Pathfinder 2. In the near future most likely also Starfinder and Pathfinder 1. Although the latter will be more complicated, as Paizo seems to have changed the chronicle layout between seaons.

* **Q: I have a Society scenario and want to use pfsct on that. But it keeps complaining that some operation is not allowed because of PDF permissions on that file. What should I do?**

  **A:** The permission settings in place on Paizos PDFs do not allow to do things like extracting pages. Which I can totally understand for the scenario and their property in general, but this makes life a little bit harder. On Windows you can open the PDF file of the scenario in a PDF viewer of your choice, go to the last page with the chronicle, and then print this using the "Microsoft Print to PDF" printer. This will produce a PDF file on your local disk that you can use together with pfsct. I assume that for macOS and linux similar options exist.

* **Q: What about scenario-specific options? For example, I have this one scenario where I need an easy way to strike out boons that the group did not get.**

  **A:** Actually, that is on my roadmap and I have some ideas for this. However, neither do I own all Society scenarios nor do I have enough time to provide the config for multiple dozens of scenarios that I own. But... contributions are welcome! When this feature is available, then no programming skills will be required for providing a scenario-specific config. Just a text editor and some time to fiddle out the coordinates for all the relevant parts. And what do I also have on the roadmap? An auto-updater for such configs! So if people would contribute scenario-specific configs, then other people don't need to update their installation of pfsct every two weeks, but just call the program with some yet-to-be-determined parameter (think of `pfsct update-config`) and you get the latest and greatest scenario-specific config automatically downloaded to your computer.

* **Q: Oh noes! I have found a bug! I HAVE FOUND A BUG!!1!**

  **A:** First, calm down. Second, please report this to me by opening an [issue](https://github.com/Blesmol/pfsct/issues). If you do so, please provide as many details as possible. What operating system are you running? What did you do (exact command line, input files), what did then happen (exact output)? I will try to find out what can be done about this and try to fix it as soon as possible.

* **Q: I have this absolutely great idea for a missing feature!**

  **A:** Yay! Feature ideas are always welcome! Please open an [issue](https://github.com/Blesmol/pfsct/issues) in this repository and describe your idea! The more details, the better! But please be also aware that I am doing this in my free time, and some things might be really, really complicated to realize. So please don't be mad if things might take a little bit longer or sometimes even won't be done. But... it's open source, contributions are always possible and welcome :) Still: If you have an idea that you think would improve this project, please create an issue here. You'll never know unless you tried!

* **Q: Will there ever be a GUI for pfsct?**

  **A:** Not by me. I suck at GUIs and am also not interested in having one for pfsct. But if you're interested in making one and using pfsct in the background, then I'm sure that lots of people will love it! Honestly! It's just not something for me, sorry.

* **Q: Why golang?**

  **A:** It produces standalone executables for all major platforms (Windows, macOS, Linux) and users do not have to install any additional software to use pfsct. Also I wanted to give golang a try anyways (first go project for me), and then it's always good to have something that you plan to use yourself.

## Legal stuff

Pathfinder, the Pathfinder logo, Pathfinder Society, Starfinder and the Starfinder logo are registered trademarks of Paizo Inc.
