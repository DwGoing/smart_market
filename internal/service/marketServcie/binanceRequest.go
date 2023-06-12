package marketServcie

type Request struct {
	Id     string `json:"id"`
	Method string `json:"method"`
}

type Response struct {
	Id         string      `json:"id"`
	Status     int         `json:"status"`
	Error      Error       `json:"error"`
	RateLimits []RateLimit `json:"rateLimits"`
}

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type RateLimit struct {
	RateLimitType string `json:"rateLimitType"`
	Interval      string `json:"interval"`
	IntervalNum   int    `json:"intervalNum"`
	Limit         int    `json:"limit"`
	Count         int    `json:"count"`
}
