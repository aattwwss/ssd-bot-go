package ssd

type EsSSDRepository struct {
	Host         string
	Port         string
	AccessKey    string
	AccessSecret string
}

func NewEsSSDRepository(host, port, accessKey, accessSecret string) *EsSSDRepository {
	return &EsSSDRepository{
		Host:         host,
		Port:         port,
		AccessKey:    accessKey,
		AccessSecret: accessSecret,
	}
}

func (tpu *EsSSDRepository) GetById(id string) (SSD, error) {
	//TODO implement this
	return SSD{}, nil
}

func (tpu *EsSSDRepository) Search(s string) ([]SSD, error) {
	//TODO implement this
	return nil, nil
}

func (tpu *EsSSDRepository) Insert(ssd SSD) error {
	//TODO implement this
	return nil
}

func (tpu *EsSSDRepository) Update(ssd SSD) error {
	//TODO implement this
	return nil
}
