package stocks

type Response struct {
	Messages   []string  `json:"messages"`
	Status     string    `json:"status"`
}
