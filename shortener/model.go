package shortener

type ModelShorten struct {
	Id    string `json:"-" bson:"_id"`
	Url   string `json:"url" bson:"url"`
	Count int64  `json:"-" bson:"count"`
}

type Operation struct {
	Result string `json:"result"`
}
