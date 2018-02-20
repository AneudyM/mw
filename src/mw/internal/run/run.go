package run

import (
	"mw/internal/cmd"
)

var CmdRun = &cmd.Command{
	CmdName:  "run",
	CmdUsage: "usage: mw run",
	Run:      testFunction,
}

func testFunction(c *cmd.Command, args []string) {
	println("You invoked the 'run' command.")
}

/*
func servePage(cmd *cobra.Command, args []string) {
	fmt.Println("Serving site at http://localhost:8080")

	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := cdp.New(ctxt, cdp.WithRunnerOptions(
		runner.Path("/opt/google/chrome/chrome"),
	))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Run(ctxt, goToSite())
	if err != nil {
		log.Fatal(err)
	}

	http.ListenAndServe(":8080", http.FileServer(http.Dir("./build")))

}

func init() {
	RootCmd.AddCommand(cmdServe)
}

func goToSite() cdp.Tasks {
	return cdp.Tasks{
		cdp.Navigate("localhost:8080"),
	}
}

func myWait() {
	time.Sleep(200 * time.Second)
}
*/
