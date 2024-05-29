package service

type Task struct {
	ID        int64   `json:"id"`
	Arg1      float64 `json:"arg1"`
	Arg2      float64 `json:"arg2"`
	Operation string  `json:"operation"`
}
