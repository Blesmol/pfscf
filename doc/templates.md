# Creating Template Files for Pfscf

After some time of using pfscf, you might be tempted to create your own templates or modify existing ones.
Or you are just curious on what is currently supported by pfscf.
Besides, it's always useful to provide proper documentation.
Even if perhaps only for the reason to show that the weird behavior for something might be actually intended!

## Template Structure on Disk

The template files are stored in a folder `templates` that is located in the same folder as the pfscf executable file.
Within this folder, each template is stored in a single file with file extension `.yml`.
It is also possible to create subfolder within the `templates` folder to organize the template files.
For example, this could be used to have separate folders for each supported game system, e.g. one folder for Pathfinder, one for Starfinder.

The filenames of the template files are irrelevant.
They should provide a hint to what is included, but besides from that they are not used anywhere.
This especially means that they are not used for the unique identifier which each template must have; this is done by the `id` field within the template file.
See also the `Template file format` section below on this.

## Template File Format

### YAML Format

Template files are stored in YAML format, one template per file.
The official spec of this format is maintained at [yaml.org](https://yaml.org/), and if you search the web you will find lots of examples, introductions and explanations to the format.
In this document, however, I will try to spare you the irrelevant parts and only include a basic overview over what you might need.

A YAML file is a plain text file and can be opened with any text editor, if need be with the Windows Notepad application.
What is not possible is to use MS Word, LibreOffice and the like to modify such files.
Many modern text editors even bringt support for writing/modifying YAML files.
One such editor that I can recommend it [Notepad++](https://notepad-plus-plus.org) (not to be mixed up with the already mentioned MS Notepad).

YAML is intended to be a human-readable format.
Most lines in these files follow the format `<some identifier>: <some value>`, e.g. `description: This is a template for PFS`.
And until I have more time to write a proper YAML introduction, I would suggest you take a look at the existing files in the `templates` folder or search the web.

### File Layout

Template files have the following structure:
```yaml
id: <unique template identifier>
description: <short description on what this template contains>
inherit: <id of template that should be inherited>

presets:
  <presetId>:
    <content>
  <presetId>:
    <content>
  ...

content:
  <contentId>:
    <content>
  <contentId>:
    <content>
  ...
```

<details>
  <summary>Example</summary>

```yaml
id: myId
description: This is an example template
inherit: idOfMyOtherCoolTemplate

presets:
  topline:
    y1: 100
    y2: 200
  bottomline:
    y1: 400
	y2: 500

content:
  name:
    type: textCell
	presets: [ topline ]
```
</details>

The only mandatory top-level field in such a template is the `id`, all other top-level fields (`description`, `inherit`, `presets`, `content`) in this basic structure are optional.
Of course, a template that only consists of an `id`  does not make much sense.
But who am I to judge?

#### Example

### Preset Mechanism

### Content Types

## Template Inheritance
