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

// NOTE: All interactions with a backend should occur in Middleware
// to faciliate the ability to easily pivot from local to remote backends

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	yaml "gopkg.in/yaml.v2"
)

const querykey = "key"

// Merge values from second into first. We will only handle nested maps,
// slices will always favor second over first.
func mergeMaps(first, second map[string]interface{}) map[string]interface{} {
	for key, secondVal := range second {
		if firstVal, present := first[key]; present {
			switch firstVal.(type) {
			case map[string]interface{}:
				// value is also a map interface, so recurse into it
				first[key] = mergeMaps(firstVal.(map[string]interface{}), secondVal.(map[string]interface{}))
				continue
			default:
				first[key] = secondVal
			}
		} else {
			// key not in first so add it
			first[key] = secondVal
		}
	}
	return first
}

func mapLookup(m map[string]interface{}, keys ...string) (interface{}, error) {
	var ok bool
	var foundObj interface{}

	if len(keys) == 0 {
		return nil, fmt.Errorf("mapLookup needs at least one key")
	}
	if foundObj, ok = m[keys[0]]; !ok {
		return nil, fmt.Errorf("key not found in map; keys: %v", keys)
	} else if len(keys) == 1 {
		return foundObj, nil
	} else if m, ok = foundObj.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("malformed structure at %#v", foundObj)
	}
	return mapLookup(m, keys[1:]...)
}

func generateInstanceID() string {
	b := make([]byte, 6)
	rand.Read(b)
	return strings.ToLower(fmt.Sprintf("i-%X", b))
}

// InstanceInfo Various meta for the basecamp instance.
type InstanceInfo struct {
	PublicKeyDSA   string `form:"pub_key_dsa" json:"pub_key_dsa" binding:"required"`
	PublicKeyRSA   string `form:"pub_key_rsa" json:"pub_key_rsa" binding:"required"`
	PublicKeyECDSA string `form:"pub_key_ecdsa" json:"pub_key_ecdsa" binding:"required"`
	InstanceID     string `form:"instance_id" json:"instance_id" binding:"required"`
	Hostname       string `form:"hostname" json:"hostname" binding:"required"`
}

func phoneHome(c *gin.Context) {
	var info InstanceInfo
	mac := c.MustGet("basecamp_mac").(string)
	if mac == DEFAULTKEY {
		c.String(http.StatusBadRequest, "Ignoring default mac.")
		return
	}

	c.Bind(&info)
	log.Printf("Saving public info for %s with an iid of %s", mac, string(info.InstanceID))
	config := c.MustGet("basecamp_metadata").(map[string]interface{})
	config["phone-home"] = info
	c.Set("basecamp_metadata", config)
	c.Set("basecamp_write", true)

	c.String(http.StatusOK, "ok")
}

func metaData(c *gin.Context) {
	config := c.MustGet("basecamp_metadata").(map[string]interface{})
	globalconfig := c.MustGet("basecamp_globaldata").(map[string]interface{})
	roleconfig := c.MustGet("basecamp_roledata").(map[string]interface{})

	data := make(map[string]interface{})
	if config["meta-data"] != nil {
		data = config["meta-data"].(map[string]interface{})
	}

	roledata := make(map[string]interface{})
	if roleconfig["user-data"] != nil {
		roledata = roleconfig["user-data"].(map[string]interface{})
	}

	mergedData := mergeMaps(roledata, data)

	log.Printf("role: %v", roleconfig)
	queries := c.Request.URL.Query()
	// Industry seems to use the same iid for the life of the VM
	// Since we don't have VMs, we are generating a new iid every boot.
	// This will force the VM to rerun all the user-data each boot, picking
	// up any changes.
	mergedData["instance-id"] = generateInstanceID()

	// Add global data before filtering a query request
	mergedData[GLOBALKEY] = globalconfig["meta-data"]

	lookupKeys, ok := queries[querykey]
	if ok && len(lookupKeys) > 0 {
		// Query string provided in request, return it.
		lookupKey := strings.Split(lookupKeys[0], ".")
		rval, err := mapLookup(mergedData, lookupKey...)
		if err != nil {
			c.String(http.StatusNotFound, "Not Found")
			return
		}
		c.JSON(http.StatusOK, rval)
	} else {
		// No query, return all data
		c.JSON(http.StatusOK, mergedData)
	}
}

func userData(c *gin.Context) {

	config := c.MustGet("basecamp_metadata").(map[string]interface{})
	roleconfig := c.MustGet("basecamp_roledata").(map[string]interface{})

	data := make(map[string]interface{})
	if config["user-data"] != nil {
		data = config["user-data"].(map[string]interface{})
	}
	roledata := make(map[string]interface{})
	if roleconfig["user-data"] != nil {
		roledata = roleconfig["user-data"].(map[string]interface{})
	}

	mergedData := mergeMaps(roledata, data)

	databytes, err := yaml.Marshal(mergedData)
	if err != nil {
		c.Error(err)
		return
	}
	final := "#cloud-config\n" + string(databytes)
	c.Data(http.StatusOK, "text/yaml", []byte(final))
}
