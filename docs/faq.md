# Frequently Asked Questions

!!! info ""

	This is the section about general questions. In case you are having problems with using the program, please have a look at the [troubleshooting chapter](troubleshooting.md)

??? question "What types of chronicles / game systems are supported at the moment?"

	At the moment the chronicles for Pathfinder 2 and Starfinder are supported.
	Pathfinder 1 is in general also possible, but I did not yet have time to provide the proper configuration.
	Also Pathfinder 1 may take some more time, as Paizo seems to have changed the chronicle layout between seaons.

??? question "What about scenario-specific options? For example, I have this one scenario where I need an easy way to strike out boons that the group did not get."

	The program already supports in general to have scenario-specific options, and there are already templates for some scenarios.
	However, neither do I own all Society scenarios nor do I have enough time to provide the config for multiple dozens of scenarios that I own.
	But... contributions are welcome!
	No programming skills are required for providing a scenario-specific config.
	Just a text editor and some time to fiddle out the coordinates for all the relevant parts.
	And what do I also have on the roadmap?
	An auto-updater for such configs!
	So if people would contribute scenario-specific configs, then other people don't need to update their installation of pfscf every two weeks, but just call the program with some yet-to-be-determined parameter (think of `pfscf update-config`) and you get the latest and greatest scenario-specific config automatically downloaded to your computer.

??? question "I have this absolutely great idea for a missing feature! Can you implement this?"

	Yay! Feature ideas are always welcome!
	Please open an [issue](https://github.com/Blesmol/pfscf/issues) in this repository and describe your idea!
	The more details, the better!
	But please be also aware that I am doing this in my free time, and some things might be really, really complicated to realize.
	So please don't be mad if things might take a little bit longer or sometimes even won't be done.
	But... it's open source, contributions are always possible and welcome :)
	Still: If you have an idea that you think would improve this project, please create an issue here.
	You'll never know unless you tried!


??? question "Will there ever be a GUI for pfscf?"

	Quite unlikely, at least not by me.
	I suck at GUIs and am also not that interested in having one for pfscf (although I also see the advantages).
	But if you're interested in making one and using pfscf in the background, then I'm sure that lots of people will love it!
	Honestly! It's just not something for me, sorry.

??? question "Why golang?"

	It produces standalone executables for all major platforms (Windows, macOS, Linux) and users do not have to install any additional software to use pfscf.
	Also I wanted to give golang a try anyways (first go project for me), and then it's always good to have something that you plan to use yourself.
