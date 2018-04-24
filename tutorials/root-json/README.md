# Usage of JSON converter in ROOT

Macros here showing basic functionality of TBufferJSON class to create JSON data from existing 
objects and read JSON back

ROOT version 6.12 and higher is required to run macros.

## createJSON.C macro

This macros creates user-written class, fills it with data and create JSON for it.
To run it, call:

    [shell] root -l createJSON.C+ -q

As result, `file.json` will be created, which contains object data

    {
      "_typename" : "UserClass",
      "fInt" : 12,
      "fVect" : [0, 0.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5],
      "fStr" : "AnyString"
    }
    
## readJSON.C macro

This macros reads file.json data, using empty object:

    [shell] root -l readJSON.C+ -q

Macro reconstructs object content from JSON and prints result to the std output.

    