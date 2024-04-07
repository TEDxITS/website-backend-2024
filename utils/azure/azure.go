package azure

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func StopOnNewDeployment() {
	// creates a new file watcher for App_offline.htm
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR creating watcher:", err)
		return
	}
	defer watcher.Close()

	// get the current working directory and watch it
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("ERROR getting current working directory:", err)
		return
	}

	if err := watcher.Add(currentDir); err != nil {
		fmt.Println("ERROR adding directory to watcher:", err)
		return
	}

	// watch for App_offline.htm and exit the program if present
	// This allows continuous deployment on App Service as the .exe will not be
	// terminated otherwise
	for {
		select {

		case event := <-watcher.Events:

			if strings.HasSuffix(event.Name, "app_offline.htm") {
				fmt.Println("Exiting due to app_offline.htm being present")
				os.Exit(0)
			}

			watcher.Add(currentDir)
		}
	}
}
