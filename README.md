# Frontmatter Replacment (fmr)
`fmr` allows you to take YAML frontmatter defined in a Markdown file, and replace content in the file with that frontmatter.

It also allows you to replace data in a different file if necessary, along with replacing the content in files using templates.

All templates are based on the formatting using the Go [`text/template`](https://pkg.go.dev/text/template) package.

This tool was developed to help generate live Markdown documentation from templates, where data is updated as it's acquired and the heavy use of task list items `- [ ]` to `- [X]` 

> [!NOTE]
> Replacements are simply text replacements and do not handle file formatting or specific formats at this time. This does mean that any file with plain text is supported though. Additionally, this does not currently utilize the `html/template` package for `.html` files
>
> This tool is a personal project, learning experience, and a work in progress. There are likely better ways to acheive some of what this does. Feel free to submit a PR and make suggestions!

# Example Usage
1. Define some frontmatter in a Markdown file
    ```yaml
    ---
    companyName: Water Corp
    customerId: 1234
    mainServer:
        fqdn: main.example.com
        ip: 10.0.0.1
    directoryServers:
        - dc1.example.com
        - dc2.example.com
    ---
    ```
2. Add some documentation in your file and add replacement placeholders
    ```
    1. Create a new entry in the client documentation
        ```
        Company: {{ .companyName }}
        Customer Id: {{ .customerId }}
        ```
    2. Configure the main server with the following information
        | Setting | Value                    |
        | ------- | ------------------------ |
        | fqdn    | `{{ .mainServer.fqdn }}` |
        | ip      | `{{ .mainServer.ip }}`   |
    ```

3. Run `fmr`
    ```
    fmr replace source --source WaterCorp.md
    ```

# Custom Template Functions
Some custom functions have been defined to help with some specific template needs

**`dn`**: Contructs a LDAP Distinguished Name. Takes up to 3 arguments
  - Domain name as FQDN e.g. `one.exmaple.com` converted to the `DC=` portion
  - (Optional) OU as FQDN e.g. `computers.org` converted to the `OU=` portion in the order specified e.g. `OU=computers,OU=org`
  - (Optional) CN as a normal value added as `CN=`

**`join`**: Maps to `strings.Join`

**`joinstr`**: Joins all provided strings as arguments by the delimiter. Format is `joinstr <delim> <parts>...`

**`json`**: Returns the JSON encoded string of the input.

**`lower`**: Maps to `strings.ToLower`

**`part`**: Returns the specific index of a string split by the delimiter. Format is `part <input> <index> <delim>`

**`replace`**: Performs a regex replacement using `regexp.ReplaceAllString`. Format is `replace <pattern> <sourceString> <replacement>`

**`shortFqdn`**: Returns the first part of a domain name. e.g. `one.example.com` -> `one`

**`title`**: Title cases the text using `language.Und`

**`trimprefix`**: Maps to `strings.TrimPrefix` with the argument order swapped to support the pipeline

**`trimsuffix`**: Maps to `strings.TrimSuffix` with the argument order swapped to support the pipeline

**`upper`**: Maps to `strings.ToUpper`

# Usage
```
fmr [global options] [command [command options]]
```

# Commands & flags
```
GLOBAL OPTIONS:
    --help, -h      show help    
    --version, -v   print the version

COMMANDS:
    replace         Replace data using frontmatter
    validate        Performs validation of files
    help, h         Shows a list of commands or help for one command
```

## `replace`
### `source`
Replace data directly in the source markdown file using frontmatter in the same file

| Flag(s)          | Description                                       |
| ---------------- | ------------------------------------------------- |
| `--source`, `-s` | File path to the source file with the frontmatter |

### `template`
Replace all data in the source file using the template file. This will merge the frontmatter between the two files with precedence given to the source file.

| Flag(s)                                    | Description                                                                                                                                                                                                                                                         |
| ------------------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `--source`, `-s`                           | File path to the source file with the frontmatter                                                                                                                                                                                                                   |
| `--template`, `-t`                         | File path to the template file                                                                                                                                                                                                                                      |
| `--[no-]retain-task-list-items`, `-[no-]r` | Defaults to true unless you specify the `no` version. This will search the source file for checked Markdown task list items (`- [X] Step Here`) and find a matching task list item in the template, and if found, re-check the item when replacing the source file. |

### `other`
Replace data in a non-markdown file, optionally from a template, using frontmatter defined in the source file.

If the template is a `.json` or `.jsonc` file, use `<<` and `>>` for the delimiters instead of `{{` and `}}` to avoid automatic formatting causing issues with your template.

| Flag(s)               | Description                                                                                                 |
| --------------------- | ----------------------------------------------------------------------------------------------------------- |
| `--source`, `-s`      | File path to the source file with the frontmatter                                                           |
| `--template`, `-t`    | (Optional) File path to the template file                                                                   |
| `--destination`, `-d` | File path to the destination file that will be replaced using the frontmatter, optionally from the template |

## `validate`
Performs validation of files

### `task-list-items`
Checks task list items in markdown files. Checks for items in both the source and template files that are not uniquely named and for items missing from the template that are checked in the source file. Both would currently would cause issues with re-checking task list items items.

| Flag(s)            | Description                                                                    |
| ------------------ | ------------------------------------------------------------------------------ |
| `--source`, `-s`   | File path to the source file with the frontmatter and checked task list items  |
| `--template`, `-t` | File path to the template file with the template and unchecked task list items |

# TODO
- [X] Tests
- [X] Version updates
- [X] Releases with artifacts built automatically
- [ ] Some refactor / cleanup