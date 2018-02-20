package get

import (
	"fmt"

	"mw/internal/cmd"
	//"io/ioutil"
	//"log"
	//"os"
)

var CmdGet = &cmd.Command{
	CmdName:  "get",
	CmdUsage: "usage: mw get [library]",
	Run:      getLibrary,
}

// getLibrary is a test function
func getLibrary(c *cmd.Command, args []string) {
	fmt.Println("This command needs implementation:")
	fmt.Println("You invoked the 'get' command")
	fmt.Println("With arguments", args)
}

// GetNew is just a test as well
// name: is the name of the file as a string
// path: is the absolute path to the file as a string
func GetNew(nanme string, path string) {

}

/*
func cmdGetLibrary(cmd *cobra.Command, args []string) {
	if len(args) <= 0 {
		fmt.Println("You need to specify a project name.")
		os.Exit(1)
	} else if len(args) > 1 {
		fmt.Println("Specify only one file.")
		os.Exit(1)
	}
	libraryName := args[0]
	fmt.Println(libraryName)
	htmlFile, err := ioutil.ReadFile(libraryName)
	if err != nil {
		log.Fatal(err)
	}
	ioutil.WriteFile("netFile.html", htmlFile, 0666)
	fmt.Println(string(htmlFile), "Hello")
}

func init() {
	RootCmd.AddCommand(cmdGet)
}
*/
