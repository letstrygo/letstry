# letstry

[![Sponsor Me!](https://img.shields.io/badge/%F0%9F%92%B8-Sponsor%20Me!-blue)](https://github.com/sponsors/nathan-fiscaletti)
[![Go Report Card](https://goreportcard.com/badge/github.com/letstrygo/letstry)](https://goreportcard.com/report/github.com/letstrygo/letstry)

**letstry** is a lightweight yet powerful tool designed to give developers **templated workspaces** directly within their preferred IDE. Written in Go, it lets you spin up new projects quickly, save them as templates, and export them to a permanent location—**all from your VSCode terminal.**

## Index

- [Installation](#installation)
- [Usage](#usage)
    - [Configuration](#configuration)
    - [Create a new session or project](#creating-a-new-session-or-project)
    - [Export a session](#exporting-a-session)
    - [List active sessions](#listing-active-sessions)
    - [Managing Templates](#managing-templates)
- [Contributing](#contributing)
- [Development](#development)

## Installation

**letstry** requires Go to be installed on your system. If you do not have Go installed, you can download it from the [official website](https://golang.org/dl/).

Once Go is installed, to install letstry, run the following command:

```sh
go install github.com/letstrygo/letstry@latest
```

### Optional: Configure `lt` alias

letstry is easier to use when you configure the `lt` alias. This allows you to type `lt` instead of typing out the full `letstry` command when you use it.

**Windows Powershell**
> Assuming you already have `$profile` configured
```powershell
"`nset-alias lt letstry" | out-file -append -encoding utf8 $profile; . $profile
```

**Bash**
```sh
echo "alias lt='letstry'" >> ~/.bashrc && source ~/.bashrc
```

## Usage

### Configuration

> [!TIP]\
> You can retrieve the path to the config file using the `lt path config` command. By default, the configuration is stored in `~/.letstry/config.json`.

By default, letstry is set-up as a **temporary workspace manager**. This means calls to `lt new` will result in a temporary workspace being created in your systems temporary directory that will be deleted once it's associated editor window is closed. This behavior can be customized using the `projects_path` and `require_export` configuration fields.

**Windows Config Example**

```jsonc
{
    // Projects Path
    //
    // The path in which to store projects. By default, letstry
    // will use a temporary directory.
    //
    // If no projects path is set, `require_export` will be
    // forcibly enabled. This is the default behavior.
    "projects_path": "",

    // Require Export
    //
    // When true:  New projects will be created as letstry sessions.
    //             These sessions will be automatically deleted once
    //             the editor window is closed.
    //
    // When false: New projects will be stored in `projects_path`
    //             and will be persisted after the editor window
    //             is closed. No letstry session will be created.
    //
    // You can force `require_export` by passing `--temp` to `lt new`.
    "require_export": true,

    // Default Editor
    //
    // The default editor to use for new sessions/projects.
    "default_editor": "vscode",

    // Editors
    //
    // Available editors for new sessions/projects.
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

### Creating a new Session or Project

Creating a new session or project with letstry is simple and efficient. Use the `lt new` command to initialize a new project or session and open it in the default editor.

```sh
$ lt new
```

Lets try sessions can be created from a directory path, a git repository URL, or a template name.

```sh
$ lt new <repository-url>
$ lt new <directory-path>
$ lt new <template-name>
```

> [!IMPORTANT]
> If `require_export` is enabled in your configuration or if you have not set a custom `projects_path`, when the VSCode window is closed the sessions temporary directory will be deleted. This is the default behavior for letstry. Therefore, you should either export your project using `lt export <path>` or save it as a template using `lt save <template-name>` (these commands must be run from within the sessions directory.)

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
