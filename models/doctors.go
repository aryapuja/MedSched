package models

import "time"

type Doctor struct {
	ID             int       `json:"id"`             // ID dokter
	Name           string    `json:"name"`           // Nama dokter
	Specialization string    `json:"specialization"` // Spesialisasi dokter
	ContactInfo    string    `json:"contact_info"`   // Kontak dokter
	CreatedAt      time.Time `json:"created_at"`     // Tanggal dibuat
}
