// Copyright 2019 The go-daq Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package tdaq_test // import "github.com/go-daq/tdaq"

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-daq/tdaq"
	"github.com/go-daq/tdaq/fsm"
	"github.com/go-daq/tdaq/internal/iomux"
)

func TestCommands(t *testing.T) {
	for _, tt := range []struct {
		name string
		want interface{}
	}{
		{
			name: "join",
			want: &tdaq.JoinCmd{
				Name: "n1",
				InEndPoints: []tdaq.EndPoint{
					{"n11", "addr11", "type11"},
					{"n12", "addr12", "type12"},
				},
				OutEndPoints: []tdaq.EndPoint{
					{"n11", "addr11", "type11"},
					{"n12", "addr12", "type12"},
					{"n13", "addr13", "type13"},
				},
			},
		},
		{
			name: "config",
			want: &tdaq.ConfigCmd{
				Name: "n1",
				InEndPoints: []tdaq.EndPoint{
					{"n11", "addr11", "type11"},
					{"n12", "addr12", "type12"},
				},
				OutEndPoints: []tdaq.EndPoint{
					{"n11", "addr11", "type11"},
					{"n12", "addr12", "type12"},
					{"n13", "addr13", "type13"},
				},
			},
		},
		{
			name: "status-unconf",
			want: &tdaq.StatusCmd{Name: "n1", Status: fsm.UnConf},
		},
		{
			name: "status-conf",
			want: &tdaq.StatusCmd{Name: "n1", Status: fsm.Conf},
		},
		{
			name: "status-init",
			want: &tdaq.StatusCmd{Name: "n1", Status: fsm.Init},
		},
		{
			name: "status-stopped",
			want: &tdaq.StatusCmd{Name: "n1", Status: fsm.Stopped},
		},
		{
			name: "status-running",
			want: &tdaq.StatusCmd{Name: "n1", Status: fsm.Running},
		},
		{
			name: "status-exiting",
			want: &tdaq.StatusCmd{Name: "n1", Status: fsm.Exiting},
		},
		{
			name: "status-error",
			want: &tdaq.StatusCmd{Name: "n1", Status: fsm.Error},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			{
				buf := new(iomux.Socket)
				enc := tdaq.NewEncoder(buf)
				err := enc.Encode(tt.want)
				if err != nil {
					t.Fatalf("could not encode %v: %+v", tt.want, err)
				}

				dec := tdaq.NewDecoder(buf)
				got := reflect.New(reflect.TypeOf(tt.want).Elem()).Elem()
				err = dec.Decode(got.Addr().Interface())
				if err != nil {
					t.Fatalf("could not decode %v: %+v", tt.want, err)
				}

				if got, want := got.Addr().Interface(), tt.want; !reflect.DeepEqual(got, want) {
					t.Fatalf("invalid round-trip:\ngot = %#v\nwant= %#v\n", got, want)
				}
			}

			{
				buf := new(iomux.Socket)
				ctx := context.Background()

				err := tdaq.SendCmd(ctx, buf, tt.want.(tdaq.Cmder))
				if err != nil {
					t.Fatalf("could not send cmd-frame: %+v", err)
				}

				frame, err := tdaq.RecvFrame(ctx, buf)
				if err != nil {
					t.Fatalf("could not recv cmd-frame: %+v", err)
				}

				got := reflect.New(reflect.TypeOf(tt.want).Elem()).Elem()
				err = got.Addr().Interface().(tdaq.Cmder).UnmarshalTDAQ(frame.Body[1:])
				if err != nil {
					t.Fatalf("could not unmarshal cmd-frame: %+v", err)
				}

				if got, want := got.Addr().Interface(), tt.want; !reflect.DeepEqual(got, want) {
					t.Fatalf("invalid cmd-frame round-trip:\ngot = %#v\nwant= %#v\n", got, want)
				}
			}
		})
	}
}

func TestCmdType(t *testing.T) {
	for _, tt := range []struct {
		cmd    tdaq.CmdType
		want   string
		panics bool
	}{
		{cmd: tdaq.CmdUnknown, want: "/unknown"},
		{cmd: tdaq.CmdJoin, want: "/join"},
		{cmd: tdaq.CmdConfig, want: "/config"},
		{cmd: tdaq.CmdInit, want: "/init"},
		{cmd: tdaq.CmdReset, want: "/reset"},
		{cmd: tdaq.CmdStart, want: "/start"},
		{cmd: tdaq.CmdStop, want: "/stop"},
		{cmd: tdaq.CmdQuit, want: "/quit"},
		{cmd: tdaq.CmdStatus, want: "/status"},
		{cmd: tdaq.CmdType(255), panics: true},
	} {
		t.Run("", func(t *testing.T) {
			if tt.panics {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic")
					}
					if got, want := err.(error).Error(), "invalid cmd-type 255"; got != want {
						t.Fatalf("invalid panic string.\ngot = %q\nwant= %q\n", got, want)
					}
				}()
			}

			got := tt.cmd.String()
			if got != tt.want {
				t.Fatalf("invalid stringer value.\ngot = %q\nwant= %q\n", got, tt.want)
			}
		})
	}
}
