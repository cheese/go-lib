/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package kv reads key value files.
package kv

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"unicode"
)

type Reader struct {
	// Delim is key value delimiter
	Delim byte
	// Comment, if not 0, is the comment character. Lines begin with the
	// Comment character are ignored.
	Comment byte
	// TrimSpace determines the behavior of trim space
	TrimSpace TrimSpaceFlag
	r         *bufio.Reader
}

type TrimSpaceFlag uint

const (
	TrimLeadingSpace TrimSpaceFlag = iota + 1
	TrimTailingSpace
	TrimDelimLeftSpace
	TrimDelimRightSpace
)
const TrimAllSpace TrimSpaceFlag = TrimLeadingSpace | TrimTailingSpace | TrimDelimLeftSpace | TrimDelimRightSpace
const TrimLeadingTailingSpace TrimSpaceFlag = TrimLeadingSpace | TrimTailingSpace

// NewReader returns a new Reader that reads from r.
// The Delim field default to '=', the TrimSpace field default to TrimAllSpace.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Delim:     '=',
		r:         bufio.NewReader(r),
		TrimSpace: TrimAllSpace,
	}
}

type Pair struct {
	Key   string
	Value string
}

// Read reads one pair from r.
func (r *Reader) Read() (pair *Pair, err error) {
	for {
		pair, err = r.parseLine()
		if pair != nil {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	return pair, nil
}

// ReadAll reads all the remaining pairs from r.
// A successful call returns err == nil, not err == io.EOF.
// Because ReadAll is defined to read until EOF,
// it does not treat end of file as an error to be reported.
func (r *Reader) ReadAll() (pairs []*Pair, err error) {
	for {
		pair, err := r.Read()
		if err == io.EOF {
			return pairs, nil
		}
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}
}

func (r *Reader) skipLine() error {
	_, err := r.r.ReadBytes('\n')
	return err
}

var ErrBadLine = errors.New("bad line")

func (r *Reader) parseLine() (*Pair, error) {
	line, err := r.r.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	if r.TrimSpace&TrimLeadingSpace != 0 {
		line = bytes.TrimLeftFunc(line, unicode.IsSpace)
	}
	if r.TrimSpace&TrimTailingSpace != 0 {
		line = bytes.TrimRightFunc(line, unicode.IsSpace)
	}

	// skip empty line
	if len(line) == 0 {
		return nil, nil
	}

	b1 := line[0]
	if r.Comment != 0 && b1 == r.Comment {
		// skip comment line
		return nil, nil
	}

	parts := bytes.SplitN(line, []byte{r.Delim}, 2)
	if len(parts) != 2 {
		return nil, ErrBadLine
	}

	key := parts[0]
	if r.TrimSpace&TrimDelimLeftSpace != 0 {
		key = bytes.TrimRightFunc(key, unicode.IsSpace)
	}

	value := parts[1]
	if r.TrimSpace&TrimDelimRightSpace != 0 {
		value = bytes.TrimLeftFunc(value, unicode.IsSpace)
	}
	return &Pair{
		Key:   string(key),
		Value: string(value),
	}, nil
}
