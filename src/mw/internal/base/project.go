package base

type Project struct {
	ProjectName string
	ProjectRoot string
}

type RootDir int
type SrcDir int

const (
	DIR_SRC RootDir = iota
	DIR_BUILD
	CFG_FILE
)
const (
	DIR_JS SrcDir = iota
	DIR_IMG
	DIR_CSS
	DIR_SECTIONS
	DIR_TEMPLATES
)

var MWROOT = [...]string{
	DIR_SRC:   "src",
	DIR_BUILD: "build",
	CFG_FILE:  ".mwrc",
}

var MWSRC = [...]string{
	DIR_CSS:       "css",
	DIR_IMG:       "img",
	DIR_JS:        "js",
	DIR_SECTIONS:  "sections",
	DIR_TEMPLATES: "templates",
}

func (p *Project) Name() string {
	return p.ProjectName
}

func (p *Project) RootDir() string {
	return p.ProjectRoot
}
