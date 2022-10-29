package parser

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/chavacava/changelogger/model"
)

type Default struct{}

type tokenKind int

const (
	kindEmpty tokenKind = iota
	kindTitle
	kindVersion
	kindSubsection
	kindEntry
	kindPlain
	kindEOF
	kindError
)

type state int

const (
	initial state = iota
	title
	version
	subsection
	entry
)

type token struct {
	fullText string
	kind     tokenKind
	pos      int
}

func (p Default) Parse(r io.Reader, config any) (*model.Changelog, error) {
	result := model.NewChangelog()
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	tokens := make(chan token)
	go func() {
		pos := 0
		for scanner.Scan() {
			pos++
			line := scanner.Text()
			lineKind := p.retrieveLineKind(line)
			if lineKind == kindEmpty {
				continue
			}
			tokens <- token{fullText: line, kind: lineKind, pos: pos}
		}
		tokens <- token{fullText: "", kind: kindEOF, pos: pos + 1}
		close(tokens)
	}()

	state := initial
	tok := <-tokens
	var currentVersion *model.Version
	var currentSubsection *model.Subsection
	for {
		switch state {
		case initial:
			switch tok.kind {
			case kindTitle:
				state = title
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line: %s\n expecting empty line or main title", tok.fullText)
			}
		case title:
			tok = <-tokens
			switch tok.kind {
			case kindPlain:
				result.Header = append(result.Header, tok.fullText)
				tok = <-tokens
			case kindVersion:
				state = version
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line: %s\nexpecting standard text line or version", tok.fullText)
			}
		case version:
			newVersion := &model.Version{Version: tok.fullText}
			result.Versions = append(result.Versions, newVersion)
			currentVersion = newVersion
			tok = <-tokens
			switch tok.kind {
			case kindVersion:
				state = version
			case kindSubsection:
				state = subsection
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line:%s\nexpecting subsection or version", tok.fullText)
			}
		case subsection:
			newSubsection := &model.Subsection{Name: tok.fullText}
			currentVersion.Subsections = append(currentVersion.Subsections, newSubsection)
			currentSubsection = newSubsection
			tok = <-tokens
			switch tok.kind {
			case kindSubsection:
				state = subsection
			case kindVersion:
				state = version
			case kindEntry:
				state = entry
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line:%s\nexpecting subsection, version or change description", tok.fullText)
			}
		case entry:
			newEntry := model.ChangeLine{Summary: tok.fullText}
			currentSubsection.History = append(currentSubsection.History, &newEntry)
			tok = <-tokens
			switch tok.kind {
			case kindSubsection:
				state = subsection
			case kindVersion:
				state = version
			case kindEntry:
				state = entry
			case kindEOF:
				return result, nil
			default:
				return nil, fmt.Errorf("unexpected line:%s\nexpecting subsection, version or change description", tok.fullText)
			}
		}
	}
}

func (p Default) retrieveLineKind(line string) tokenKind {
	trimedLine := strings.Trim(line, " ")
	if strings.HasPrefix(trimedLine, "###") {
		return kindSubsection
	}

	if strings.HasPrefix(trimedLine, "##") {
		return kindVersion
	}

	if strings.HasPrefix(trimedLine, "#") {
		return kindTitle
	}

	if strings.HasPrefix(trimedLine, "-") || strings.HasPrefix(trimedLine, "*") {
		return kindEntry
	}

	if trimedLine == "" {
		return kindEmpty
	}

	return kindPlain
}
