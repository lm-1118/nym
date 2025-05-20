package app

import (
	"fmt"
	"nym/core/node"
	"os"

	"github.com/b9o2/tabby"
)

type ListCommand struct {
	*tabby.BaseApplication
	ma *MainApp
}

func NewListCommand(subApps ...tabby.Application) *ListCommand {
	return &ListCommand{
		tabby.NewBaseApplication(0, 0, subApps),
		nil,
	}
}

func (l *ListCommand) Detail() (string, string) {
	return "list", "List all installed Node.js versions"
}

func (l *ListCommand) Init(parent tabby.Application) error {
	return nil
}

func (l *ListCommand) Main(args tabby.Arguments) (*tabby.TabbyContainer, error) {
	versions, err := node.ListInstalledVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing versions: %v\n", err)
		os.Exit(1)
	}

	currentVersion, _ := node.GetCurrentVersion()

	fmt.Println("Installed Node.js versions:")
	for _, v := range versions {
		if currentVersion == v {
			fmt.Printf("=> v%s (current)\n", v)
		} else {
			fmt.Printf("   v%s\n", v)
		}
	}

	return nil, nil
}
