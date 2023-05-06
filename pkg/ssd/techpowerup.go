package ssd

type TpuSSDRepository struct {
	apikey string
}

func NewTpuSSDRepository(apiKey string) *TpuSSDRepository {
	return &TpuSSDRepository{
		apikey: apiKey,
	}
}

func (tpu *TpuSSDRepository) GetById(id string) (SSD, error) {
	//TODO implement this
	return SSD{}, nil
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
