package ssd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testUsername = "username"
	testApiKey   = "apiKey"

	testSsd  = `{"status":"success","result":{"driveId":"1461","url":"https://www.techpowerup.com/ssd-specs/corsair-mp600-mini-1-tb.d1461","mfgr":"Corsair","name":"MP600 Mini","capacity":"1 TB","formFactor":"M.2 2280","interface":"PCIe 4.0 x4","protocol":"NVMe 1.4","dram":"N/A","hmb":"64 MB","released":"Apr 25th, 2023","endurance":"600 TBW","warranty":"5 Years","seqRead":"4,800 MB/s","seqWrite":"4,800 MB/s","controller":{"mfgr":"Phison","name":"PS5021-E21T","nameShort":"Phison E21T","channels":"4"},"flash":{"mfgr":"Micron","name":"B47R FortisFlash","type":"TLC","layers":"176-layer"}}}`
	testSsd2 = `{"status":"success","result":{"driveId":"123","url":"https://www.techpowerup.com/ssd-specs/xpg-gammix-s70-2-tb.d123","mfgr":"XPG","name":"Gammix S70","capacity":"2 TB","formFactor":"M.2 2280","interface":"PCIe 4.0 x4","protocol":"NVMe 1.4","dram":"2048 MB","hmb":"Unknown","released":"Sep 2020","endurance":"1480 TBW","warranty":"5 Years","seqRead":"7,400 MB/s","seqWrite":"6,400 MB/s","controller":{"mfgr":"InnoGrit","name":"IG5236 (Rainier)","nameShort":"IG5236","channels":"8"},"flash":{"mfgr":"Micron","name":"B27B FortisFlash","type":"TLC","layers":"96-layer"}}}`
)

func setup(t *testing.T) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == fmt.Sprintf("/ssd-specs/api/%s/v1/query", testUsername) {
			w.Write([]byte(testSsd))
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

	tpu := NewTpuSSDRepository(server.URL, "username", "apikey")
	// Call the getUserData function with the mock server URL
	ssd, err := tpu.FindById("1461")
	if err != nil {
		t.Errorf("Error getting ssd data: %s", err)
	}

	// Check that the returned user has the expected values
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
