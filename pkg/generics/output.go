package generics

type ItemsOutput[T any] struct {
	Success bool  `json:"success"`
	Total   int64 `json:"total"`
	Items   []T   `json:"items"`
	Error   error `json:"error,omitempty"`
}

type ItemOutput[T any] struct {
	Success bool  `json:"success"`
	Item    T     `json:"item"`
	Error   error `json:"error,omitempty"`
}
