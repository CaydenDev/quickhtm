package main

import (
	"strings"
)

type Parser struct {
	text   string
	pos    int
	result strings.Builder
}

func NewParser(text string) *Parser {
	return &Parser{text: text}
}

func (p *Parser) Parse() string {
	p.result.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"UTF-8\">\n<style>\nbody { font-family: system-ui; line-height: 1.5; max-width: 800px; margin: 0 auto; padding: 2rem; }\ncode { background: #f4f4f4; padding: 0.2em 0.4em; border-radius: 3px; }\npre code { display: block; padding: 1em; overflow-x: auto; }\n</style>\n</head>\n<body>\n")

	lines := strings.Split(p.text, "\n")
	inCodeBlock := false

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				p.result.WriteString("</code></pre>\n")
				inCodeBlock = false
			} else {
				p.result.WriteString("<pre><code>")
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			p.result.WriteString(line + "\n")
			continue
		}

		if line == "" {
			p.result.WriteString("\n")
			continue
		}

		if strings.HasPrefix(line, "#") {
			p.parseHeading(line)
		} else {
			p.result.WriteString("<p>")
			p.parseInline(line)
			p.result.WriteString("</p>\n")
		}
	}

	p.result.WriteString("</body>\n</html>")
	return p.result.String()
}

func (p *Parser) parseHeading(line string) {
	headingLevel := 1
	for i := 1; i < len(line); i++ {
		if line[i] == '#' {
			headingLevel++
		} else {
			break
		}
	}

	text := strings.TrimSpace(line[headingLevel:])
	p.result.WriteString("<h" + string(rune('0'+headingLevel)) + ">")
	p.parseInline(text)
	p.result.WriteString("</h" + string(rune('0'+headingLevel)) + ">\n")
}

func (p *Parser) parseInline(text string) {
	for i := 0; i < len(text); i++ {
		switch text[i] {
		case '`':
			end := strings.IndexByte(text[i+1:], '`')
			if end != -1 {
				p.result.WriteString("<code>")
				p.result.WriteString(text[i+1 : i+1+end])
				p.result.WriteString("</code>")
				i += end + 1
			} else {
				p.result.WriteByte(text[i])
			}

		case '*':
			if i+1 < len(text) && text[i+1] == '*' {
				end := strings.Index(text[i+2:], "**")
				if end != -1 {
					p.result.WriteString("<strong>")
					p.result.WriteString(text[i+2 : i+2+end])
					p.result.WriteString("</strong>")
					i += end + 3
				} else {
					p.result.WriteByte(text[i])
				}
			} else {
				end := strings.IndexByte(text[i+1:], '*')
				if end != -1 {
					p.result.WriteString("<em>")
					p.result.WriteString(text[i+1 : i+1+end])
					p.result.WriteString("</em>")
					i += end + 1
				} else {
					p.result.WriteByte(text[i])
				}
			}

		default:
			p.result.WriteByte(text[i])
		}
	}
}
