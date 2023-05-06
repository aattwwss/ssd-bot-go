package ssd

import (
	"fmt"
	"strings"
)

type SSDRepository interface {
	GetById(id string) (SSD, error)
	Search(s string) ([]SSD, error)
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
		fmt.Sprintf("The %s %s is a *%s* **%s** SSD.", ssd.Brand, ssd.Model, ssd.NandType, ssd.Category),
		fmt.Sprintf("* Interface: **%s**", ssd.Interface),
		fmt.Sprintf("* Form Factor: **%s**", ssd.FormFactor),
		fmt.Sprintf("* Controller: **%s**", ssd.Controller),
		fmt.Sprintf("* Configuration: **%s**", ssd.Configuration),
		fmt.Sprintf("* DRAM: **%s**", ssd.DRAM),
		fmt.Sprintf("* HMB: **%s**", ssd.HMB),
		fmt.Sprintf("* NAND Brand: **%s**", ssd.NandBrand),
		fmt.Sprintf("* NAND Type: **%s**", ssd.NandType),
		// fmt.Sprintf("* 2D/3D NAND: **%s**", ssd.NandDimension),
		fmt.Sprintf("* Layers: **%s**", ssd.Layers),
		fmt.Sprintf("* R/W: **%s**", ssd.ReadWrite),
		fmt.Sprintf("* Price History: **[camelcamelcamel](https://camelcamelcamel.com/search?sq=%s)**", ssd.Brand+" "+ssd.Model),
		fmt.Sprintf("---\n%s", ref),
	}
	return strings.Join(arr, "\n\n")
}

func main() {

	ssd := SSD{}
	fmt.Println(ssd.ToMarkdown())
}
