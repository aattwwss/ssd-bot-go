package reddit

type Listing[T any] struct {
	Kind string         `json:"kind"`
	Data ListingData[T] `json:"data"`
}

type ListingData[T any] struct {
	After     string        `json:"after"`
	Dist      int           `json:"dist"`
	Modhash   string        `json:"modhash"`
	GeoFilter string        `json:"geo_filter"`
	Children  []Children[T] `json:"children"`
	Before    string        `json:"before"`
}

type Children[T any] struct {
	Kind string `json:"kind"`
	Data T      `json:"data"`
}
