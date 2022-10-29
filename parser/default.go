package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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
	kindOther
	kindEOF
	kindError
)

type token struct {
	fullText string
	match    string
	kind     tokenKind
	pos      int
	err      error
}

func (p Default) Parse(r io.Reader, config any) (*model.Changelog, error) {
	patterns, err := p.defaultOrConfPatterns(config)
	if err != nil {
		return nil, fmt.Errorf("bad parser configuration: %v", err)
	}

	result := model.NewChangelog()
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	tokens := make(chan token)
	go func() {
		pos := 0
		for scanner.Scan() {
			pos++
			line := scanner.Text()
			lineKind, match, err := p.retrieveLineKind(line, patterns)
			if lineKind == kindEmpty {
				continue
			}
			tokens <- token{fullText: line, match: match, kind: lineKind, pos: pos, err: err}
		}
		tokens <- token{fullText: "", kind: kindEOF, pos: pos + 1}
		close(tokens)
	}()

	state := "INITIAL"
	tok := <-tokens
	var currentVersion *model.Version
	var currentSubsection *model.Subsection
	for {
		switch state {
		case "INITIAL":
			switch tok.kind {
			case kindTitle:
				state = "TITLE"
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line: %s\n expecting empty line or main title", tok.fullText)
			}
		case "TITLE":
			tok = <-tokens
			if tok.err != nil {
				return nil, tok.err
			}
			switch tok.kind {
			case kindOther:
				result.Header = append(result.Header, tok.fullText)
				tok = <-tokens
				if tok.err != nil {
					return nil, tok.err
				}
			case kindVersion:
				state = "VERSION"
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line: %s\nexpecting standard text line or version", tok.fullText)
			}
		case "VERSION":
			newVersion := &model.Version{Version: tok.match}
			result.Versions = append(result.Versions, newVersion)
			currentVersion = newVersion
			tok = <-tokens
			if tok.err != nil {
				return nil, tok.err
			}
			switch tok.kind {
			case kindVersion:
				state = "VERSION"
			case kindSubsection:
				state = "SUBSECTION"
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line:%s\nexpecting subsection or version", tok.fullText)
			}
		case "SUBSECTION":
			newSubsection := &model.Subsection{Name: tok.match}
			currentVersion.Subsections = append(currentVersion.Subsections, newSubsection)
			currentSubsection = newSubsection
			tok = <-tokens
			if tok.err != nil {
				return nil, tok.err
			}
			switch tok.kind {
			case kindSubsection:
				state = "SUBSECTION"
			case kindVersion:
				state = "VERSION"
			case kindEntry:
				state = "ENTRY"
			case kindEOF:
				return nil, fmt.Errorf("unexpected end of file")
			default:
				return nil, fmt.Errorf("unexpected line:%s\nexpecting subsection, version or change description", tok.fullText)
			}
		case "ENTRY":
			newEntry := model.ChangeLine{Summary: tok.match}
			currentSubsection.History = append(currentSubsection.History, &newEntry)
			tok = <-tokens
			if tok.err != nil {
				return nil, tok.err
			}
			switch tok.kind {
			case kindSubsection:
				state = "SUBSECTION"
			case kindVersion:
				state = "VERSION"
			case kindEntry:
				state = "ENTRY"
			case kindEOF:
				return result, nil
			default:
				return nil, fmt.Errorf("unexpected line:%s\nexpecting subsection, version or change description", tok.fullText)
			}
		}
	}
}

const mdSectionTitle = "title"
const mdSectionVersion = "version"
const mdSectionSubsection = "subsection"
const mdSectionEntry = "entry"
const mdSectionPlain = "plain"

func (p Default) defaultOrConfPatterns(config any) (map[string]*regexp.Regexp, error) {
	result := map[string]*regexp.Regexp{
		mdSectionTitle:      regexp.MustCompile(`^# (Changelog).*$`),
		mdSectionVersion:    regexp.MustCompile(`^## \[?(Unreleased|\d+\.\d+\.\d+)\]?( - \d{4}-\d\d-\d\d)?[ ]*$`),
		mdSectionSubsection: regexp.MustCompile(`^### ([A-Z][a-z]*).*$`),
		mdSectionEntry:      regexp.MustCompile(`^[\*|\-] (.+)$`),
		mdSectionPlain:      regexp.MustCompile(`^[A-Za-z].*$`),
	}

	userRegexp, ok := config.(map[string]string)
	if !ok {
		return nil, fmt.Errorf("expected map[string]string, got %T", config)
	}
	for k, v := range userRegexp {
		re, err := regexp.Compile(v)
		if err != nil {
			return nil, fmt.Errorf("invalid regular expression %q for %q: %v", v, k, err)
		}

		result[k] = re
	}

	return result, nil
}

func (p Default) match(re *regexp.Regexp, text string) (string, error) {
	matches := re.FindStringSubmatch(text)
	if len(matches) > 1 {
		return matches[1], nil
	}
	return "", fmt.Errorf("not enought (>1) matches: %v", matches)
}

func (p Default) retrieveLineKind(line string, patterns map[string]*regexp.Regexp) (tokenKind, string, error) {
	trimedLine := strings.Trim(line, " ")
	if strings.HasPrefix(trimedLine, "###") {
		re := patterns[mdSectionSubsection]
		match, err := p.match(re, line)
		if err != nil {
			return kindOther, match, fmt.Errorf("bad formated subsection line:%s\ndoes not match %s\n%v", line, re.String(), err)
		}
		return kindSubsection, match, nil
	}

	if strings.HasPrefix(trimedLine, "##") {
		re := patterns[mdSectionVersion]
		match, err := p.match(re, line)
		if err != nil {
			return kindOther, match, fmt.Errorf("bad formated version line:%s\ndoes not match %s\n%v", line, re.String(), err)
		}
		return kindVersion, match, nil
	}

	if strings.HasPrefix(trimedLine, "#") {
		re := patterns[mdSectionTitle]
		match, err := p.match(re, line)
		if err != nil {
			return kindOther, match, fmt.Errorf("bad formated title line:%s\ndoes not match %s\n%v", line, re.String(), err)
		}
		return kindTitle, match, nil
	}

	if strings.HasPrefix(trimedLine, "-") || strings.HasPrefix(trimedLine, "*") {
		re := patterns[mdSectionEntry]
		match, err := p.match(re, line)
		if err != nil {
			return kindOther, match, fmt.Errorf("bad formated change description line:%s\ndoes not match %s\n%v", line, re.String(), err)
		}
		return kindEntry, match, nil
	}

	if trimedLine == "" {
		return kindEmpty, "", nil
	}

	return kindOther, "", nil
}
