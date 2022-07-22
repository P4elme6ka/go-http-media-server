package tpl

import (
	"bytes"
	"github.com/P4elme6ka/go-http-media-server/src/util"
	"io"
)

type asset struct {
	contentType string
	readSeeker  io.ReadSeeker
}

type assets map[string]asset

func (assets assets) set(path string, content []byte) error {
	rd := bytes.NewReader(content)
	ctype, err := util.GetContentType(path, rd)
	if err != nil {
		return err
	}

	asset := asset{
		contentType: ctype,
		readSeeker:  rd,
	}
	assets[path] = asset
	return nil
}
