// Copyright 2017 The vizmon-demo Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgsvg"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
)

type Plots struct {
	update time.Time
	plots  ControlPlots
	data   MonData
}

func (ps *Plots) MarshalJSON() ([]byte, error) {
	var raw struct {
		Plot   string `json:"plot"`
		Update string `json:"update"`
		Data   string `json:"data"`
	}

	raw.Plot = renderPlot(ps.plots.tile)
	raw.Update = ps.update.Format("2006-01-02 15:04:05 (MST)")

	str := new(bytes.Buffer)
	w := tabwriter.NewWriter(str, 8, 4, 1, ' ', 0)
	for _, d := range ps.data.Values {
		fmt.Fprintf(w, "%s\t%v\n", d.Name, d.Value)
	}
	w.Flush()
	raw.Data = string(str.Bytes())

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(raw)
	if err != nil {
		log.Printf("plots-marshal: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

type ControlPlots struct {
	update time.Time
	tile   *hplot.TiledPlot
}

func renderPlot(p *hplot.TiledPlot) string {
	size := 30 * vg.Centimeter
	canvas := vgsvg.New(size, size/vg.Length(math.Phi))
	p.Draw(draw.New(canvas))
	out := new(bytes.Buffer)
	_, err := canvas.WriteTo(out)
	if err != nil {
		panic(err)
	}
	return string(out.Bytes())
}

func newControlPlots(data []MonData) (ControlPlots, error) {
	var (
		ps  ControlPlots
		err error
	)

	ps.update = time.Now().UTC()
	const pad = 10
	ps.tile, err = hplot.NewTiledPlot(draw.Tiles{
		Cols:      2,
		Rows:      2,
		PadBottom: pad,
		PadLeft:   pad,
		PadRight:  pad,
		PadTop:    pad,
		PadX:      pad,
		PadY:      pad,
	})
	if err != nil {
		return ps, err
	}

	for i, pl := range []*hplot.Plot{
		ps.tile.Plot(0, 0),
		ps.tile.Plot(0, 1),
		ps.tile.Plot(1, 0),
		ps.tile.Plot(1, 1),
	} {
		pl.Title.Text = strings.Title(data[0].Values[i].Name)
		err = setupPlot(pl, data, i)
		if err != nil {
			return ps, err
		}
	}

	return ps, err
}

func (ps *ControlPlots) MarshalJSON() ([]byte, error) {
	var raw struct {
		Plot   string `json:"plot"`
		Update string `json:"update"`
	}

	raw.Plot = renderPlot(ps.tile)
	raw.Update = ps.update.Format("2006-01-02 15:04:05 (MST)")

	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(raw)
	if err != nil {
		log.Printf("plots-marshal: %v", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func setupPlot(pl *hplot.Plot, table []MonData, idx int) error {
	h := hbook.NewH1D(100, 0, 40)
	for _, tbl := range table {
		h.Fill(tbl.Values[idx].Value, 1)
	}
	hh, err := hplot.NewH1D(h)
	if err != nil {
		return err
	}
	hh.LineStyle.Color = color.NRGBA{255, 0, 0, 128}
	hh.FillColor = color.NRGBA{255, 0, 0, 128}

	pl.Add(hh)
	pl.Add(plotter.NewGrid())
	return nil
}
