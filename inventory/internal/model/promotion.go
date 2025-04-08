package model

import (
	"fmt"
	"time"
)

type Promotion struct {
	ID                 string
	Name               string
	Description        string
	DiscountPercentage float64
	ApplicableProducts []string
	Products           *[]Product
	StartDate          time.Time
	EndDate            time.Time
	IsActive           bool
}

func (p *Promotion) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("promotion name cannot be empty")
	}
	if p.DiscountPercentage <= 0.0 {
		return fmt.Errorf("promotion discount must be positive: %v", p)
	}
	if p.ApplicableProducts == nil {
		return fmt.Errorf("you need to have at least one applicable product: %v", p)
	}
	return nil
}
