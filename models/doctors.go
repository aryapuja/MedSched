package models

import "time"

type Doctor struct {
	ID             int       `json:"id"`             
	Name           string    `json:"name"`           
	Specialization string    `json:"specialization"` 
	ContactInfo    string    `json:"contact_info"`   
	CreatedAt      time.Time `json:"created_at"`     
}
