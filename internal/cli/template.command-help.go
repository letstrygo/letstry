package cli

var commandHelpTemplate string = `
{{getCallerName}}: {{whiteCommand .Name}} — {{.ShortDescription}}

Description
    {{- $wrappedText := wrap .Description getMaxWidth 4}}

    {{$wrappedText}}

Usage

    {{white getCallerName}} {{whiteCommand .Name}} {{range .Arguments}}{{if .Required}}<{{.Name}}>{{else}}[{{.Name}}]{{end}} {{end}}
{{if .Arguments}}

Arguments
{{range .Arguments}}
    {{- $argDesc := wrap .Description getMaxWidth 4}}
    {{.Label}}

    {{$argDesc}}
{{end}}
{{end}}
Run '{{getCallerName}} help' for information on additional commands.

`
