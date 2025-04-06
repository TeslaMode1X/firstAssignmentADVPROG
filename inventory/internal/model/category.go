package model

import (
	"fmt"
	"time"
)

type Category struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (c *Category) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	return nil
}
