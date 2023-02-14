package internal

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const ServiceName = "dat(A)tom"

type version struct {
	major uint
	minor uint
	patch uint
}

func (v version) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

type Info struct {
	title       string
	description string
	v           version
}

func NewInfo(title, description string) *Info {
	return &Info{
		title:       strings.TrimSpace(title),
		description: strings.TrimSpace(description),
	}
}

func (i Info) String() string {
	parts := []string{
		fmt.Sprintf("%s v%s", ServiceName, i.Version()),
	}
	if i.title != "" {
		parts = append(parts, cases.Title(language.Und, cases.NoLower).String(i.title))
	}
	if i.description != "" {
		description := []rune(i.description)
		description[0] = unicode.ToUpper(description[0])
		parts = append(parts, string(description))
	}
	return strings.Join(parts, "\n")
}

func (i *Info) SetVersion(major, minor, patch uint) {
	i.v.major = major
	i.v.minor = minor
	i.v.patch = patch
}

func (i Info) Version() string {
	return i.v.String()
}

func (i Info) Title() string {
	return i.title
}

func (i Info) Description() string {
	return i.description
}
