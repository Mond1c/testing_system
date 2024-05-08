package internal

import (
	"encoding/json"
	"os"
)

type Language struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Languages struct {
	Langs        []Language        `json:"languages"`
	CompileFiles map[string]string `json:"compileFiles"`
}

func newLanguages() *Languages {
	return &Languages{}
}

func (l *Languages) GetLanguages() []Language {
	return l.Langs
}

func (l *Languages) GetCompileFileName(language string) string {
	return l.CompileFiles[language]
}

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
