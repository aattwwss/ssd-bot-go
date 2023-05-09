package ssd

type TpuSSDRepository struct {
	username string
	apikey   string
}

func NewTpuSSDRepository(username, apiKey string) *TpuSSDRepository {
	return &TpuSSDRepository{
		username: username,
		apikey:   apiKey,
	}
}

func (tpu *TpuSSDRepository) FindById(id string) (*SSD, error) {
	//TODO implement this
	return &SSD{}, nil
}

func (tpu *TpuSSDRepository) Search(s string) ([]SSD, error) {
	//TODO implement this
	return nil, nil
}

func (tpu *TpuSSDRepository) Insert(ssd SSD) error {
	//TODO implement this
	return nil
}

func (tpu *TpuSSDRepository) Update(ssd SSD) error {
	//TODO implement this
	return nil
}
