package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/chavacava/changelog-lint/model"
)

// Default CHANGELOG parser
// Parses Markdown changelogs
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
	entryContinuation
)

type token struct {
	fullText string
	kind     tokenKind
	pos      int
}

type decoratorConfig struct {
	titlePattern      *regexp.Regexp
	versionPattern    *regexp.Regexp
	subsectionPattern *regexp.Regexp
	entryPattern      *regexp.Regexp
}

// Parse parses a changelog
func (p Default) Parse(r io.Reader, config *Config) (*model.Changelog, error) {
	cl, err := p.parse(r)
	if err != nil {
		return nil, err
	}

	err = p.decorateChangelog(cl, p.extractDecoratorConfig(config))

	return cl, err
}

func (p Default) decorateChangelog(cl *model.Changelog, config decoratorConfig) error {
	if !config.titlePattern.Match([]byte(cl.Header[0])) {
		return p.normalizeError(
			fmt.Sprintf("the title\n\t%s\ndoes not match %s", cl.Header[0], config.titlePattern.String()),
			1,
		)
	}
	for _, v := range cl.Versions {
		err := p.decorateVersion(v, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p Default) extractDecoratorConfig(conf *Config) decoratorConfig {
	return decoratorConfig{
		titlePattern:      conf.TitlePattern,
		versionPattern:    conf.VersionPattern,
		subsectionPattern: conf.SubsectionPattern,
		entryPattern:      conf.EntryPattern,
	}
}

func (p Default) decorateVersion(v *model.Version, config decoratorConfig) error {
	matches := config.versionPattern.FindStringSubmatch(v.SourceLine)
	if len(matches) < 2 {
		return p.normalizeError(
			fmt.Sprintf("the version\n\t%s\ndoes not match %s", v.SourceLine, config.versionPattern.String()),
			v.Position,
		)
	}

	v.Version = matches[1]

	for _, s := range v.Subsections {
		err := p.decorateSubsection(s, config)
		if err != nil {
			return fmt.Errorf("version %q contains a bad formated subsection: %v", v.Version, err)
		}
	}

	return nil
}

func (p Default) decorateSubsection(s *model.Subsection, config decoratorConfig) error {
	matches := config.subsectionPattern.FindStringSubmatch(s.SourceLine)
	if len(matches) < 2 {
		return p.normalizeError(
			fmt.Sprintf("the subsection\n\t%s\ndoes not match %s", s.SourceLine, config.subsectionPattern.String()),
			s.Position,
		)
	}
	s.Name = matches[1]

	for _, e := range s.History {
		err := p.decorateEntry(e, config)
		if err != nil {
			return fmt.Errorf("subsection %q contains a bad formated entry: %v", s.Name, err)
		}
	}

	return nil
}

func (p Default) decorateEntry(e *model.Entry, config decoratorConfig) error {
	matches := config.entryPattern.FindStringSubmatch(e.Summary)
	if len(matches) < 1 {
		return p.normalizeError(
			fmt.Sprintf("the entry\n\t%s\ndoes not match %s", e.Summary, config.entryPattern.String()),
			e.Position)
	}

	return nil
}

func (p Default) parse(r io.Reader) (*model.Changelog, error) {
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
	var currentEntry *model.Entry
	for {
		switch state {
		case initial:
			switch tok.kind {
			case kindTitle:
				state = title
			case kindEOF:
				return nil, p.normalizeError("unexpected end of file", tok.pos)
			default:
				return nil, p.normalizeError(
					fmt.Sprintf("unexpected line: %s\nexpecting empty line or main title", tok.fullText),
					tok.pos,
				)
			}
		case title:
			result.Header = append(result.Header, tok.fullText)
			tok = <-tokens
			switch tok.kind {
			case kindPlain:
				// do nothing -> it will return to case title
			case kindVersion:
				state = version
			case kindEOF:
				return nil, p.normalizeError("unexpected end of file", tok.pos)
			default:
				return nil, p.normalizeError(
					fmt.Sprintf("unexpected line: %s\nexpecting plain text line or version", tok.fullText),
					tok.pos,
				)
			}
		case version:
			newVersion := &model.Version{SourceLine: tok.fullText, Position: tok.pos}
			result.Versions = append(result.Versions, newVersion)
			currentVersion = newVersion
			tok = <-tokens
			switch tok.kind {
			case kindVersion:
				state = version
			case kindSubsection:
				state = subsection
			case kindEOF:
				return nil, p.normalizeError("unexpected end of file", tok.pos)
			default:
				return nil, p.normalizeError(
					fmt.Sprintf("unexpected line:%s\nexpecting subsection or version", tok.fullText),
					tok.pos,
				)
			}
		case subsection:
			newSubsection := &model.Subsection{SourceLine: tok.fullText, Position: tok.pos}
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
				return nil, p.normalizeError("unexpected end of file", tok.pos)
			default:
				return nil, p.normalizeError(
					fmt.Sprintf("unexpected line:%s\nexpecting subsection, version or change description", tok.fullText),
					tok.pos,
				)
			}
		case entry:
			newEntry := &model.Entry{Summary: tok.fullText, Position: tok.pos}
			currentSubsection.History = append(currentSubsection.History, newEntry)
			currentEntry = newEntry
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
			case kindPlain:
				state = entryContinuation
			default:
				return nil, p.normalizeError(
					fmt.Sprintf("unexpected line:%s\nexpecting subsection, version or change description", tok.fullText),
					tok.pos,
				)
			}
		case entryContinuation:
			currentEntry.Summary += " " + tok.fullText
			tok = <-tokens
			switch tok.kind {
			case kindSubsection:
				state = subsection
			case kindVersion:
				state = version
			case kindEntry:
				state = entry
			case kindPlain:
				state = entryContinuation
			case kindEOF:
				return result, nil
			default:
				return nil, p.normalizeError(
					fmt.Sprintf("unexpected line:%s\nexpecting subsection, version or change description", tok.fullText),
					tok.pos,
				)
			}
		}
	}
}

func (Default) retrieveLineKind(line string) tokenKind {
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

func (Default) normalizeError(msg string, pos int) error {
	return fmt.Errorf("%s (line %d)", msg, pos)
}
