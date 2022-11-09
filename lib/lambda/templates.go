package lambda

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Template struct {
	Arguments map[string]string `json:"arguments"`
	Lambdas   map[string]string `json:"lambda_functions"`
}

func ReadTemplate(name string) (*Template, error) {
	templatePath := filepath.Join(os.Getenv("LAMBDA_TEMPLATE_DIR"), name)
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return nil, fmt.Errorf("unable to read lambda template: %w", err)
	}

	tpl := &Template{}
	if err := json.Unmarshal(content, tpl); err != nil {
		return nil, fmt.Errorf("unable to parse template as json: %w", err)
	}

	return tpl, nil
}

func WriteTemplate(name string, tpl *Template) error {
	content, err := json.Marshal(tpl)
	if err != nil {
		return fmt.Errorf("error creating json string: %w", err)
	}

	templatePath := filepath.Join(os.Getenv("LAMBDA_TEMPLATE_DIR"), name)
	if err := ioutil.WriteFile(templatePath, content, os.ModePerm); err != nil {
		return fmt.Errorf("error writing template file: %w", err)
	}

	return nil
}
