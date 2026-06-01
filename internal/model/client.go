package model

type Client struct {
	Id       int64
	Language Language
}

func NewClient(id int64) *Client {
	return &Client{Id: id, Language: EnglishLanguage}
}
