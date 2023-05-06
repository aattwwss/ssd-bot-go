package techpowerup

import "fmt"

type TechPowerUpService interface {
	GetById(id string)
	Search(id string)
}

type TPUImpl struct {
	apiKey string
}

func (ti TPUImpl) GetById(id string) {
	fmt.Println("vim-go")
}

func (ti TPUImpl) Search(id string) {
	fmt.Println("vim-go")
}
