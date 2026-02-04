// Copyright 2025 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package genaitools_test

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/maruel/genai"
	"github.com/maruel/genai/adapters"
	"github.com/maruel/genai/providers/anthropic"
	"github.com/maruel/genaitools"
	"github.com/maruel/roundtrippers"
)

func Example_genSyncWithToolCallLoop_with_custom_HTTP_Header() {
	// Modified version of the example in package adapters, with a custom header.
	//
	// As of June 2025, interleaved thinking can be enabled with a custom header.
	// https://docs.anthropic.com/en/docs/build-with-claude/extended-thinking#interleaved-thinking
	wrapper := genai.ProviderOptionTransportWrapper(func(h http.RoundTripper) http.RoundTripper {
		return &roundtrippers.Header{
			Transport: h,
			Header:    http.Header{"anthropic-beta": []string{"interleaved-thinking-2025-05-14"}},
		}
	})
	ctx := context.Background()
	c, err := anthropic.New(ctx, genai.ProviderOptionModel("claude-sonnet-4-20250514"), wrapper)
	if err != nil {
		log.Fatal(err)
	}
	msgs := genai.Messages{
		genai.NewTextMessage("What is 3214 + 5632? Leverage the tool available to you to tell me the answer. Do not explain. Be terse. Include only the answer."),
	}
	opts := genai.GenOptionsTools{
		Tools: []genai.ToolDef{genaitools.Arithmetic},
		// Force the LLM to do a tool call first.
		Force: genai.ToolCallRequired,
	}
	newMsgs, _, err := adapters.GenSyncWithToolCallLoop(ctx, c, msgs, &opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", newMsgs[len(newMsgs)-1].String())
	// Remove this comment line to run the example.
	// Output:
	// 8846
}
