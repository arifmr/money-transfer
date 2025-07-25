package models

type Response struct {
	StatusCode     int         `json:"status_code"`
	Data           interface{} `json:"data,omitempty"`
	Error          *ErrorLog   `json:"error,omitempty"`
	Total          int64       `json:"total,omitempty"`
	Message        string      `json:"message,omitempty"`
	RequestPayload interface{} `json:"request_payload,omitempty"`
}

type ErrorLog struct {
	Line              string      `json:"line,omitempty" bson:"line"`
	Filename          string      `json:"filename,omitempty" bson:"filename"`
	Function          string      `json:"function,omitempty" bson:"function"`
	Message           interface{} `json:"message,omitempty" bson:"message"`
	SystemMessage     interface{} `json:"system_message,omitempty" bson:"system_message"`
	Url               string      `json:"url,omitempty" bson:"url"`
	Method            string      `json:"method,omitempty" bson:"method"`
	Fields            interface{} `json:"fields,omitempty" bson:"fields"`
	ConsumerTopic     string      `json:"consumer_topic,omitempty" bson:"consumer_topic"`
	ConsumerPartition int         `json:"consumer_partition,omitempty" bson:"consumer_partition"`
	ConsumerName      string      `json:"consumer_name,omitempty" bson:"consumer_name"`
	ConsumerOffset    int64       `json:"consumer_offset,omitempty" bson:"consumer_offset"`
	ConsumerKey       string      `json:"consumer_key,omitempty" bson:"consumer_key"`
	Err               error       `json:"-"`
	StatusCode        int         `json:"-"`
}
