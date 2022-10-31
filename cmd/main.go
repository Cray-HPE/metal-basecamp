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

package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	bcMiddleware "github.com/Cray-HPE/metal-basecamp/internal/middleware"
	bcRoutes "github.com/Cray-HPE/metal-basecamp/internal/routes"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	config := viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("server")
	config.AddConfigPath("./configs")
	err := config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	config.WatchConfig()

	flag.String("bind", ":8888", "Address to bind on")
	flag.Bool("local-mode", true, "Serve in local mode")
	flag.Bool("serve-static", true, "Serve in static files")
	flag.Bool("init-local", false, "When using remote, initialize with local data?")
	flag.String("local-data", "", "Path of local data file")
	flag.String("static-dir", "", "Path to static files to serve")
	flag.String("remote-endpoint", "", "Remote backend url")
	flag.String("remote-creds", "", "Path to file holding remote credentials")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()
	config.BindPFlags(pflag.CommandLine)

	router := gin.Default()
	router.Use(bcMiddleware.AppErrorReporter())

	// Register Routes
	bcRoutes.Register(router, config)

	router.Run(config.GetString("bind"))
}
