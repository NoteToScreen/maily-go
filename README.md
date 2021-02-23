# maily-go
[![CI status](https://github.com/NoteToScreen/maily-go/workflows/CI/badge.svg)](https://github.com/NoteToScreen/maily-go/actions)

Maily is a simple Go library for sending templated email messages. The built-in `text/template` and `html/template` packages are used for templating, and `net/smtp` is used to send the actual email. Maily takes care of evaluating your templates, building a MIME multipart message, and sending your email.

## Usage
To use Maily, you must first create a `Context`, which contains all of the options you can configure.

An example is below:
```go
context := maily.Context{
	FromAddress: "sample@example.com",
	FromDisplay: "Test Sender <sample@example.com>",

	SendDomain: "example.com",

	SMTPHost:     "smtp.example.com",
	SMTPPort:     25,
	SMTPUsername: "sample@example.com",
	SMTPPassword: "password123",

	TemplatePath: "templates/", // path can be relative or absolute
}
```

The `TemplatePath` must be a path to a folder which contains all of your email templates. This folder must contain a file, `base.html`, which is included in all HTML templates, and a folder for each template. A template's folder must contain three files:

* `subject.txt` - A text template which is evaluated to create the email's subject line.
* `template.txt` - A text template which is evaluated to create the email's text body.
* `template.html` - An HTML template which is evaluated to create the email's HTML body.

Then, you can send an email like so (in this example, the folder containing the template files is called `importantMessage`):
```go
func getFirstName(name string) string {
	return strings.Split(name, " ")[0]
}

templateData := maily.TemplateData{
	"nameOfPerson": "John Q. Smith",
	"message": "Hello this is important.",
	"randomParameter": 42,
	"isImportant": true,
}
funcs := maily.FuncMap{
	"fname": getFirstName,
}
result, err := context.SendMail("John Smith", "john_smith@example.com", "importantMessage", templateData, funcs, funcs)
```

`templateData` is a `TemplateData` map (really a `map[string]interface{}`) which is passed to the template. `funcs` is a `FuncMap` which contain functions that can be used from templates.

For example, if a template file contained `{{ .Data.nameOfPerson | fname }}`, this would be evaluated to `John` in the above example. See the [text/template package documentation](https://golang.org/pkg/text/template/) for more details about templating.

`result` is a struct of type `EmailResult`. Currently, it only contains `MessageID`, which lets you know the value of the `Message-ID` header in the sent email. This can be useful if you're building an application where it's important to track replies.