package client

import "cloud.google.com/go/bigquery"

type Client struct {
	client *bigquery.Client
}

func New(client *bigquery.Client) *Client {
	return &Client{
		client: client,
	}
}
