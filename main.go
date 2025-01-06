package main

import (
	"os"

	_ "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/docs"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/startup"
)

// @BasePath /api/v1
func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate-up" {
		startup.Migrate()
		return
	}

	startup.Execute()
}
