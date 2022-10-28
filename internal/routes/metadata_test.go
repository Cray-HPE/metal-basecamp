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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func Test_mapLookup(t *testing.T) {
	testData := map[string]interface{}{
		"Global": map[string]interface{}{
			"meta-data": map[string]interface{}{
				"one": "global-1",
				"two": "global-2",
			},
			"user-data": map[string]interface{}{},
		},
		"Default": map[string]interface{}{
			"meta-data": map[string]interface{}{
				"one": "default-1",
				"two": "default-2",
			},
			"user-data": map[string]interface{}{},
		},
		"Storage": map[string]interface{}{
			"meta-data": map[string]interface{}{
				"one": "storage-1",
				"two": "storage-2",
			},
			"user-data": map[string]interface{}{},
		},
		"Nodes": map[string]interface{}{
			"f0:18:98:8c:ea:7d": map[string]interface{}{
				"meta-data": map[string]interface{}{
					"one": "mac-defined-1",
					"two": "mac-defined-2",
				},
				"user-data": map[string]interface{}{},
			},
		},
		"f0:18:98:8c:ea:7e": map[string]interface{}{
			"meta-data": map[string]interface{}{
				"one": "mac-defined-1",
				"two": "mac-defined-2",
			},
			"user-data": map[string]interface{}{},
		},
	}
	// Test that we can find a nested key
	result, err := mapLookup(testData, "Nodes", "f0:18:98:8c:ea:7d")
	if err != nil {
		t.Errorf("Encountered unexpected error in mapLookup: %s", err)
	} else {
		resultMap := result.(map[string]interface{})
		metaDataMap := resultMap["meta-data"].(map[string]interface{})
		if metaDataMap["one"] != "mac-defined-1" {
			t.Errorf("Got unexpected results from mapLookup: %v", result)
		}
	}

	// Test that we can find an un-nested key
	result, err = mapLookup(testData, "f0:18:98:8c:ea:7e")
	if err != nil {
		t.Errorf("Encountered unexpected error in mapLookup: %s", err)
	} else {
		resultMap := result.(map[string]interface{})
		metaDataMap := resultMap["meta-data"].(map[string]interface{})
		if metaDataMap["one"] != "mac-defined-1" {
			t.Errorf("Got unexpected results from mapLookup: %v", result)
		}
	}

}

type baseData struct {
	basecampMetadata   map[string]interface{}
	basecampGlobaldata map[string]interface{}
	basecampRoledata   map[string]interface{}
}

func testMetaDataGeneral(t *testing.T, route string, payload baseData) {
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Set("basecamp_metadata", payload.basecampMetadata)
	c.Set("basecamp_globaldata", payload.basecampGlobaldata)
	c.Set("basecamp_roledata", payload.basecampRoledata)
	req, _ := http.NewRequest("GET", route, nil)
	c.Request = req
	metaData(c)
	if w.Code != http.StatusOK {
		t.Error("Response was not OK: ", w.Code)
	}
}

func TestMetaDataBare(t *testing.T) {

	var payload baseData

	payload.basecampGlobaldata = make(map[string]interface{})
	payload.basecampMetadata = make(map[string]interface{})
	payload.basecampRoledata = make(map[string]interface{})

	testMetaDataGeneral(t, "/user-data/", payload)

}

func TestMetaDataWithGlobalData(t *testing.T) {
	var payload baseData

	payload.basecampGlobaldata = make(map[string]interface{})
	payload.basecampMetadata = make(map[string]interface{})
	payload.basecampRoledata = make(map[string]interface{})

	payload.basecampGlobaldata["meta-data"] = make(map[string]string)
	payload.basecampGlobaldata["meta-data"].(map[string]string)["region"] = "US-East-1"

	testMetaDataGeneral(t, "/user-data/", payload)
}
