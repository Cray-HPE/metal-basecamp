package middleware

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mostlygeek/arp"
	"github.com/spf13/viper"
)

// MacNotFoundError error thrown when mac is not found in arp table
type MacNotFoundError struct {
	Msg string
}

func (e *MacNotFoundError) Error() string {
	return fmt.Sprintf(e.Msg)
}

func getMac(c *gin.Context, defaultkey string) (string, error) {
	remoteaddr := c.Request.Header.Get("X-Forwarded-For")
	if remoteaddr == "" {
		remoteaddr = strings.Split(c.Request.RemoteAddr, ":")[0]
	}
	log.Printf("Updating arp cache")
	arp.CacheUpdate()
	mac := arp.Search(remoteaddr)
	if mac == "" {
		log.Printf("Could not find mac for ip '%s', returning '%s'", remoteaddr, defaultkey)
		return defaultkey, nil
	}
	return mac, nil
}

func getMetadataConfig(configfile string) (*viper.Viper, error) {
	dirname, filename := path.Split(configfile)
	extenstion := filepath.Ext(filename)
	name := strings.TrimSuffix(filename, extenstion)

	config := viper.New()
	config.SetConfigType(strings.TrimPrefix(extenstion, "."))
	config.SetConfigName(name)
	config.AddConfigPath(dirname)
	err := config.ReadInConfig()
	if err != nil {
		return config, err
	}
	config.WatchConfig()
	return config, nil
}

// LocalMetaData middleware to inject all data found for MAC within local files
func LocalMetaData(localdata, defaultkey, globalkey string) gin.HandlerFunc {
	metadata, err := getMetadataConfig(localdata)
	if err != nil {
		// We were unable to parse the data. This is fatal because we now have
		// no meta-data to serve.
		log.Fatal(err)
	}
	return func(c *gin.Context) {
		mac, err := getMac(c, defaultkey)
		if err != nil {
			c.Error(&AppError{Code: http.StatusNotFound, Message: err.Error()})
			c.Abort()
		}
		if !metadata.IsSet(mac) {
			c.Error(&AppError{Code: http.StatusNotFound, Message: fmt.Sprintf("Could not find metadata for mac '%s'", mac)})
			c.Abort()
		}

		md := metadata.GetStringMap(mac)
		gd := metadata.GetStringMap(globalkey)

		rd := make(map[string]interface{})

		m := md["meta-data"].(map[string]interface{})
		if m != nil && m["shasta-role"] != nil {
			rd = metadata.GetStringMap(m["shasta-role"].(string))
		}

		c.Set("basecamp_mac", mac)
		c.Set("basecamp_metadata", md)
		c.Set("basecamp_roledata", rd)
		c.Set("basecamp_globaldata", gd)

		c.Next()

		if c.GetBool("basecamp_write") {
			metadata.Set(mac, c.MustGet("basecamp_metadata"))
			metadata.WriteConfig()
			// For some reason it seems that if you WriteConfig, the watch is lost
			// TODO: Investigate
			metadata.WatchConfig()
		}
	}
}
