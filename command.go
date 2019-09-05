// Copyright 2019 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tdaq is a minimal toolkit to implement a tiny data acquisition system.
package tdaq // import "github.com/go-daq/tdaq"

//go:generate stringer -type CmdType -output z_cmdtype_string.go .

import (
	"bytes"
	"context"
	"io"

	"golang.org/x/xerrors"
)

type CmdType byte

const (
	CmdUnknown CmdType = iota
	CmdJoin
	CmdConnect
	CmdConfig
	CmdInit
	CmdReset
	CmdStart
	CmdStop
	CmdTerm
	CmdStatus
	CmdLog
)

var cmdNames = [...][]byte{
	CmdUnknown: []byte("/unknown"),
	CmdJoin:    []byte("/join"),
	CmdConnect: []byte("/connect"),
	CmdConfig:  []byte("/config"),
	CmdInit:    []byte("/init"),
	CmdReset:   []byte("/reset"),
	CmdStart:   []byte("/start"),
	CmdStop:    []byte("/stop"),
	CmdTerm:    []byte("/term"),
	CmdStatus:  []byte("/status"),
	CmdLog:     []byte("/log"),
}

func cmdTypeToPath(cmd CmdType) []byte {
	return cmdNames[cmd]
}

type Cmder interface {
	Marshaler
	Unmarshaler

	CmdType() CmdType
}

type Cmd struct {
	Type CmdType
	Body []byte
}

func CmdFrom(frame Frame) (Cmd, error) {
	if frame.Type != FrameCmd {
		return Cmd{}, xerrors.Errorf("invalid frame type %v", frame.Type)
	}
	cmd := Cmd{
		Type: CmdType(frame.Body[0]),
		Body: frame.Body[1:],
	}
	return cmd, nil
}

func (raw Cmd) cmd() (cmd Cmder, err error) {
	switch raw.Type {
	case CmdJoin:
		var c JoinCmd
		err = c.UnmarshalTDAQ(raw.Body[1:])
		cmd = &c
	case CmdConnect:
		panic("not implemented")
	case CmdInit:
		panic("not implemented")
	case CmdConfig:
		panic("not implemented")
	case CmdReset:
		panic("not implemented")
	case CmdStart:
		panic("not implemented")
	case CmdStop:
		panic("not implemented")
	case CmdTerm:
		panic("not implemented")
	case CmdStatus:
		panic("not implemented")
	case CmdLog:
		panic("not implemented")
	default:
		return nil, xerrors.Errorf("invalid cmd type %q", raw.Type)
	}
	return cmd, err
}

func SendCmd(ctx context.Context, w io.Writer, cmd Cmder) error {
	raw, err := cmd.MarshalTDAQ()
	if err != nil {
		return xerrors.Errorf("could not marshal cmd: %w", err)
	}

	ctype := cmd.CmdType()
	path := cmdTypeToPath(cmd.CmdType())
	return sendFrame(ctx, w, FrameCmd, path, append([]byte{byte(ctype)}, raw...))
}

func sendCmd(ctx context.Context, w io.Writer, ctype CmdType, body []byte) error {
	path := cmdTypeToPath(ctype)
	return sendFrame(ctx, w, FrameCmd, path, append([]byte{byte(ctype)}, body...))
}

func recvCmd(ctx context.Context, r io.Reader) (cmd Cmd, err error) {
	frame, err := RecvFrame(ctx, r)
	if err != nil {
		return cmd, xerrors.Errorf("could not receive TDAQ cmd: %w", err)
	}
	if frame.Type != FrameCmd {
		return cmd, xerrors.Errorf("did not receive a TDAQ cmd")
	}
	return Cmd{Type: CmdType(frame.Body[0]), Body: frame.Body[1:]}, nil
}

type JoinCmd struct {
	Name     string
	InPorts  []Port
	OutPorts []Port
}

func (cmd JoinCmd) CmdType() CmdType { return CmdJoin }

