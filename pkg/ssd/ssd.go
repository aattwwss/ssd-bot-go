package ssd

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// URLs for Reddit comment references
const (
	TechPowerUpURL    = "https://www.techpowerup.com/ssd-specs"
	TechPowerUpQueryURL = "https://www.techpowerup.com/ssd-specs/?q="
	GitHubURL         = "https://github.com/aattwwss/ssd-bot-go"
	GitHubIssuesURL   = "https://github.com/aattwwss/ssd-bot-go/issues"
	CamelCamelURL     = "https://camelcamelcamel.com/search?sq="
)

// Repository defines the interface for SSD data storage and retrieval.
type Repository interface {
	FindById(ctx context.Context, id string) (*SSD, error)
	Insert(ctx context.Context, ssd SSD) error
	Update(ctx context.Context, ssd SSD) error
	SearchBasic(ctx context.Context, s string) ([]SSDBasic, error)
	Search(ctx context.Context, s string) ([]SSD, error)
}

// SSD represents a solid-state drive with its specifications.
type SSD struct {
	DriveID      string     `json:"driveId"`
	URL          string     `json:"url"`
	Manufacturer string     `json:"mfgr"`
	Name         string     `json:"name"`
	Capacity     string     `json:"capacity"`
	FormFactor   string     `json:"formFactor"`
	Interface    string     `json:"interface"`
	Protocol     string     `json:"protocol"`
	Dram         string     `json:"dram"`
	Hmb          string     `json:"hmb"`
	Released     string     `json:"released"`
	Endurance    string     `json:"endurance"`
	Warranty     string     `json:"warranty"`
	SeqRead      string     `json:"seqRead"`
	SeqWrite     string     `json:"seqWrite"`
	Controller   Controller `json:"controller"`
	Flash        Flash      `json:"flash"`
}

// Controller represents the SSD controller information.
type Controller struct {
	Manufacturer string `json:"mfgr"`
	Name         string `json:"name"`
	NameShort    string `json:"nameShort"`
	Channels     string `json:"channels"`
}

// Flash represents the flash memory information.
type Flash struct {
	Manufacturer string `json:"mfgr"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Layers       string `json:"layers"`
}

// SSDBasic contains basic SSD information for search results.
type SSDBasic struct {
	DriveID      string `json:"driveId"`
	Manufacturer string `json:"mfgr"`
	Name         string `json:"name"`
	Capacity     string `json:"capacity"`
	FormFactor   string `json:"formFactor"`
}

// GetHMBSize returns the HMB (Host Memory Buffer) size, or "N/A" if unknown.
func (ssd SSD) GetHMBSize() string {
	if ssd.Hmb == "Unknown" {
		return "N/A"
	}
	return ssd.Hmb
}

// GetDramSize returns the DRAM cache size, or "N/A" if unknown.
func (ssd SSD) GetDramSize() string {
	if ssd.Dram == "Unknown" {
		return "N/A"
	}
	return ssd.Dram
}

// ToMarkdown converts SSD to Markdown format to support
// formatting in a reddit comment submission
func (ssd SSD) ToMarkdown() string {
	ref := fmt.Sprintf(
		"[^(TechPowerup Database)](%s) ^| [^( Github)](%s) ^| [^(Issues)](%s)",
		TechPowerUpURL,
		GitHubURL,
		GitHubIssuesURL,
	)

	arr := []string{
		fmt.Sprintf("The %s %s %s is a *%s* SSD.", ssd.Manufacturer, ssd.Name, ssd.Capacity, ssd.Flash.Type),
		fmt.Sprintf("* Interface: **%s**", ssd.Interface),
		fmt.Sprintf("* Form Factor: **%s**", ssd.FormFactor),
		fmt.Sprintf("* Controller: **%s %s**", ssd.Controller.Manufacturer, ssd.Controller.Name),
		fmt.Sprintf("* DRAM: **%s**", ssd.GetDramSize()),
		fmt.Sprintf("* HMB: **%s**", ssd.GetHMBSize()),
		fmt.Sprintf("* NAND Brand: **%s**", ssd.Flash.Manufacturer),
		fmt.Sprintf("* NAND Type: **%s**", ssd.Flash.Type),
		fmt.Sprintf("* R/W: **%s - %s**", ssd.SeqRead, ssd.SeqWrite),
		fmt.Sprintf("* Endurance: **%s**", ssd.Endurance),
		fmt.Sprintf("* Price History: **[camelcamelcamel](%s)**", CamelCamelURL+url.QueryEscape(ssd.Manufacturer+" "+ssd.Name+" "+ssd.Capacity)),
		fmt.Sprintf("* Detailed Link: **[TechPowerUp SSD Database](%s)**", ssd.URL),
		fmt.Sprintf("* Variations: **[TechPowerUp SSD](%s)**", TechPowerUpQueryURL+url.QueryEscape(ssd.Manufacturer+" "+ssd.Name)),
		fmt.Sprintf("---\n%s", ref),
	}
	return strings.Join(arr, "\n\n")
}
