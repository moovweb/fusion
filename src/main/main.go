package main

import(
	"fusion"
	"os"
)

func main() {
		
	bundler := fusion.NewQuickBundler(os.Args[1])		
	bundler.Run()
	
}