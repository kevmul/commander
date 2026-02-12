package template

import (
	"strings"
)

// Parser handles variable interpolation in strings
type Parser struct {
	variables map[string]string
}

// NewParser creates a new template parser
func NewParser() *Parser {
	return &Parser{
		variables: make(map[string]string),
	}
}

// Set sets a variable value
func (p *Parser) Set(key, value string) {
	p.variables[key] = value
}

// Get gets a variable value
func (p *Parser) Get(key string) (string, bool) {
	val, ok := p.variables[key]
	return val, ok
}

// Parse replaces all {{variable}} placeholders with their values
func (p *Parser) Parse(template string) string {
	result := template
	for key, value := range p.variables {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// Reset clears all variables
func (p *Parser) Reset() {
	p.variables = make(map[string]string)
}
