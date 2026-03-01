package ssd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testUsername = "username"
	testApiKey   = "apiKey"

	testSsd = `{"status":"success","result":{"driveId":"1461","url":"https://www.techpowerup.com/ssd-specs/corsair-mp600-mini-1-tb.d1461","mfgr":"Corsair","name":"MP600 Mini","capacity":"1 TB","formFactor":"M.2 2280","interface":"PCIe 4.0 x4","protocol":"NVMe 1.4","dram":"N/A","hmb":"64 MB","released":"Apr 25th, 2023","endurance":"600 TBW","warranty":"5 Years","seqRead":"4,800 MB/s","seqWrite":"4,800 MB/s","controller":{"mfgr":"Phison","name":"PS5021-E21T","nameShort":"Phison E21T","channels":"4"},"flash":{"mfgr":"Micron","name":"B47R FortisFlash","type":"TLC","layers":"176-layer"}}}`

	testSsdList = `{"status":"success","result":[{"driveId":"1142","mfgr":"Magix","name":"Alpha EVO","capacity":"960","formFactor":"2.5\""},{"driveId":"1143","mfgr":"Magix","name":"Alpha EVO","capacity":"480","formFactor":"2.5\""},{"driveId":"1144","mfgr":"Magix","name":"Alpha EVO","capacity":"240","formFactor":"2.5\""},{"driveId":"1145","mfgr":"Magix","name":"Alpha EVO","capacity":"120","formFactor":"2.5\""}]}`

	testMagixSsd1142 = `{"status":"success","result":{"driveId":"1142","url":"https://www.techpowerup.com/ssd-specs/magix-alpha-evo.d1142","mfgr":"Magix","name":"Alpha EVO","capacity":"960 TB","formFactor":"2.5\"","interface":"SATA 3.2","protocol":"AHCI","dram":"N/A","hmb":"N/A","released":"Jan 1st, 2022","endurance":"300 TBW","warranty":"3 Years","seqRead":"550 MB/s","seqWrite":"500 MB/s","controller":{"mfgr":"Magix","name":"MX100","nameShort":"MX100","channels":"2"},"flash":{"mfgr":"Magix","name":"MagixFlash","type":"TLC","layers":"96-layer"}}}`

	testMagixSsd1143 = `{"status":"success","result":{"driveId":"1143","url":"https://www.techpowerup.com/ssd-specs/magix-alpha-evo.d1143","mfgr":"Magix","name":"Alpha EVO","capacity":"480 TB","formFactor":"2.5\"","interface":"SATA 3.2","protocol":"AHCI","dram":"N/A","hmb":"N/A","released":"Jan 1st, 2022","endurance":"300 TBW","warranty":"3 Years","seqRead":"550 MB/s","seqWrite":"500 MB/s","controller":{"mfgr":"Magix","name":"MX100","nameShort":"MX100","channels":"2"},"flash":{"mfgr":"Magix","name":"MagixFlash","type":"TLC","layers":"96-layer"}}}`

	testMagixSsd1144 = `{"status":"success","result":{"driveId":"1144","url":"https://www.techpowerup.com/ssd-specs/magix-alpha-evo.d1144","mfgr":"Magix","name":"Alpha EVO","capacity":"240 TB","formFactor":"2.5\"","interface":"SATA 3.2","protocol":"AHCI","dram":"N/A","hmb":"N/A","released":"Jan 1st, 2022","endurance":"300 TBW","warranty":"3 Years","seqRead":"550 MB/s","seqWrite":"500 MB/s","controller":{"mfgr":"Magix","name":"MX100","nameShort":"MX100","channels":"2"},"flash":{"mfgr":"Magix","name":"MagixFlash","type":"TLC","layers":"96-layer"}}}`

	testMagixSsd1145 = `{"status":"success","result":{"driveId":"1145","url":"https://www.techpowerup.com/ssd-specs/magix-alpha-evo.d1145","mfgr":"Magix","name":"Alpha EVO","capacity":"120 TB","formFactor":"2.5\"","interface":"SATA 3.2","protocol":"AHCI","dram":"N/A","hmb":"N/A","released":"Jan 1st, 2022","endurance":"300 TBW","warranty":"3 Years","seqRead":"550 MB/s","seqWrite":"500 MB/s","controller":{"mfgr":"Magix","name":"MX100","nameShort":"MX100","channels":"2"},"flash":{"mfgr":"Magix","name":"MagixFlash","type":"TLC","layers":"96-layer"}}}`
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/ssd-specs/api/%s/v1/query", testUsername) {
			// Return different SSDs based on the ID parameter
			id := r.URL.Query().Get("id")
			switch id {
			case "1461":
				w.Write([]byte(testSsd))
			case "1142":
				w.Write([]byte(testMagixSsd1142))
			case "1143":
				w.Write([]byte(testMagixSsd1143))
			case "1144":
				w.Write([]byte(testMagixSsd1144))
			case "1145":
				w.Write([]byte(testMagixSsd1145))
			}
		} else if r.URL.Path == fmt.Sprintf("/ssd-specs/api/%s/v1/lookup", testUsername) {
			w.Write([]byte(testSsdList))
		} else {
			http.NotFound(w, r)
		}
	}))
	return server
}

func TestTpuFindById(t *testing.T) {
	// Create a mock server that responds with JSON data
	server := setup(t)
	defer server.Close()

	tpu := NewTpuRepository(server.URL, "username", "apikey")
	ssd, err := tpu.FindById(context.Background(), "1461")
	if err != nil {
		t.Errorf("Error getting ssd data: %s", err)
	}

	if ssd.DriveID != "1461" {
		t.Errorf("Expected drive id to be 1461, got %s", ssd.DriveID)
	}

	if ssd.Controller.Manufacturer != "Phison" {
		t.Errorf("Expected controller manufacturer to be Phison, got %s", ssd.Controller.Manufacturer)
	}

	if ssd.Flash.Type != "TLC" {
		t.Errorf("Expected flash type to be TLC, got %s", ssd.Flash.Type)
	}
}

func TestTpuSearch(t *testing.T) {
	// Create a mock server that responds with JSON data
	server := setup(t)
	defer server.Close()

	tpu := NewTpuRepository(server.URL, "username", "apikey")
	ssds, err := tpu.Search(context.Background(), "search")
	if err != nil {
		t.Errorf("Error searching ssd data: %s", err)
	}

	if len(ssds) != 4 {
		t.Errorf("Expected length of ssds to be 4, got %v", len(ssds))
	}

	if ssds[0].Manufacturer != "Magix" {
		t.Errorf("Expected controller manufacturer to be Magix, got %s", ssds[0].Manufacturer)
	}

	if ssds[0].FormFactor != "2.5\"" {
		t.Errorf("Expected form factor to be 2.5\", got %s", ssds[0].FormFactor)
	}
}
