package main

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
	"github.com/watercorp/fmr/process"
)

type cliConfig struct {
	Source              string
	Template            string
	Destination         string
	RetainTaskListItems bool
}

// Create a new cli config
var cliConf = &cliConfig{}

func main() {
	// Define the commands
	cli := &cli.Command{
		Version:                "0.0.1",
		UseShortOptionHandling: true,
		Commands: []*cli.Command{
			{
				Name:  "replace",
				Usage: "Replace data using frontmatter",
				Commands: []*cli.Command{
					{
						Name:  "source",
						Usage: "Replace data directly in the source markdown file using frontmatter in the same file",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "source",
								Usage:       "File path to the source file with the frontmatter",
								TakesFile:   true,
								Aliases:     []string{"s"},
								Destination: &cliConf.Source,
								Required:    true,
							},
						},
						Action: processReplace,
					},
					{
						Name:  "template",
						Usage: "Replace data in the source markdown file, from a template markdown file, using frontmatter from the source, combined with the template",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "source",
								Usage:       "File path to the source file with the frontmatter",
								TakesFile:   true,
								Aliases:     []string{"s"},
								Destination: &cliConf.Source,
								Required:    true,
							},
							&cli.StringFlag{
								Name:        "template",
								Usage:       "File path to the template file",
								TakesFile:   true,
								Aliases:     []string{"t"},
								Destination: &cliConf.Template,
								Required:    true,
							},
							&cli.BoolWithInverseFlag{
								Name:        "retain-task-list-items",
								Usage:       "Retain checked task list items when replacing with a template",
								Aliases:     []string{"r"},
								Destination: &cliConf.RetainTaskListItems,
								Value:       true,
							},
						},
						Action: processReplace,
					},
					{
						Name:  "other",
						Usage: "Replace data in a non-markdown file, optionally from a template, using frontmatter from the source. .json and .jsonc template files use \"<<\" and \">>\" for replacemenmt delimiters instead of \"{{\" and \"}}\"",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "source",
								Usage:       "File path to the source file with the frontmatter",
								TakesFile:   true,
								Aliases:     []string{"s"},
								Destination: &cliConf.Source,
								Required:    true,
							},
							&cli.StringFlag{
								Name:        "template",
								Usage:       "Files path to the template file",
								TakesFile:   true,
								Aliases:     []string{"t"},
								Destination: &cliConf.Template,
							},
							&cli.StringFlag{
								Name:        "destination",
								Usage:       "File path to the destination file that will be replaced using the frontmatter, optionally from the template",
								TakesFile:   true,
								Aliases:     []string{"d"},
								Destination: &cliConf.Destination,
							},
						},
						Action: processOther,
					},
				},
			},
			{
				Name:  "validate",
				Usage: "Performs validation of files",
				Commands: []*cli.Command{
					{
						Name:  "task-list-items",
						Usage: "Checks markdown files for task list items",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:        "source",
								Usage:       "File path to the source file with the frontmatter and checked task list items",
								TakesFile:   true,
								Aliases:     []string{"s"},
								Destination: &cliConf.Source,
								Required:    true,
							},
							&cli.StringFlag{
								Name:        "template",
								Usage:       "File path to the template file with the template and unchecked task list items",
								TakesFile:   true,
								Aliases:     []string{"t"},
								Destination: &cliConf.Template,
								Required:    true,
							},
						},
						Action: processValidateTaskListItems,
					},
				},
			},
		},
	}

	if err := cli.Run(context.Background(), os.Args); err != nil {
		panic(err)
	}
}

func processReplace(ctx context.Context, cmd *cli.Command) error {
	// Check if we should process the source document itself
	// or replace the source using a template

	if cliConf.Template != "" {
		// Process the source and template
		err := process.ReplaceSourceFromTemplate(cliConf.Source, cliConf.Template, cliConf.RetainTaskListItems)
		if err != nil {
			return fmt.Errorf("error processing template from source: %w", err)
		}
	} else {
		// Process source file only
		err := process.ReplaceSourceOnly(cliConf.Source)
		if err != nil {
			return fmt.Errorf("error processing source: %w", err)
		}
	}

	return nil
}

func processOther(ctx context.Context, cmd *cli.Command) error {
	// Check if we need to replace the file using a template
	var err error
	if cliConf.Template != "" {
		err = process.ReplaceOther(cliConf.Source, cliConf.Destination, &cliConf.Template)
	} else {
		err = process.ReplaceOther(cliConf.Source, cliConf.Destination, nil)
	}

	if err != nil {
		return err
	}

	return nil
}

func processValidateTaskListItems(ctx context.Context, cmd *cli.Command) error {
	err := process.ValidateTaskListItems(cliConf.Source, cliConf.Template)
	if err != nil {
		return err
	}

	return nil
}
