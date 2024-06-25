package config

import (
	"encoding/json"
	"os"
)


// Langauge represents information about programming language
// Name is a full name of programming language
// Value is a short name of programming language
// For example: Name = "C++ 20 (gcc)", Value = "cpp20"
type Language struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Languages represents information about all using programming languages in the contest
// Langs is array of Language with base information about each language
// CompileFiles is map that associate short name of the programming language with path to the compilation shell script
type Languages struct {
	Langs        []Language        `json:"languages"`
	CompileFiles map[string]string `json:"compileFiles"`
}

// newLanguages create pointer to empty Languages struct
func newLanguages() *Languages {
	return &Languages{}
}

// GetLangauges returns array of basic information about programming languages in the contest as array of Language
func (l *Languages) GetLanguages() []Language {
	return l.Langs
}

// GetCompileFileName returns path to the compilation shell script for the specified programming languages
// language specifies programming language
func (l *Languages) GetCompileFileName(language string) string {
	return l.CompileFiles[language]
}

// ParseLanguages parses Languages from the specifed json file.
// returns pointer to Languages 
func ParseLangauges(path string) (*Languages, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := newLanguages()
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

var LangaugesConfig *Languages
