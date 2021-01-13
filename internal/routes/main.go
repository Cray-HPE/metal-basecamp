package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	bcMiddleware "stash.us.cray.com/mtl/basecamp/internal/middleware"
)

const (
	//GLOBALKEY Global key name.
	GLOBALKEY = "Global"
	//DEFAULTKEY Default key name.
	DEFAULTKEY = "Default"
)

// Register middleware configuration.
func Register(router *gin.Engine, config *viper.Viper) {

	md := router.Group("/")

	localdata := config.GetString("local-data")
	if _, err := os.Stat(localdata); os.IsNotExist(err) {
		fmt.Printf("Local data path [%s] does not exist.\n", localdata)
		os.Exit(1)
	}

	// TODO: This feels wrong, look into refactoring to be idiomatic go
	var mw gin.IRoutes
	if config.GetBool("local-mode") {
		mw = md.Use(bcMiddleware.LocalMetaData(localdata, DEFAULTKEY, GLOBALKEY))
	} else {
		mw = md.Use(bcMiddleware.RemoteMetaData())
	}

	mw.GET("/meta-data", metaData)
	mw.GET("/user-data", userData)
	mw.POST("/phone-home/:iid", phoneHome)

	if config.GetBool("serve-static") {
		log.Println("Severing static files located here: ", config.GetString("static-dir"))
		router.StaticFS("/static", http.Dir(config.GetString("static-dir")))
	}
}
