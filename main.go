package main

import "MediaWarp/cmd"

func main() {
	cmd.PrintLOGO()
	cmd.InitFlag()

	defer cmd.ShutDown()
	cmd.SetUP()
}
