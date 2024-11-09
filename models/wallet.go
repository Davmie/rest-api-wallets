package models

type Tabler interface {
	TableName() string
}

func (Wallet) TableName() string {
	return "wallet"
}

type Wallet struct {
	ID     int    `json:"id" db:"id"`
	UID    string `json:"uid" db:"uid"` // a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11
	Amount int    `json:"amount" db:"amount"`
}
