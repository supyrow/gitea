// Copyright 2022 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

//go:generate go run invisible/generate.go -v -o ./invisible_gen.go

//go:generate go run ambiguous/generate.go -v -o ./ambiguous_gen.go ambiguous/ambiguous.json

package charset

import (
	"io"
	"strings"

	"code.gitea.io/gitea/modules/log"
	"code.gitea.io/gitea/modules/translation"
)

// RuneNBSP is the codepoint for NBSP
const RuneNBSP = 0xa0

// EscapeControlHTML escapes the unicode control sequences in a provided html document
func EscapeControlHTML(text string, locale translation.Locale, allowed ...rune) (escaped *EscapeStatus, output string) {
	sb := &strings.Builder{}
	outputStream := &HTMLStreamerWriter{Writer: sb}
	streamer := NewEscapeStreamer(locale, outputStream, allowed...).(*escapeStreamer)

	if err := StreamHTML(strings.NewReader(text), streamer); err != nil {
		streamer.escaped.HasError = true
		log.Error("Error whilst escaping: %v", err)
	}
	return streamer.escaped, sb.String()
}

// EscapeControlReaders escapes the unicode control sequences in a provider reader and writer in a locale and returns the findings as an EscapeStatus and the escaped []byte
func EscapeControlReader(reader io.Reader, writer io.Writer, locale translation.Locale, allowed ...rune) (escaped *EscapeStatus, err error) {
	outputStream := &HTMLStreamerWriter{Writer: writer}
	streamer := NewEscapeStreamer(locale, outputStream, allowed...).(*escapeStreamer)

	if err = StreamHTML(reader, streamer); err != nil {
		streamer.escaped.HasError = true
		log.Error("Error whilst escaping: %v", err)
	}
	return streamer.escaped, err
}

// EscapeControlString escapes the unicode control sequences in a provided string and returns the findings as an EscapeStatus and the escaped string
func EscapeControlString(text string, locale translation.Locale, allowed ...rune) (escaped *EscapeStatus, output string) {
	sb := &strings.Builder{}
	outputStream := &HTMLStreamerWriter{Writer: sb}
	streamer := NewEscapeStreamer(locale, outputStream, allowed...).(*escapeStreamer)

	if err := streamer.Text(text); err != nil {
		streamer.escaped.HasError = true
		log.Error("Error whilst escaping: %v", err)
	}
	return streamer.escaped, sb.String()
}
