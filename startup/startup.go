package startup

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/gammazero/workerpool"
)

func Migrate() {
	// Open the database connection
	db := database.Open()

	database.MigrateUp(db)
}

func registerDependencies() *controller.ApiContainer {
	// Open the database connection
	db := database.Open()

	return internal.InitializeContainer(db)
}

func Execute() {
	container := registerDependencies()

	wp := workerpool.New(2)

	wp.Submit(container.HttpServer.Run)

	wp.Submit(container.WebsocketServer.Run)

	wp.StopWait()
}
