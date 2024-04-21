package internal

import (
	"encoding/json"
	"log"
	"os"
)

type Languages struct {
	Names        []string          `json:"names"`
	CompileFiles map[string]string `json:"compileFiles"`
}

func newLanguages() *Languages {
	return &Languages{}
}

func (l *Languages) GetLanguages() []string {
	return l.Names
}

func (l *Languages) GetCompileFileName(language string) string {
	return l.CompileFiles[language]
}

func ParseLangauges(path string) (*Languages, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	log.Printf("DATA LENGTH: %v", len(data))
	config := newLanguages()
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

var LangaugesConfig *Languages
