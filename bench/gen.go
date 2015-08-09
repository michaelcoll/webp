// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package main

// This program generates x_benchmark_test.go. Invoke it as
//	go run gen.go -output x_benchmark_test.go

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

var (
	flagOutputFilename = flag.String("output", "x_benchmark_test.go", "output file name")
)

var testFilenames = []string{
	"1_webp_a.webp",
	"1_webp_ll.webp",
	"2_webp_a.webp",
	"2_webp_ll.webp",
	"3_webp_a.webp",
	"3_webp_ll.webp",
	"4_webp_a.webp",
	"4_webp_ll.webp",
	"5_webp_a.webp",
	"5_webp_ll.webp",
	"blue-purple-pink-large.lossless.webp",
	"blue-purple-pink-large.no-filter.lossy.webp",
	"blue-purple-pink-large.normal-filter.lossy.webp",
	"blue-purple-pink-large.simple-filter.lossy.webp",
	"blue-purple-pink.lossless.webp",
	"blue-purple-pink.lossy.webp",
	"gopher-doc.1bpp.lossless.webp",
	"gopher-doc.2bpp.lossless.webp",
	"gopher-doc.4bpp.lossless.webp",
	"gopher-doc.8bpp.lossless.webp",
	"tux.lossless.webp",
	"video-001.lossy.webp",
	"video-001.webp",
	"yellow_rose.lossless.webp",
	"yellow_rose.lossy-with-alpha.webp",
	"yellow_rose.lossy.webp",
}

func main() {
	flag.Parse()

	var buf bytes.Buffer
	printHeader(&buf, *flagOutputFilename)
	for _, filename := range testFilenames {
		printTestCase(&buf, filename)
	}

	data, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(*flagOutputFilename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func printHeader(w io.Writer, outputFilename string) {
	fmt.Fprintf(w, `
// Copyright 2015 <chaishushan{AT}gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// generated by go run gen.go -output %s; DO NOT EDIT

package webp_bench

import (
	"bytes"
	"io/ioutil"
	"testing"

	chai2010_webp "github.com/chai2010/webp"
	x_image_webp "golang.org/x/image/webp"
)

type CBuffer interface {
	chai2010_webp.CBuffer
}

func tbLoadData(tb testing.TB, filename string) []byte {
	data, err := ioutil.ReadFile("../testdata/" + filename)
	if err != nil {
		tb.Fatal(err)
	}
	return data
}

func tbLoadCData(tb testing.TB, filename string) CBuffer {
	data, err := ioutil.ReadFile("../testdata/" + filename)
	if err != nil {
		tb.Fatal(err)
	}
	cbuf := chai2010_webp.NewCBuffer(len(data))
	copy(cbuf.CData(), data)
	return cbuf
}

`[1:], outputFilename)
}

func printTestCase(w io.Writer, filename string) {
	s := `
func BenchmarkDecode_{{.goodBaseName}}_chai2010_webp(b *testing.B) {
	data := tbLoadData(b, "{{.filename}}")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err := chai2010_webp.Decode(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
		_ = m
	}
}

func BenchmarkDecode_{{.goodBaseName}}_chai2010_webp_cbuf(b *testing.B) {
	cbuf := tbLoadCData(b, "{{.filename}}")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, pix, err := chai2010_webp.DecodeRGBAEx(cbuf.CData(), cbuf)
		if err != nil {
			b.Fatal(err)
		}
		_ = m
		pix.Close()
	}
}

func BenchmarkDecode_{{.goodBaseName}}_x_image_webp(b *testing.B) {
	data := tbLoadData(b, "{{.filename}}")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m, err := x_image_webp.Decode(bytes.NewReader(data))
		if err != nil {
			b.Fatal(err)
		}
		_ = m
	}
}

`
	s = strings.Replace(s, "{{.goodBaseName}}", goodBaseName(filename), -1)
	s = strings.Replace(s, "{{.filename}}", filename, -1)
	fmt.Fprintln(w, s)
}

func goodBaseName(name string) string {
	name = filepath.Base(name)
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[:idx]
	}
	temp := []rune(name)
	for i := 0; i < len(temp); i++ {
		switch temp[i] {
		case '.', '-':
			temp[i] = '_'
		}
	}
	return string(temp)
}