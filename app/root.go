package app

import (
	"fmt"

	"github.com/b9o2/tabby"
)

type MainApp struct {
	*tabby.BaseApplication
}

func NewRootCommand(subApps ...tabby.Application) *MainApp {
	return &MainApp{
		tabby.NewBaseApplication(0, 0, subApps),
	}
}

func (r *MainApp) Detail() (string, string) {
	return "nym", "Node Version Manager written in Go"
}

func (r *MainApp) Init(parent tabby.Application) error {
	return nil
}

func (r *MainApp) Main(args tabby.Arguments) (*tabby.TabbyContainer, error) {
	fmt.Println("Welcome to nym - Node Version Manager")
	return nil, nil
}
