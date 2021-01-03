# Quickstart

This chapter presents quickstart instructions for filling out multiple chronicle sheets for a complete game group. You should have the scenario as PDF file at hand.

## Quickstart

1. [Download the program](https://github.com/Blesmol/pfscf/releases) and extract into a new folder. For details see the [installation instructions](installation.md).
2. Get a blank society chronicle as single-paged PDF file. I'd suggest to follow [these instructions](extraction.md). Put the resulting PDF page in the same directory to which you extracted `pfscf`.
3. Open a command line prompt / terminal in this directory. On windows this can be done by double-clicking the `openCmdHere.bat` file included in the installation folder.
4. Check whether the scenario from your chronicle is already explicitly supported. Call `pfscf template list` to get a complete list of all supported scenarios, or use `pfscf template search <part of scenario title>` to search for it. Write down the ID from the beginning of the line, e.g. something like `pfs2.s1-06` for PFS2 scenario "#1-06: Lost on the Spirit Road".
6. Create a CSV file for your scenario where you collect all the required information about the players and the event. Do this by calling `pfscf batch create <filename> -t <id from step 4>`. This could look something like `pfscf batch create session.csv -t pfs2.s1-06`.
7. If you have a program like Excel or LibreOffice installed on your PC, then chances are good that this will now automatically open up and present the contents of the CSV file. If not, you have to open the CSV file manually, either with one of the named programs or with a text editor like [Notepad++](https://notepad-plus-plus.org).
8. Fill in the values for each player and then save the file.
9. Finally create the chronicle files by calling `pfscf batch fill <csv file from step 6> -i <blank chronicle from step 2> -o <output directory>`. For our example, this could look like the following: `pfscf batch fill session.csv -i blank-s1-06.pdf -o outputDir`. This should now create a single PDF file per player in the listed dir. If the specified directory does not yet exist, `pfscf` will try to create it.

## Basic Troubleshooting

- If something does not work as described above, please check whether you followed the exact instructions.
- If chronicle files were created, but all the printed values seem to be slightly off, please try to create a new blank chronicle with the instructions from [this chapter](extraction.md).
- Check whether the [troubleshooting chapter](troubleshooting.md) provides additional help.

