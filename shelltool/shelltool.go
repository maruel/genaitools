// Copyright 2025 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package shelltool makes a sandboxed shell available as a tool to the LLM.
package shelltool

import (
	"fmt"
	"os"

	"github.com/maruel/genai"
)

// New return a shell tool that works on the current OS.
//
// If allowNetwork is false, the script will not have network access.
//
//   - On macOS, it runs /bin/zsh under sandbox-exec.
//   - On Windows, it runs powershell under a restricted user token. It is currently disabled due to a crash in the Go runtime.
//   - On other platforms, it runs bash under bubblewrap. bubblewrap must be installed separately.
func New(allowNetwork bool) (*genai.OptionsTools, error) {
	return getShellTool(allowNetwork)
}

// arguments is the shell tool argument.
type arguments struct {
	Script string `json:"script"`
}

func writeTempFile(g, content string) (string, error) {
	f, err := os.CreateTemp("", g)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	n := f.Name()
	if _, err = f.WriteString(content); err != nil {
		_ = os.Remove(n)
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}
	err = f.Close()
	return n, err
}
