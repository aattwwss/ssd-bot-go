package ssd

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

type SSDRepository interface {
	FindById(ctx context.Context, id string) (*SSD, error)
	Insert(ctx context.Context, ssd SSD) error
	Update(ctx context.Context, ssd SSD) error
	SearchBasic(ctx context.Context, s string) ([]BasicSSD, error)
	Search(ctx context.Context, s string) ([]SSD, error)
}

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

type Controller struct {
	Manufacturer string `json:"mfgr"`
	Name         string `json:"name"`
	NameShort    string `json:"nameShort"`
	Channels     string `json:"channels"`
}

type Flash struct {
	Manufacturer string `json:"mfgr"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Layers       string `json:"layers"`
}

type BasicSSD struct {
	DriveID      string `json:"driveId"`
	Manufacturer string `json:"mfgr"`
	Name         string `json:"name"`
	Capacity     string `json:"capacity"`
	FormFactor   string `json:"formFactor"`
}

func (ssd SSD) getHMB() string {
	if ssd.Hmb == "Unknown" {
		return "N/A"
	}
	return ssd.Hmb
}

func (ssd SSD) getDram() string {
	if ssd.Dram == "Unknown" {
		return "N/A"
	}
	return ssd.Dram
}

func (ssd SSD) ToMarkdown() string {

	ref := fmt.Sprintf(
		"[^(TechPowerup Database)](%s) ^| [^( Github)](%s) ^| [^(Issues)](%s)",
		"https://www.techpowerup.com/ssd-specs",
		"https://github.com/aattwwss/ssd-bot-go",
		"https://github.com/aattwwss/ssd-bot-go/issues",
	)

	arr := []string{
		fmt.Sprintf("The %s %s is a *%s* SSD.", ssd.Manufacturer, ssd.Name, ssd.Flash.Type),
		fmt.Sprintf("* Interface: **%s**", ssd.Interface),
		fmt.Sprintf("* Form Factor: **%s**", ssd.FormFactor),
		fmt.Sprintf("* Controller: **%s %s**", ssd.Controller.Manufacturer, ssd.Controller.Name),
		fmt.Sprintf("* DRAM: **%s**", ssd.getDram()),
		fmt.Sprintf("* HMB: **%s**", ssd.getHMB()),
		fmt.Sprintf("* NAND Brand: **%s**", ssd.Flash.Manufacturer),
		fmt.Sprintf("* NAND Type: **%s**", ssd.Flash.Type),
		fmt.Sprintf("* R/W: **%s - %s**", ssd.SeqRead, ssd.SeqWrite),
		fmt.Sprintf("* Endurance: **%s**", ssd.Endurance),
		fmt.Sprintf("* Price History: **[camelcamelcamel](https://camelcamelcamel.com/search?sq=%s)**", url.QueryEscape(ssd.Manufacturer+" "+ssd.Name)),
		fmt.Sprintf("* Detailed Link: **[TechPowerUp](https://www.techpowerup.com/ssd-specs/#%s)**", url.QueryEscape(filterName(ssd.Name))),
		fmt.Sprintf("---\n%s", ref),
	}
	return strings.Join(arr, "\n\n")
}

func filterName(s string) string {
	return strings.ReplaceAll(s, "/", "")
}
