package new

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"mw/internal/base"
	"mw/internal/cmd"
)

var CmdNew = &cmd.Command{
	CmdName:    "new",
	CmdUsage:   "usage: mw new <project-name>",
	HasNoFlags: true,
	Run:        createProject,
}

var srcDir = base.MWROOT[base.DIR_SRC]

const projectConfig = `# {{.ProjectName}} project configuration file
projectName="{{.ProjectName}}"
projectRoot="{{.ProjectRoot}}"
`

const indexTemplate = `<!DOCTYPE html>
<html class="no-js" lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="x-ua-compatible" content="ie=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="generator" content="Mint Web 0.0.2" />
    <title>Mint Web Page</title>
    <link rel="stylesheet" href="css/main.css">
  </head>
  <body>
    <h1>Hello, Mint Web!</h1>

    <script src="js/main.js"></script>
  </body>
</html>
`

func createProject(c *cmd.Command, args []string) {
	var project base.Project

	if len(args) == 0 {
		log.Fatalf("No project name specified.")
	}

	if len(args) > 1 {
		fmt.Printf(c.CmdUsage + "\n")
		os.Exit(2)
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	project.ProjectName = args[0]
	project.ProjectRoot = filepath.Join(cwd, project.ProjectName)

	err = os.Mkdir(project.ProjectName, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := os.Create(filepath.Join(project.ProjectRoot, base.MWROOT[base.CFG_FILE]))
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("config").Parse(projectConfig)
	if err != nil {
		log.Fatal(err)
	}

	err = tmpl.Execute(cfg, project)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(project.ProjectName, srcDir), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir(filepath.Join(project.ProjectName, base.MWROOT[base.DIR_BUILD]), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	for _, dir := range base.MWSRC {
		err := os.Mkdir(filepath.Join(project.ProjectName, srcDir, dir), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
	}

	file := filepath.Join(project.ProjectName, srcDir, "index.html")
	cf, err := os.Create(file)
	defer cf.Close()
	if err != nil {
		log.Fatal(err)
	}

	_, err = cf.WriteString(indexTemplate)
	if err != nil {
		log.Fatal(err)
	}
}
