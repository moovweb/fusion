package main

import(
	"fusion/src/fusion"
	"os"
	"log4go"
)

func main() {
		
	bundler, err := fusion.NewQuickBundler(os.Args[1], make(log4go.Logger))		

	if err != nil {
		println("Error:", err.Error() )
	} else {
		bundler.Run()
	}
}
