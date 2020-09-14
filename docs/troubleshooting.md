# Troubleshooting

* **Problem: The values in the right sidebar (xp, credits/gold, ...) are misplaced**

  **Answer:** Chronicles for different scenarios might have sidebars with a different size and thus need some minor adaptions to the template configuration. First check whether there already exists a specialized template for your scenario (`pfscf template list`) and try using that one. If this does not help, you could either [try to fix the configuration yourself](templates.md) or send me the chronicle via [email](mailto:github@pecebe.de) and I will see what I can do.

* **Problem: All values on my chronicle are misplaced**

  **Answer:** This can happen if your input chronicle PDF file has slightly different dimensions than expected. There is only so much that can be done for auto-correcting this on side of `pfscf`. The best advice that can be given at the moment is to have a look at section [Getting Blank Chronicle Sheets](extraction.md), follow the steps described there and see if that helps.

* **Problem: The program keeps complaining that some operation is now allowed because of PDF permissions on my file**

  **Answer:** The permission settings in place on Paizos PDFs do not allow to do things like extracting pages. Which I can totally understand for the scenario and their property in general, but this makes life a little bit harder. Please see the section on [Getting Blank Chronicle Sheets](extraction.md) for guidance.

* **Problem: Oh noes! I have found a bug! I HAVE FOUND A BUG!!1!**

  **Answer:** First, calm down. Second, please report this to me by opening an [issue](https://github.com/Blesmol/pfscf/issues). If you do so, please provide as many details as possible. What operating system are you running? What did you do (exact command line, input files), what did then happen (exact output)? I will try to find out what can be done about this and try to fix it as soon as possible.

<!--
* **Problem:**

  **Answer:**
-->
