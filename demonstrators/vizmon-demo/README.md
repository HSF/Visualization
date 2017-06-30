# vizmon-demo

`vizmon-demo` is a simple web server written in [Go](https://golang.org).
It demonstrates serving histograms over WebSockets.
The plots are converted to SVG on the server side and transfered thru the WebSockets (PNGs could be used as well, but SVGs are lightweighter.)
In a real slow-control application ([here](https://github.com/go-lsst/fcs-lpc-motor-ctl)), the data is displayed as histograms and video frame stills.


As Go compiles to completely statically compiled executables, deployment is really easy.
Also, cross-compiling is simple:

```
$> GOARCH=arm                go build -o my-app-linux-arm.exe
$> GOARCH=386   GOOS=linux   go build -o my-app-linux-386.exe
$> GOARCH=amd64 GOOS=windows go build -o my-app-win64.exe
```

The histograms are created via [go-hep/hbook](https://go-hep.org/x/hep/hbook) and displayed with [go-hep/hplot](https://go-hep.org/x/hep/hplot).

## Example

Once the Go toolkit has been installed, `vizmon-demo` can be installed like so:

```sh
$> go get github.com/HEP-SF/Visualization/demonstrators/vizmon-demo
$> vizmon-demo
vizmon-demo: starting up web-server on: :8080
vizmon-demo: daq: {Values:[{Name:vtx-x Value:-0.9292493828416044} {Name:vtx-y Value:9.935965908773898} {Name:vtx-z Value:0.4648125393918842} {Name:ele-tof Value:17.23179836038509}]}
```

Alternatively, binaries are available here (so no Go toolchain needed):

- https://cern.ch/binet/hsf/vizmon-demo-linux-386.exe
- https://cern.ch/binet/hsf/vizmon-demo-linux-amd64.exe
- https://cern.ch/binet/hsf/vizmon-demo-linux-arm.exe (_e.g._ for RaspBerry Pi)
- https://cern.ch/binet/hsf/vizmon-demo-macos-386.exe
- https://cern.ch/binet/hsf/vizmon-demo-macos-amd64.exe
- https://cern.ch/binet/hsf/vizmon-demo-windows-386.exe
- https://cern.ch/binet/hsf/vizmon-demo-windows-amd64.exe

Then, one can open their favorite web-browser and go to `localhost:8080`.

Additionally, one can also retrieve informations from the `"/echo"` end-point:

```sh
$> curl localhost:8080/echo
{"values":[{"name":"vtx-x","value":0.007457764495316657},{"name":"vtx-y","value":9.285480166569435},{"name":"vtx-z","value":-0.0024393582606496394},{"name":"ele-tof","value":22.441684011690615}]}
```

![image](https://github.com/HEP-SF/Visualization/raw/master/demonstrators/vizmon-demo/screenshot.png)

## Example reading ROOT files

Another example is [root-srv](https://godoc.org/go-hep.org/x/hep/rootio/cmd/root-srv): another web-browser based tool to inspect [ROOT](https://root.cern.ch) files without any ROOT installation needed.

An instance is available on Google AppEngine:

- http://rootio-inspector.appspot.com

(not every single ROOT file and TTree is supported: `go-hep/rootio` is still a work in progress.)
