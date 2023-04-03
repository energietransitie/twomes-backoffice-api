package twomes

// A Measurement is a measured value for a specific property.
type Measurement struct {
	ID       uint     `json:"id"`
	UploadID uint     `json:"upload_id"`
	Property Property `json:"property"`
	Time     Time     `json:"time"`
	Value    string   `json:"value"`
}
