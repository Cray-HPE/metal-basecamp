/*
 MIT License
 
 (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 
 Permission is hereby granted, free of charge, to any person obtaining a
 copy of this software and associated documentation files (the "Software"),
 to deal in the Software without restriction, including without limitation
 the rights to use, copy, modify, merge, publish, distribute, sublicense,
 and/or sell copies of the Software, and to permit persons to whom the
 Software is furnished to do so, subject to the following conditions:
 
 The above copyright notice and this permission notice shall be included
 in all copies or substantial portions of the Software.
 
 THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 OTHER DEALINGS IN THE SOFTWARE.
*/

package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	bcMiddleware "github.com/Cray-HPE/metal-basecamp/internal/middleware"
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
