package model

import (
	_ "encoding/json"
)

// A HealthzResponse expresses health check message.
type HealthzResponse struct{
	 Message string `json:"message"`
}
