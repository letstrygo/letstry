# letstry

[![Sponsor Me!](https://img.shields.io/badge/%F0%9F%92%B8-Sponsor%20Me!-blue)](https://github.com/sponsors/nathan-fiscaletti)
[![Go Report Card](https://goreportcard.com/badge/github.com/letstrygo/letstry)](https://goreportcard.com/report/github.com/letstrygo/letstry)

**letstry** is a lightweight yet powerful tool designed to give developers **temporary workspaces** directly within their preferred IDE. Written in Go, it lets you spin up new projects quickly, save them as templates, and export them to a permanent location—**all from your VSCode terminal.** In addition to managing your own templates, letstry also supports a public repository for templates, making it easy to discover, share, and use templates from the community.

## Index

- [Installation](#installation)
- [Usage](#usage)
    - [Create a new session](#creating-a-new-session)
    - [Export a session](#exporting-a-session)
    - [List active sessions](#listing-active-sessions)
    - [Managing Templates](#managing-templates)
    - [Configuration](#configuration)
- [Contributing](#contributing)
- [Development](#development)

## Installation

Directly downloading **letstry** is the best option if you do not have [Golang](https://golang.org/dl/) installed on your system.

![Download Windows Installer](https://img.shields.io/badge/Windows%20Installer%20(x86_64)-blue?label=Download&color=7dccf0)
![Download macOS Installer](https://img.shields.io/badge/macOS%20Installer%20(arm64)-blue?label=Download&color=f2f2f7)
![Download Debian Installer](https://img.shields.io/badge/Debian%20Installer%20(x86_64)-blue?label=Download&color=d15a84)

> More download options are available on the [latest release](https://github.com/letstrygo/letstry/releases/latest) page.

### Install using GoLang

If you have [Golang](https://golang.org/dl/) already installed, you can install **letstry** more easily using the `go install` command.

```sh
go install github.com/letstrygo/letstry@latest
```

### Build from Source

If you wish to build it from source, you can clone the repository directly.

```powershell
git clone https://github.com/letstrygo/letstry.git
cd letstry
# Run directly
go run . [... arguments]
# Or, Build
go build -o output.exe .
```

> Note: **If you've installed letstry manually using Golang, or by compiling from source**; after installing you will need to manually enable the `lt` alias.
> 
> ```powershell
> # Windows Powershell
> "`nset-alias lt letstry" | out-file -append -encoding utf8 $profile; . $profile
> # Bash (or .zshrc, etc.)
> echo "alias lt='letstry'" >> ~/.bashrc && source ~/.bashrc
> ```

## Usage

### Creating a new Session

Creating a new session with letstry is simple and efficient. Use the `lt new` command to initialize a temporary project directory and open it in the default editor. This allows for quick prototyping. If you like the results, you can export the session to a more permanent location or save it as a template. 

```sh
$ lt new
```

If the VSCode window is closed, the temporary directory will be deleted. Therefore, you should either export your project using `lt export <path>` or save it as a template using `lt save <template-name>`.

Lets try sessions can be created from a directory path, a git repository URL, or a template name.

```sh
$ lt new <repository-url>
$ lt new <directory-path>
$ lt new <template-name>
```

### Exporting a Session

To export a session, use the `lt export` command from within the sessions directory. This will copy the session to the directory you specify.

```sh
$ lt export <path>
```

### Listing active sessions

To list all active sessions, use the `lt list` command.

```sh
$ lt list
```

### Managing Templates

**Creating a template**

Templates are a powerful feature of letstry. They allow you to save a project as a template and quickly create new projects based on that template.

To save an active session as a template, use the `lt save` command from within the sessions directory.

```sh
$ lt save [template-name]
```

If the session was initially created from an existing template, you can omit the name argument and the original template will be updated with the new session.

**Importing a Template**

You can easily import git repositories as templates using the `lt import` command.

```sh
$ lt import <template-name> <repository-url>
```

**Updating Templates**

If you've imported a template from a git repository using `lt import`, or if the template is stored as a git repository (i.e. contains a `.git` directory), you can use the `lt update` command to update the template with the latest version from it's associated git repository.

```sh
$ lt update <template-name>
```

**Listing Templates**

To list all available templates, use the `lt templates` command.

```sh
$ lt templates
```

**Deleting a Template**

To delete a template, use the `lt delete` command.

```sh
$ lt delete <template-name>
```

## Configuration

letstry can be configured using a configuration file. The configuration file is located at `~/.letstry/config.json`.

The config file allows you to specify different editors if you do not use VSCode.

**Windows Config Example**

`~/.letstry/config.json`
```json
{
    "default_editor": "vscode",
    "editors": [
        {
            "name": "vscode",
            "run_type": "run",
            "path": "C:\\Users\\natef\\AppData\\Local\\Programs\\Microsoft VS Code\\Code.exe",
            "args": "-n",
            "process_capture_delay": 2000000000,
            "tracking_type": "file_access"
        }
    ]
}
```

## Contributing

We welcome contributions to improve letstry. If you have suggestions or bug reports, please open an issue or submit a pull request.

## Development

To install letstry for development, run the following command from the root of the project:

```sh
$ go install ./
```

**Attaching a Debugger in VSCode**

Open the "Run and Debug" tab in VSCode (Ctrl+Shift+D on Windows) and select the `Run Letstry` configuration.

## License

This project is licensed under the MIT License.
