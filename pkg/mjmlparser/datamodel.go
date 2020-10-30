package mjmlparser

import (
	"bytes"
	"html/template"
)

func GenerateMJMLWithData(model string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("MJML").Parse(model)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
