# Using Pfscf

## Table of Contents
* [Filling Out a Single Chronicle](#filling-out-a-single-chronicle)
* [Finding the Right Chronicle Template](#finding-the-right-chronicle-template)
* [Extracting a Chronicle PDF from a Scenario](#extracting-a-chronicle-pdf-from-a-scenario)
* [Troubleshooting](#troubleshooting)

## Filling Out a Single Chronicle

To fill out a chronicle, you will basically need to things to start:
1. The ID of the template to use for filling the chronicle. Details can be found in [this section](#finding-the-right-chronicle-template), but when in doubt or you just want to get this going, simply use the following:
  * `pfs2` for Pathfinder 2
  * `sfs` for Starfinder
  * Pathfinder 1 is not yet supported
2. An empty chronicle to be filled. This should be in PDF format and contain only a single page. For information on how to create such a file if you have purchased and downloaded a scenario PDF file from Paizo, read the [section on how to extract a chronicle PDF](#extracting-a-chronicle-pdf-from-a-scenario).

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

Now you probably want to add some more stuff than just the things shown in the example above. To get the complete list of supported values for a specific chronicle, please call `pfscf template describe <template>`


## Finding the Right Chronicle Template

Short version: Call `pfscf template list` and have a look at the resulting list

TBD

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
