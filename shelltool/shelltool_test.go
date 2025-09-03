// Copyright 2025 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package shelltool

import (
	"encoding/json"
	"os"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strings"
	"testing"

	"github.com/maruel/genai"
)

func TestGetSandbox(t *testing.T) {
	ipV4 := regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	t.Run("with network access", func(t *testing.T) {
		opts, err := New(true)
		if err != nil {
			t.Fatal(err)
		}
		if opts == nil {
			t.Fatal("excepted opts")
		}

		t.Run("stderr", func(t *testing.T) {
			script, want := "", ""
			if runtime.GOOS == "windows" {
				script = "Write-Output \"hi\"\n[System.Console]::Error.WriteLine(\"hello\")\n"
				want = "hi\r\nhello\r\n"
			} else {
				script = "echo hi\necho hello >&2\n"
				want = "hi\nhello\n"
			}
			b, _ := json.Marshal(&arguments{Script: script})
			msg := genai.Message{Replies: []genai.Reply{{ToolCall: genai.ToolCall{Name: opts.Tools[0].Name, Arguments: string(b)}}}}
			res, err := msg.DoToolCalls(t.Context(), opts.Tools)
			if err != nil {
				t.Log(res.ToolCallResults)
				t.Fatalf("Got error: %v", err)
			}
			if got := res.ToolCallResults[0].Result; got != want {
				t.Fatalf("unexpected output\nwant: %q\ngot:  %q", want, got)
			}
		})

		t.Run("subprocess", func(t *testing.T) {
			script := "/bin/ls\n"
			if runtime.GOOS == "windows" {
				// Intentionally use a subprocess to list the files.
				script = "cmd /c dir /b"
			}

			dirEntries, err := os.ReadDir(".")
			if err != nil {
				t.Fatal(err)
			}
			want := make([]string, 0, len(dirEntries))
			for _, entry := range dirEntries {
				want = append(want, entry.Name())
			}
			sort.Strings(want)
			b, _ := json.Marshal(&arguments{Script: script})
			msg := genai.Message{Replies: []genai.Reply{{ToolCall: genai.ToolCall{Name: opts.Tools[0].Name, Arguments: string(b)}}}}
			res, err := msg.DoToolCalls(t.Context(), opts.Tools)
			if err != nil {
				t.Log(res.ToolCallResults)
				t.Fatalf("Got error: %v", err)
			}
			got := strings.Fields(strings.TrimSpace(res.ToolCallResults[0].Result))
			sort.Strings(got)
			if !slices.Equal(got, want) {
				t.Fatalf("unexpected output\nwant: %q\ngot:  %q", want, got)
			}
		})

		t.Run("network", func(t *testing.T) {
			script := "curl -sS ifconfig.co\n"
			if runtime.GOOS == "windows" {
				script = "(Invoke-WebRequest -Uri https://ifconfig.co -UserAgent curl).Content\n"
			}
			b, _ := json.Marshal(&arguments{Script: script})
			msg := genai.Message{Replies: []genai.Reply{{ToolCall: genai.ToolCall{Name: opts.Tools[0].Name, Arguments: string(b)}}}}
			res, err := msg.DoToolCalls(t.Context(), opts.Tools)
			if err != nil {
				t.Log(res.ToolCallResults)
				t.Fatalf("Got error: %v", err)
			}
			if got := strings.TrimSpace(res.ToolCallResults[0].Result); !ipV4.MatchString(got) {
				t.Fatalf("unexpected output\nwant: IPv4\ngot:  %q", got)
			}
		})
	})

	t.Run("no network access", func(t *testing.T) {
		opts, err := New(false)
		if err != nil {
			if runtime.GOOS == "windows" {
				t.Skip("please send a RP")
			}
			t.Fatal(err)
		} else if runtime.GOOS == "windows" {
			t.Fatal("should have failed")
		}
		if opts == nil {
			t.Fatal("excepted opts")
		}

		t.Run("stderr", func(t *testing.T) {
			script, want := "", ""
			if runtime.GOOS == "windows" {
				script = "Write-Output \"hi\"\n[System.Console]::Error.WriteLine(\"hello\")\n"
				want = "hi\r\nhello\r\n"
			} else {
				script = "echo hi\necho hello >&2\n"
				want = "hi\nhello\n"
			}
			b, _ := json.Marshal(&arguments{Script: script})
			msg := genai.Message{Replies: []genai.Reply{{ToolCall: genai.ToolCall{Name: opts.Tools[0].Name, Arguments: string(b)}}}}
			res, err := msg.DoToolCalls(t.Context(), opts.Tools)
			if err != nil {
				t.Log(res.ToolCallResults)
				t.Fatalf("Got error: %v", err)
			}
			if got := res.ToolCallResults[0].Result; got != want {
				t.Fatalf("unexpected output\nwant: %q\ngot:  %q", want, got)
			}
		})

		t.Run("subprocess", func(t *testing.T) {
			script := "/bin/ls\n"
			if runtime.GOOS == "windows" {
				// Intentionally use a subprocess to list the files.
				script = "cmd /c dir /b"
			}
			dirEntries, err := os.ReadDir(".")
			if err != nil {
				t.Fatal(err)
			}
			want := make([]string, 0, len(dirEntries))
			for _, entry := range dirEntries {
				want = append(want, entry.Name())
			}
			sort.Strings(want)
			b, _ := json.Marshal(&arguments{Script: script})
			msg := genai.Message{Replies: []genai.Reply{{ToolCall: genai.ToolCall{Name: opts.Tools[0].Name, Arguments: string(b)}}}}
			res, err := msg.DoToolCalls(t.Context(), opts.Tools)
			if err != nil {
				t.Log(res.ToolCallResults)
				t.Fatalf("Got error: %v", err)
			}
			got := strings.Fields(strings.TrimSpace(res.ToolCallResults[0].Result))
			sort.Strings(got)
			if !slices.Equal(got, want) {
				t.Fatalf("unexpected output\nwant: %q\ngot:  %q", want, got)
			}
		})

		t.Run("network", func(t *testing.T) {
			script := "curl -sS ifconfig.co\n"
			if runtime.GOOS == "windows" {
				script = "(Invoke-WebRequest -Uri https://ifconfig.co -UserAgent curl).Content\n"
			}
			b, _ := json.Marshal(&arguments{Script: script})
			msg := genai.Message{Replies: []genai.Reply{{ToolCall: genai.ToolCall{Name: opts.Tools[0].Name, Arguments: string(b)}}}}
			res, err := msg.DoToolCalls(t.Context(), opts.Tools)
			if err != nil {
				// That's okay.
				t.Logf("Got error: %v", err)
			} else if got := strings.TrimSpace(res.ToolCallResults[0].Result); ipV4.MatchString(got) {
				t.Fatalf("unexpected output\ndo not want: IPv4\ngot:  %q", got)
			}
		})
	})
}
