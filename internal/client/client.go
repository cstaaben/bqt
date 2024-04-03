/*
	bqt - BigQuery TUI
	Copyright (C) 2024  Corbin Staaben<cstaaben@gmail.com>

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package client wraps the BigQuery clients used by bqt.
package client

import (
	"context"
	"errors"
	"fmt"

	"cloud.google.com/go/bigquery"
)

// Client wraps a map of BigQuery clients, stored by project ID.
type Client struct {
	clients map[string]*bigquery.Client
}

// New creates a new Client.
func New() *Client {
	return &Client{
		clients: make(map[string]*bigquery.Client),
	}
}

// GetClient returns a BigQuery client for the given project ID.
func (c *Client) GetClient(ctx context.Context, projectID string) (*bigquery.Client, error) {
	if client, ok := c.clients[projectID]; ok {
		return client, nil
	}

	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("creating %s client: %w", projectID, err)
	}

	c.clients[projectID] = client

	return client, nil
}

// Close closes all BigQuery clients, returning any errors it encounters.
func (c *Client) Close() error {
	var err error
	for _, client := range c.clients {
		err = errors.Join(err, client.Close())
	}

	return err
}

// Query runs a query against the given project ID and returns a row iterator.
func (c *Client) Query(ctx context.Context, projectID, query string) (*bigquery.RowIterator, error) {
	client, err := c.GetClient(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("getting client: %w", err)
	}
	q := client.Query(query)
	return q.Read(ctx)
}