func (cmd JoinCmd) MarshalTDAQ() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.WriteStr(cmd.Name)

	enc.WriteU64(uint64(len(cmd.InPorts)))
	for _, p := range cmd.InPorts {
		enc.WriteStr(p.Name)
		enc.WriteStr(p.Addr)
		enc.WriteStr(p.Type)
	}

	enc.WriteU64(uint64(len(cmd.OutPorts)))
	for _, p := range cmd.OutPorts {
		enc.WriteStr(p.Name)
		enc.WriteStr(p.Addr)
		enc.WriteStr(p.Type)
	}
	return buf.Bytes(), enc.err
}

func (cmd *JoinCmd) UnmarshalTDAQ(p []byte) error {
	dec := NewDecoder(bytes.NewReader(p))

	cmd.Name = dec.ReadStr()
	n := int(dec.ReadU64())
	cmd.InPorts = make([]Port, n)
	for i := range cmd.InPorts {
		p := &cmd.InPorts[i]
		p.Name = dec.ReadStr()
		p.Addr = dec.ReadStr()
		p.Type = dec.ReadStr()
	}

	n = int(dec.ReadU64())
	cmd.OutPorts = make([]Port, n)
	for i := range cmd.OutPorts {
		p := &cmd.OutPorts[i]
		p.Name = dec.ReadStr()
		p.Addr = dec.ReadStr()
		p.Type = dec.ReadStr()
	}

	return dec.err
}

type ConfigCmd struct {
	Name     string
	InPorts  []Port
	OutPorts []Port
}

func newConfigCmd(frame Frame) (ConfigCmd, error) {
	var (
		cfg ConfigCmd
		err error
	)

	raw, err := CmdFrom(frame)
	if err != nil {
		return cfg, xerrors.Errorf("not a /config cmd: %w", err)
	}

	if raw.Type != CmdConfig {
		return cfg, xerrors.Errorf("not a /config cmd")
	}

	err = cfg.UnmarshalTDAQ(raw.Body)
	return cfg, err
}

func (cmd ConfigCmd) CmdType() CmdType { return CmdConfig }

func (cmd ConfigCmd) MarshalTDAQ() ([]byte, error) {
	buf := new(bytes.Buffer)
	enc := NewEncoder(buf)
	enc.WriteStr(cmd.Name)

	enc.WriteU64(uint64(len(cmd.InPorts)))
	for _, p := range cmd.InPorts {
		enc.WriteStr(p.Name)
		enc.WriteStr(p.Addr)
		enc.WriteStr(p.Type)
	}

	enc.WriteU64(uint64(len(cmd.OutPorts)))
	for _, p := range cmd.OutPorts {
		enc.WriteStr(p.Name)
		enc.WriteStr(p.Addr)
		enc.WriteStr(p.Type)
	}
	return buf.Bytes(), enc.err
}

func (cmd *ConfigCmd) UnmarshalTDAQ(p []byte) error {
	dec := NewDecoder(bytes.NewReader(p))

	cmd.Name = dec.ReadStr()
	n := int(dec.ReadU64())
	cmd.InPorts = make([]Port, n)
	for i := range cmd.InPorts {
		p := &cmd.InPorts[i]
		p.Name = dec.ReadStr()
		p.Addr = dec.ReadStr()
		p.Type = dec.ReadStr()
	}

	n = int(dec.ReadU64())
	cmd.OutPorts = make([]Port, n)
	for i := range cmd.OutPorts {
		p := &cmd.OutPorts[i]
		p.Name = dec.ReadStr()
		p.Addr = dec.ReadStr()
		p.Type = dec.ReadStr()
	}

	return dec.err
}

var (
	_ Cmder       = (*JoinCmd)(nil)
	_ Marshaler   = (*JoinCmd)(nil)
	_ Unmarshaler = (*JoinCmd)(nil)

	_ Cmder       = (*ConfigCmd)(nil)
	_ Marshaler   = (*ConfigCmd)(nil)
	_ Unmarshaler = (*ConfigCmd)(nil)
)