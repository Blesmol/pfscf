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
    x1: 50
    x2: 150
  date:
    type: textCell
    presets: [ bottomline ]
    x1: 200
    x2: 230
```
</details>

The only mandatory top-level field in such a template is the `id`, all other top-level fields (`description`, `inherit`, `presets`, `content`) in this basic structure are optional.
Of course, a template that only consists of an `id`  does not make much sense.
But who am I to judge?

Each entry below `content` and `presets` must have an ID that is unique within that section.
That means, no two entries in the `content` section may have the same ID, and no two entries in the `presets` section may have the same ID.
The description on what may be included in a `content` or `presets` entry is described below.

### Content Types

The pfsct app supports different types of content entries.
Each content entry requires a mandatory field `type` where the type of the content entry must be added.

#### `textCell`

A `textCell` describes a rectangular cell on the PDF file where user-provided text is added.

It has a couple of mandatory and some optional fields:
| Field      | Description                                                       | Input type | Required? |
|:-----------|:------------------------------------------------------------------|:----------:|:---------:|
| `desc`     | Gives a short description of what this is                         | Text       | TODO      |
| `x1`       | First coordinate for the cell on the X axis                       | Number     | Mandatory |
| `y1`       | First coordinate for the cell on the Y axis                       | Number     | Mandatory |
| `x2`       | Second coordinate for the cell on the X axis                      | Number     | Mandatory |
| `y2`       | Second coordinate for the cell on the Y axis                      | Number     | Mandatory |
| `font`     | Name of the font to use. See [here](#fonts)                       | Text       | Mandatory |
| `fontsize` | Fontsize in points                                                | Number     | Mandatory |
| `align`    | Text alignment inside the rectangle. See [here](#text-alignment). | Text       | TODO      |
| `example`  | An example input value                                            | Text       | Optional  |


Regarding the coordinates: It does not matter whether `x1` is smaller than `x2` or vice versa.
Same goes for `y1` and `y2`.

<details>
  <summary>Example</summary>

```yaml
playername:
  type: textCell
  desc: Player name
  x1: 40
  y1: 125
  x2: 140
  y2: 110
  font: Helvetica
  fontsize: 14
  align: CB
  example: Bob
```
</details>

#### `societyId`

### Misc

#### Text Alignment

##### Horizontal Alignment

* `L`: Left-bound
* `C`: Centered
* `R`: Right-bound

##### Vertical Alignment

* `T`: Top
* `M`: Middle
* `B`: Bottom
* `A`: Baseline

#### Fonts

### Preset Mechanism

## Template Inheritance
