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
