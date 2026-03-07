package ssd

import (
	"testing"
)

func TestParseCapacity(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantCapacity int
		wantOk       bool
	}{
		{
			name:         "TB capacity",
			input:        "Samsung 970 EVO 1TB",
			wantCapacity: 1,
			wantOk:       true,
		},
		{
			name:         "GB capacity",
			input:        "Crucial MX500 500GB",
			wantCapacity: 500,
			wantOk:       true,
		},
		{
			name:         "2TB capacity",
			input:        "WD Blue 2TB NVMe SSD",
			wantCapacity: 2,
			wantOk:       true,
		},
		{
			name:         "no capacity",
			input:        "Samsung 970 EVO",
			wantCapacity: 0,
			wantOk:       false,
		},
		{
			name:         "lowercase tb",
			input:        "samsung 970 evo 1tb",
			wantCapacity: 1,
			wantOk:       true,
		},
		{
			name:         "with space before unit",
			input:        "Samsung 970 EVO 1 TB",
			wantCapacity: 1,
			wantOk:       true,
		},
		{
			name:         "multiple numbers only first match",
			input:        "Samsung 970 1TB for 100 dollars",
			wantCapacity: 1,
			wantOk:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			capacity, ok := parseCapacity(tt.input)
			if ok != tt.wantOk {
				t.Errorf("parseCapacity(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if capacity != tt.wantCapacity {
				t.Errorf("parseCapacity(%q) capacity = %d, want %d", tt.input, capacity, tt.wantCapacity)
			}
		})
	}
}

func TestParseFormFactor(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantFormFactor int
		wantOk      bool
	}{
		{
			name:        "2280 form factor",
			input:       "Samsung 970 EVO M.2 2280",
			wantFormFactor: 2280,
			wantOk:      true,
		},
		{
			name:        "2230 form factor",
			input:       "WD SN530 2230 NVMe",
			wantFormFactor: 2230,
			wantOk:      true,
		},
		{
			name:        "2242 form factor",
			input:       "Crucial MX500 2242",
			wantFormFactor: 2242,
			wantOk:      true,
		},
		{
			name:        "22110 form factor",
			input:       "Intel Pro 22110 SSD",
			wantFormFactor: 22110,
			wantOk:      true,
		},
		{
			name:        "no form factor",
			input:       "Samsung 970 EVO",
			wantFormFactor: 0,
			wantOk:      false,
		},
		{
			name:        "partial match should not work",
			input:       "Samsung 22 inch monitor",
			wantFormFactor: 0,
			wantOk:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formFactor, ok := parseFormFactor(tt.input)
			if ok != tt.wantOk {
				t.Errorf("parseFormFactor(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
			}
			if formFactor != tt.wantFormFactor {
				t.Errorf("parseFormFactor(%q) formFactor = %d, want %d", tt.input, formFactor, tt.wantFormFactor)
			}
		})
	}
}

func TestSSDToMarkdown(t *testing.T) {
	ssd := SSD{
		DriveID:      "1461",
		URL:          "https://www.techpowerup.com/ssd-specs/corsair-mp600-mini-1-tb.d1461",
		Manufacturer: "Corsair",
		Name:         "MP600 Mini",
		Capacity:     "1 TB",
		FormFactor:   "M.2 2280",
		Interface:    "PCIe 4.0 x4",
		Protocol:     "NVMe 1.4",
		Dram:         "N/A",
		Hmb:          "64 MB",
		Endurance:    "600 TBW",
		SeqRead:      "4,800 MB/s",
		SeqWrite:     "4,800 MB/s",
		Controller: Controller{
			Manufacturer: "Phison",
			Name:         "PS5021-E21T",
		},
		Flash: Flash{
			Manufacturer: "Micron",
			Type:         "TLC",
		},
	}

	markdown := ssd.ToMarkdown()

	// Verify key content is present
	expectedStrings := []string{
		"Corsair",
		"MP600 Mini",
		"1 TB",
		"TLC",
		"PCIe 4.0 x4",
		"M.2 2280",
		"Phison PS5021-E21T",
		"N/A",
		"64 MB",
		"Micron",
		"4,800 MB/s",
		"600 TBW",
		"camelcamelcamel",
		TechPowerUpURL,
		GitHubURL,
		GitHubIssuesURL,
	}

	for _, expected := range expectedStrings {
		if !contains(markdown, expected) {
			t.Errorf("ToMarkdown() missing expected string %q", expected)
		}
	}
}

func TestGetHMBSize(t *testing.T) {
	tests := []struct {
		name string
		hmb  string
		want string
	}{
		{
			name: "known HMB size",
			hmb:  "64 MB",
			want: "64 MB",
		},
		{
			name: "unknown HMB",
			hmb:  "Unknown",
			want: "N/A",
		},
		{
			name: "empty HMB",
			hmb:  "",
			want: "",
		},
		{
			name: "N/A HMB",
			hmb:  "N/A",
			want: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssd := SSD{Hmb: tt.hmb}
			if got := ssd.GetHMBSize(); got != tt.want {
				t.Errorf("GetHMBSize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetDramSize(t *testing.T) {
	tests := []struct {
		name string
		dram string
		want string
	}{
		{
			name: "known DRAM size",
			dram: "1 GB",
			want: "1 GB",
		},
		{
			name: "unknown DRAM",
			dram: "Unknown",
			want: "N/A",
		},
		{
			name: "empty DRAM",
			dram: "",
			want: "",
		},
		{
			name: "N/A DRAM",
			dram: "N/A",
			want: "N/A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ssd := SSD{Dram: tt.dram}
			if got := ssd.GetDramSize(); got != tt.want {
				t.Errorf("GetDramSize() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTpuRepositoryInsertNotImplemented(t *testing.T) {
	tpu := NewTpuRepository("http://test", "user", "key")
	err := tpu.Insert(nil, SSD{})
	if err == nil {
		t.Error("Insert() expected error, got nil")
	}
	if err.Error() != "TpuRepository.Insert is not implemented" {
		t.Errorf("Insert() unexpected error message: %v", err)
	}
}

func TestTpuRepositoryUpdateNotImplemented(t *testing.T) {
	tpu := NewTpuRepository("http://test", "user", "key")
	err := tpu.Update(nil, SSD{})
	if err == nil {
		t.Error("Update() expected error, got nil")
	}
	if err.Error() != "TpuRepository.Update is not implemented" {
		t.Errorf("Update() unexpected error message: %v", err)
	}
}

func TestEsRepositoryUpdateNotImplemented(t *testing.T) {
	esRepo := &EsRepository{}
	err := esRepo.Update(nil, SSD{})
	if err == nil {
		t.Error("Update() expected error, got nil")
	}
	if err.Error() != "EsRepository.Update is not implemented" {
		t.Errorf("Update() unexpected error message: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
