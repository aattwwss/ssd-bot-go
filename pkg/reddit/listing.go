package reddit

// Listing represents a Reddit API listing response.
type Listing[T any] struct {
	Kind string         `json:"kind"`
	Data ListingData[T] `json:"data"`
}

// ListingData contains the data within a listing response.
type ListingData[T any] struct {
	After     string        `json:"after"`
	Dist      int           `json:"dist"`
	Modhash   string        `json:"modhash"`
	GeoFilter string        `json:"geo_filter"`
	Children  []Children[T] `json:"children"`
	Before    string        `json:"before"`
}

// Children represents individual items within a listing.
type Children[T any] struct {
	Kind string `json:"kind"`
	Data T      `json:"data"`
}
