package main

import (
	"nym/app"

	"github.com/b9o2/tabby"
)

func main() {

	ListVersion := app.NewListCommand()
	UseVersion := app.NewUseCommand()
	InitCommand := app.NewInitCommand()
	RootCommand := app.NewRootCommand(ListVersion, UseVersion, InitCommand)

	tabbyApp := tabby.NewTabby("nym", RootCommand)
	tabbyApp.Run(nil)
}
