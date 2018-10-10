// Copyright 2015 Monmohan Singh. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xferspdy

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/golang/glog"
)

//Patch is a wrapper on PatchFile (current version only supports patching of local files)
func Patch(delta []Block, sign Fingerprint, t io.Writer) {
	PatchFile(delta, sign.Source, t)
}

// PatchFile takes a source file and Diff as input, and writes out to the Writer.
// The source file would normally be the base version of the file  and
// the Diff is the delta computed by using the Fingerprint generated for the base file and the new version of the file
func PatchFile(delta []Block, source string, t io.Writer) error {
	s, e := ioutil.ReadFile(source)
	if e != nil {
		return e
	}
	buf := bytes.NewReader(s)

	return ProcessBlocks(delta, buf, t)
}

// ProcessBlocks takes a byte reader and applies the blocks to that data
// This is usefull when you want to process blocks in memory as apposed to an actual file
func ProcessBlocks(delta []Block, s *bytes.Reader, t io.Writer) error {
	wptr := int64(0)
	for _, block := range delta {
		if block.HasData {
			glog.V(3).Infof("Writing RawBytes block , wptr=%v , num bytes = %v \n", wptr, len(block.RawBytes))
			_, e := t.Write(block.RawBytes)
			glog.V(4).Infof("Writing bytes = %v \n", block.RawBytes)
			if e != nil {
				return e
			}
			wptr += int64(len(block.RawBytes))
		} else {
			s.Seek(block.Start, 0)
			ds := block.End - block.Start
			glog.V(3).Infof("Writing RawBytes block, Block=%v\n , wptr=%v , num bytes = %v \n", block, wptr, ds)
			if _, e := io.CopyN(t, s, block.End-block.Start); e != nil {
				return e
			}
			wptr += ds
		}
	}
	return nil
}
