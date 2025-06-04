package cli

var applicationHelpTemplate string = `
{{.Name}}: {{.ShortDescription}}

Description
    {{- $wrappedText := wrap .Description .Config.DescriptionMaxWidth 4}}

    {{$wrappedText}}

Usage

    {{.Name}} <command> [arguments]

Commands
{{$commands := commands -}}
{{$longestLen := longestStringLength $commands -}}
{{range $commands}}
    {{$padded := padEnd .Name $longestLen -}}
    {{white $padded}}  {{.ShortDescription}}{{end}}

Run '{{.Name}} help <command>' for more information on a command.

`
