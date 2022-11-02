package main

import (
	"edufund/controllers"
)

func main() {
	app := &controllers.App{}
	app.Initialize()
	app.Run(":" + controllers.Env.HTTPPort)
}
