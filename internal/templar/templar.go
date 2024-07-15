package templar

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/yaml.v3"
)

// ExecuteItemTemplate executes the template for a single item
func GetParsedTemplate(tstring string, vars map[string]interface{}) (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("template").Parse(tstring)
	if err != nil {
		return "", err
	}
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", fmt.Errorf("failed to execute item template: %v", err)
	}
	return buf.String(), nil
}

func executeTemplate(pattern string, vars map[string]interface{}, with_items interface{}) error {
	// Handle cases where with_items is a template itself
	if str, ok := with_items.(string); ok {
		template_pattern, err := template.New("template_pattern").Parse(str)
		if err != nil {
			return fmt.Errorf("failed to parse with_items template: %v", err)
		}
		vars_buf := new(bytes.Buffer)
		if err := template_pattern.Execute(vars_buf, vars); err != nil {
			return fmt.Errorf("failed to execute with_items template: %v", err)
		}
		template_string := vars_buf.String()
		err = yaml.Unmarshal([]byte(template_string), &with_items)
		if err != nil {
			return fmt.Errorf("failed to unmarshal with_items template result: %v", err)
		}
	}

	// Process the items
	// switch items := with_items.(type) {
	// case []interface{}:
	// 	for _, item := range items {
	// 		if _, err := executeItemTemplate(pattern, vars); err != nil {
	// 			return err
	// 		}
	// 	}
	// case interface{}:
	// 	if _, err := executeItemTemplate(pattern, vars); err != nil {
	// 		return err
	// 	}
	// default:
	// 	return fmt.Errorf("unexpected type for with_items: %T", with_items)
	// }

	return nil
}
