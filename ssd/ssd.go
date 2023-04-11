package main

import (
	"fmt"
	"strings"
)

type SSD struct {
	Brand    string
	Model    string
	NandType string
	Category string

	Interface string

	FormFactor string

	Controller string

	Configuration string

	DRAM string //yes/no

	HMB string //yes/no

	NandBrand string

	NandDimension string //2D or 3D

	Layers string

	ReadWrite string

	CellRow int
}

func (ssd SSD) toMarkdown() string {
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
		fmt.Sprintf("* 2D/3D NAND: **%s**", ssd.NandDimension),
		fmt.Sprintf("* Layers: **%s**", ssd.Layers),
		fmt.Sprintf("* R/W: **%s**", ssd.ReadWrite),
		fmt.Sprintf("[Click here to view this SSD in the tier list](https://docs.google.com/spreadsheets/d/1B27_j9NDPU3cNlj2HKcrfpJKHkOf-Oi1DbuuQva2gT4/edit#gid=0&amp;range=A%v:V%v)", ssd.CellRow, ssd.CellRow),
		fmt.Sprintf("[Click here to view camelcamelcamel product search page](https://camelcamelcamel.com/search?sq=%s).", ssd.Brand+" "+ssd.Model),
		fmt.Sprintf("---\n^(Suggestions, concerns, errors? Message us directly or submit an issue on), [^(Github!)](https://github.com/aattwwss/ssd-bot-go)"),
	}
	return strings.Join(arr, "\n\n")
}

func main() {

	ssd := SSD{}
	fmt.Println(ssd.toMarkdown())
}
