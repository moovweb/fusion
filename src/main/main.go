package main

import(
	"fusion"
	"os"
)

func main() {
		
	bundler, err := fusion.NewQuickBundler(os.Args[1])		

	if err != nil {
		println("Error:", (*err).String() )
	} else {
		bundler.Run()
	}
	
}