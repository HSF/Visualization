# Exporting geometry to JSON

This macros takes some existing ATLAS geometry, show approx 5000 most significant volumes and
export it to the JSON format. To run macro, call:

    [shell] root -l geomAtlas.C

It loads file from web, import geometry, visualize it and stores visible elements into JSON.
There are three files created:

   - atlas2.json [76M] without any compression
   - atlas2.json.gz [1.2M] without new lines, spaces and gzip
   - atlas2.root [1.6M] is same information in binary format
   
Macro runs approximately 30 s because of large JSON structures.

Produced geometry can be directly displayed with JSROOT

[https://root.cern/js/latest/?nobrowser&file=../files/geom/atlas2.root&item=atlas&opt=black](https://root.cern/js/latest/?nobrowser&file=../files/geom/atlas2.root&item=atlas&opt=black)
