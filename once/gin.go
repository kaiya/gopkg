package once

import "github.com/gin-gonic/gin"

func CreateApp() (app *gin.Engine) {
	//gin
	app = gin.New()
	return
}
