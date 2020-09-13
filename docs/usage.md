# Using Pfscf

## Table of Contents
* [Installation](#installation)
* [Quickstart](#quickstart)
* [Filling Out a Single Chronicle](#filling-out-a-single-chronicle)
* [Filling Out Multiple Chronicles](#filling-out-multiple-chronicles)
* [Finding the Right Chronicle Template](#finding-the-right-chronicle-template)
* [Extracting a Chronicle PDF from a Scenario](#extracting-a-chronicle-pdf-from-a-scenario)
* [Troubleshooting](#troubleshooting)

## Installation

You can download the latest release from the [Releases section](https://github.com/Blesmol/pfscf/releases). Currently there are zipped packages for Windows, macOS/Darwin and Linux. Download the release for your currently used operating system and extract into a new folder. The important things in this directory are the `pfscf` executable and the `templates` directory. The executable is, of course, the program to run this whole thing, whereas the `templates` directory contains the configuration files for the chronicle sheets, so that the program knows what to print where on the resulting sheet.

The program itself is a command line application. This means that e.g. on Windows you have to use the Windows command prompt or the Windows PowerShell to run this. You should also find a file `openCmdHere.bat` in the installation directory that will on windows open a cmd prompt window in the directory where the program is installed.

## Quickstart

1. [Download the program](https://github.com/Blesmol/pfscf/releases) and extract into a new folder. For details see the [installation instructions](#installation).
2. Get a blank society chronicle as single-paged PDF file. I'd suggest to follow the instructions in [this chapter](#extracting-a-chronicle-pdf-from-a-scenario). Put the resulting PDF page in the same directory to which you extracted `pfscf`.
3. Open a command line prompt / terminal in this directory. Can on windows be done by double-clicking the `openCmdHere.bat` file.
4. Call the program to find whether your scenario is already explicitly supported, based on the scenario title. See the example below where we want to check for PFS2 scenario "#1-06: Lost on the Spirit Road"
   ```
   $ pfscf template search lost

   Matching Templates:
   - pfs2.s1-06: #1-06: Lost on the Spirit Road
   - pfs2.s1-20: #1-20: The Lost Legend

   ```
   The important part in the test result is the template id at the beginning of the line. Here this would be `pfs2.s1-06`
5. Call the program again, to fill your first chronicle. Lets stick with the example:
   ```
   $ pfscf fill pfs2.s1-06 myBlankChronicle.pdf chronicleForBob.pdf player=Bob char="The Bobbynator" societyid=123456-2001 xp=4 gp=10
   ```
6. Open up the resulting PDF file in a PDF viewer of your choice and have a look!


## Filling Out a Single Chronicle

To fill out a chronicle, you will basically need to things to start:
1. The ID of the template to use for filling the chronicle. Details can be found in [this section](#finding-the-right-chronicle-template), but when in doubt or you just want to get this going, simply use the following:
  * `pfs2` for Pathfinder 2
  * `sfs` for Starfinder
  * Pathfinder 1 is not yet supported
2. An empty chronicle to be filled. This should be in PDF format and consist of only a single page. For information on how to create such a file if you have purchased and downloaded a scenario PDF file from Paizo, read the [section on how to extract a chronicle PDF](#extracting-a-chronicle-pdf-from-a-scenario).

Everything set so far? Good! Then we can get serious now... ok, lets do this!

To fill out a chronicle, you have to call `pfscf` with the `fill` command. The call in general looks as follows:
```
pfscf fill <template> <infile> <outfile> [<param_id>=<value> ...]
```

And here is an example on how an actual call could look like:
```
pfscf fill pfs2 s103_blank.pdf s103_bob player=Bob char="The Bobbynator" societyid=123456-2001 xp=4
```

Worked so far? Great, you've created your first filled chronicle using `pfscf`! But if something did not go as expected, weird error messages coming up or the resulting PDF looks somehow wrong, then please have a look at section [Troubleshooting](#troubleshooting).

Now you probably want to add some more stuff than just the things shown in the example above. To get the complete list of supported values for a specific chronicle, please call `pfscf template describe <template>`. This will display a list of all the supported parameters that you can use to fill out your chronicle. If you use a specialized chronicle template, e.g. template `pfs2.s1-06` from above, instead of the more generic templates like `pfs2`, you might get additional options, e.g. for striking out specific boons or other scenario-specific content.

## Filling Out Multiple Chronicles

To fill out multiple chronicles in one go, e.g. to create all chronicles for a single game session, a batch mode is included. Using this mode is (I hope) rather easy and consists of two steps that are described below.

### Create a CSV File Out of a Chronicle Template

First you have to create a CSV file (CSV: Comma-separated values) for the chronicle template that you want to use. If you do not yet know which chronicle template to use, have a look at [this section](#finding-the-right-chronicle-template).

CSV was selected here because then you can then easily use other programs like Excel or LibreOffice Calc to open this file and modify it. If you have one of the listed programs, chances are good that all you have to do is double-click on the CSV file and the correct program will open up automatically. Of course it is also possible to modify the CSV file with any texteditor of your choice.

You can create a CSV file for a specific scenario with the `pfscf batch create <template> <outputFile>` command.
```
$ pfscf batch create pfs2.s1-06 mySession.csv
```

The resulting CSV file will contain entries for all parameters supported by the selected chronicle template, like player name, society id and scenario-specific boons if they are already supported. It includes columns for up to 7 players, but you can easily add or remove columns here.

If you want to have some parameters already prefilled, you can provide additional arguments during CSV creation:
```
$ pfscf batch create pfs2.s1-06 mySession.csv event="PaizoCon" date=2020-09-12" gm="J. Doe" gmid="123456"
```

### Creating Filled Chronicles From a CSV File

So now you have already created a CSV file that contains information about your players, and now want to use that to fill out chronicles, one for each player? Great, thats what I'm talking about! For this there is the `pfscf batch fill` command, or short `pfscf b f`. The complete command with arguments is ` pfscf batch fill <template> <csv_file> <input_pdf> <output_dir> [<param_id>=<value> ...]`. An example call looks as follows:

```
$ pfscf batch fill pfs2.s1-06 mySession.csv s106_blank.pdf outputDir
Creating file outputDir\Chronicle_Player_1.pdf
Creating file outputDir\Chronicle_Player_2.pdf
Creating file outputDir\Chronicle_Player_3.pdf
Creating file outputDir\Chronicle_Player_4.pdf
Creating file outputDir\Chronicle_Player_5.pdf
Creating file outputDir\Chronicle_Player_6.pdf
Creating file outputDir\Chronicle_Player_7.pdf
```

This would then create one file per player in the specified output directory. In the example, you would have files `outputDir/Chronicle_Player_1.pdf` to `outputDir/Chronicle_Player_7.pdf`. Chronicles will only be generated if at least one value is set in the CSV file for that player.

## Finding the Right Chronicle Template

To find the right template for your chronicle, you can basically do two things: Display the complete list of supported templates, or use the builtin search function to search for a specific template

### Display list of templates

To display the complete list of supported templates, execute command `pfscf template list` (or short: `pfscf t l`)
```
$ pfscf template list
List of available templates:

- pfs2: Pathfinder 2 Society Chronicle
  - pfs2.quests: PFS2 Quests
    - pfs2.q01: Quest #01: The Sandstone Secret
	[...]
  - pfs2.s1: PFS2 Season 1: Year of the Open Road
    - pfs2.s1-00: #1-00: Origin of the Open Road
    - pfs2.s1-01: #1-01: The Absalom Initiation
	[...]
  - pfs2.specials: Specials
    - pfs2.littleTrouble: Little Trouble in Big Absalom
- sfs: Starfinder Society Chronicle
```

### Search for template

To search for a specific template, you can use the command `pfscf template search <search terms>` (or short: `pfscf t s <search terms>`). This will display all chronicle templates where all of the terms appear in the template description and template id. The search is case-insensitive.
```
$ pfscf template search star pfs2
Matching Templates:
- pfs2.s1-09: #1-09: Star-Crossed Voyages
- pfs2.s1-24: #1-24: Lightning Strikes, Stars Fall
- pfs2.s1-23: #1-23: The Star-Crossed Court
```

## Extracting a Chronicle PDF from a Scenario

### Windows

TBD using the windows PDF writer and Adobe Acrobat Reader

Short step-by-step version for the moment
1. Open scenario PDF using the Acrobat Reader
2. Switch to last page that contains the chronicle
3. Open print dialog
4. Select printer "Microsoft Print to PDF"
5. Select "Print pages: Current"
6. Select "Fit"
7. Print

### MacOS

TBD, thankful for tipps

### Linux

TBD, thankful for tipps

## Troubleshooting

* **Problem: The values in the right sidebar (xp, credits/gold, ...) are misplaced**

  **Answer:** Chronicles for different scenarios might have sidebars with a different size and thus need some minor adaptions to the template configuration. First check whether there already exists a specialized template for your scenario (`pfscf template list`) and try using that one. If this does not help, you could either [try to fix the configuration yourself](templates.md) or send me the chronicle via [email](mailto:github@pecebe.de) and I will see what I can do.

* **Problem: All values on my chronicle are misplaced**

  **Answer:** This can happen if your input chronicle PDF file has slightly different dimensions than expected. There is only so much that can be done for auto-correcting this on side of `pfscf`. The best advice that can be given at the moment is to have a look at section [Extracting a Chronicle PDF from a Scenario](#extracting-a-chronicle-pdf-from-a-scenario), follow the steps described there and see if that helps.

* **Problem:**

  **Answer:**