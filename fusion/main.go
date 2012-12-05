package main

import (
	"fusion"
	"golog"
	"os"
)

func main() {
	logger := golog.NewLogger("fusion")
	consoleProcessor := golog.NewConsoleProcessor(golog.LOG_INFO, true)
	logger.AddProcessor("console", consoleProcessor)
	bundler, err := fusion.NewQuickBundler(os.Args[1], logger)

	if err != nil {
		// We call println since fmt.Println might not be printed out.
		println("Error:", err.Error())
	} else {
		bundler.Run()
	}
}
