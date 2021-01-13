package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	bcMiddleware "stash.us.cray.com/mtl/basecamp/internal/middleware"
	bcRoutes "stash.us.cray.com/mtl/basecamp/internal/routes"
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
