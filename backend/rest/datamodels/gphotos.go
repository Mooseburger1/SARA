package datamodels

type Media struct {
	Id         string `json:"id"`
	ProductUrl string `json:"productUrl"`
	MimeType   string `json:"mimeType"`
}

type MediaItems struct {
	Items []Media `json:"mediaItems"`
}
