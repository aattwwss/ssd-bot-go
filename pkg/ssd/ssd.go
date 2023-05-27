package ssd

import (
	"context"
	"fmt"
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

// type SSD struct {
// 	Brand    string
// 	Model    string
// 	NandType string
// 	Category string
//
// 	Interface string
//
// 	FormFactor string
//
// 	Controller string
//
// 	Configuration string
//
// 	DRAM string //yes/no
//
// 	HMB string //yes/no
//
// 	NandBrand string
//
// 	NandDimension string //2D or 3D
//
// 	Layers string
//
// 	ReadWrite string
//
// 	CellRow int
//
// 	Capacity string //not used
// }

func (ssd SSD) ToMarkdown() string {

	ref := fmt.Sprintf(
		"[^(Data Sheet)](%s) ^| [^( Github)](%s) ^| [^(Issues)](%s)",
		"https://docs.google.com/spreadsheets/d/1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4/edit#gid=0",
		"https://github.com/aattwwss/ssd-bot-go",
		"https://github.com/aattwwss/ssd-bot-go/issues",
	)

	arr := []string{
		fmt.Sprintf("The %s %s is a *%s* SSD.", ssd.Manufacturer, ssd.Name, ssd.Flash.Type),
		fmt.Sprintf("* Interface: **%s**", ssd.Interface),
		fmt.Sprintf("* Form Factor: **%s**", ssd.FormFactor),
		fmt.Sprintf("* Controller: **%s %s**", ssd.Controller.Manufacturer, ssd.Controller.Name),
		fmt.Sprintf("* DRAM: **%s**", ssd.Dram),
		fmt.Sprintf("* HMB: **%s**", ssd.Hmb),
		fmt.Sprintf("* NAND Brand: **%s**", ssd.Flash.Manufacturer),
		fmt.Sprintf("* NAND Type: **%s**", ssd.Flash.Type),
		// fmt.Sprintf("* 2D/3D NAND: **%s**", ssd.NandDimension),
		fmt.Sprintf("* R/W: **%s - %s**", ssd.SeqRead, ssd.SeqWrite),
		fmt.Sprintf("* Price History: **[camelcamelcamel](https://camelcamelcamel.com/search?sq=%s)**", ssd.Manufacturer+" "+ssd.Name),
		fmt.Sprintf("* Detailed Link: **[TechPowerUp](%s)**", ssd.URL),
		fmt.Sprintf("---\n%s", ref),
	}
	return strings.Join(arr, "\n\n")
}

func main() {

	ssd := SSD{}
	fmt.Println(ssd.ToMarkdown())
}
