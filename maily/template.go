package maily

import (
	"bytes"
	htmlTemplate "html/template"
	"os"
	"path/filepath"
	textTemplate "text/template"
)

// A FuncMap is a list of functions that can be used within a template.
type FuncMap map[string]interface{}

// A TemplateData map holds information used in a template.
type TemplateData map[string]interface{}

type templateDot struct {
	ToEmail string
	Data    map[string]interface{}
}

func renderTextTemplate(filePath string, data templateDot, funcs FuncMap) (string, error) {
	baseName := filepath.Base(filePath)
	template, err := textTemplate.New(baseName).Funcs(textTemplate.FuncMap(funcs)).ParseFiles(filePath)
	if err != nil {
		return "", err
	}

	var outBuffer bytes.Buffer
	err = template.Execute(&outBuffer, data)
	if err != nil {
		return "", err
	}

	return outBuffer.String(), nil
}

func renderHTMLTemplate(basePath string, filePath string, data templateDot, funcs FuncMap) (string, error) {
	baseName := filepath.Base(filePath)
	template, err := htmlTemplate.New(baseName).Funcs(htmlTemplate.FuncMap(funcs)).ParseFiles(filePath, basePath)
	if err != nil {
		return "", err
	}

	var outBuffer bytes.Buffer
	err = template.Execute(&outBuffer, data)
	if err != nil {
		return "", err
	}

	return outBuffer.String(), nil
}

// render renders the template with the given data and returns the rendered subject, text, and HTML, in that order.
func (c *Context) render(toEmail string, templateName string, data TemplateData, textFuncs FuncMap, htmlFuncs FuncMap) (string, string, string, error) {
	baseFile := filepath.Join(c.TemplatePath, "base.html")
	templateFolder := filepath.Join(c.TemplatePath, templateName)

	subjectFile := filepath.Join(templateFolder, "subject.txt")
	textFile := filepath.Join(templateFolder, "template.txt")
	htmlFile := filepath.Join(templateFolder, "template.html")

	dot := templateDot{
		ToEmail: toEmail,
		Data:    data,
	}

	subject, err := renderTextTemplate(subjectFile, dot, textFuncs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", "", ErrTemplateMissingFile
		}

		return "", "", "", err
	}

	textBody, err := renderTextTemplate(textFile, dot, textFuncs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", "", ErrTemplateMissingFile
		}

		return "", "", "", err
	}

	htmlBody, err := renderHTMLTemplate(baseFile, htmlFile, dot, htmlFuncs)
	if err != nil {
		if os.IsNotExist(err) {
			return "", "", "", ErrTemplateMissingFile
		}

		return "", "", "", err
	}

	return subject, textBody, htmlBody, nil
}
