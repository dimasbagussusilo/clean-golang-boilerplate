package model

type Response struct {
	StatusCode    int         `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data,omitempty"`
	Error         *ErrorLog   `json:"error,omitempty"`
	Page          int         `json:"page,omitempty"`
	PerPage       int         `json:"per_page,omitempty"`
	Total         int64       `json:"total,omitempty"`
	StatCode      string      `json:"stat_code,omitempty"`
	StatMsg       string      `json:"stat_msg,omitempty"`
	FileUrl       string      `json:"file_url,omitempty"`
	MaxPage       int         `json:"max_page,omitempty"`
}

type ResponseChannel struct {
	Data  interface{} `json:"data,omitempty"`
	Error error       `json:"errors,omitempty"`
}

type ErrorLog struct {
	Line              string      `json:"line,omitempty"`
	Filename          string      `json:"filename,omitempty"`
	Function          string      `json:"function,omitempty"`
	Message           interface{} `json:"message,omitempty"`
	SystemMessage     interface{} `json:"system_message,omitempty"`
	Url               string      `json:"url,omitempty"`
	Method            string      `json:"method,omitempty"`
	Fields            interface{} `json:"fields,omitempty"`
	ConsumerTopic     string      `json:"consumer_topic,omitempty"`
	ConsumerPartition int         `json:"consumer_partition,omitempty"`
	ConsumerName      string      `json:"consumer_name,omitempty"`
	ConsumerOffset    int64       `json:"consumer_offset,omitempty"`
	ConsumerKey       string      `json:"consumer_key,omitempty"`
	Err               error       `json:"-"`
	StatusCode        int         `json:"-"`
}

type MustActiveRequest struct {
	Table         string
	ReqField      string
	Clause        string
	CustomMessage string
	Id            interface{}
}
