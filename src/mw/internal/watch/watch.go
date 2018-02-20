package watch

import (
	"cmd/internal/build"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var cmdWatch = &cobra.Command{
	Use:   "watch",
	Short: "Start watching the filesystem",
	Long:  `Start watching the filesystem`,
	Run:   fileSystemWatch,
}

func fileSystemWatch(cmd *cobra.Command, args []string) {
	currDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Watching the filesystem...")
	fswatcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer fswatcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-fswatcher.Events:
				log.Println("Event: ", event.Op)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("Modified file: ", event.Name)
				}
			case err := <-fswatcher.Errors:
				log.Println("Error: ", err)
			}
		}
	}()

	err = fswatcher.Add(currDir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

func init() {
	RootCmd.AddCommand(cmdWatch)
}
