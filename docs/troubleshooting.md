# Troubleshooting

??? note "The values in the right sidebar (xp, credits/gold, ...) are misplaced"

	Chronicles for different scenarios might have sidebars with a different size and thus need some minor adaptions to the template configuration.
	First check whether there already exists a specialized template for your scenario (`pfscf template list`) and try using that one.
	If this does not help, you could either [try to fix the configuration yourself](templates.md) or send me the chronicle via [email](mailto:github@pecebe.de) and I will see what I can do.

??? note "All values on my chronicle are misplaced"

	This can happen if your input chronicle PDF file has slightly different dimensions than expected.
	There is only so much that can be done for automatically correcting this on side of `pfscf`.
	To fix this, I would propose to first have a look at section [Getting Blank Chronicle Sheets](extraction.md) and follow the steps described there to see if that helps.
	If this does not help, lets see if the following does:

	First, execute the following command:

	``` bash
	$ pfscf fill <your_scenario_id> <your_blank_chronicle.pdf> output.pdf -d -e
	```

	Now you should have file `output.pdf` with some example value and some green boxes.
	There is one large box in the middle that covers nearly all of the content.
	In the lower right corner of this green box it says `main` in green letters.
	If everything would be correct, then it would be overlaying the edges of what I call the "main content area".
	If this is off for your chronicle, then you have to correct this for your runs.
	You can do this by adding an additional parameter `-offset-x <value>` (or short: `-x <value>` to your call.
	If the green box is larger than the main area, use positive values, e.g. `-x 10`.
	If the green box is smaller than the main area, try negative values, e.g. `-x -10`.
	Might take some trying around until you got a working value, but finally it should hopefully work.
	Then you can use that same value when filling out your chronicle.
	This parameter is supported both for the `fill` and for the `batch fill` command.

??? note "I have a Society scenario and want to use pfscf on that. But it keeps complaining that some operation is not allowed because of PDF permissions on that file."

	The permission settings in place on Paizos PDFs do not allow to do things like extracting pages.
	Which I can totally understand for the scenario and their property in general, but this makes life a little bit harder.
	Have a look at the chapter about [getting blank chronicle sheets](extraction.md) for instructions on how to extract a chronicle sheet from your scenario PDF file.

??? note "Oh noes! I have found a bug! I HAVE FOUND A BUG!!1!"

	First, calm down.
	Second, please report this to me by opening an [issue](https://github.com/Blesmol/pfscf/issues).
	If you do so, please provide as many details as possible.
	What operating system are you running?
	What did you do (exact command line, input files), what did then happen (exact output)?
	I will try to find out what can be done about this and try to fix it as soon as possible.
