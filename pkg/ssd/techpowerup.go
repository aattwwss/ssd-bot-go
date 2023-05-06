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
