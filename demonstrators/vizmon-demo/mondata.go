// Copyright 2017 The vizmon-demo Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

type MonData struct {
	Values []MonValue `json:"values"`
}

type MonValue struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}
