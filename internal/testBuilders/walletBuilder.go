package testBuilders

import (
	"github.com/Davmie/javaCode/models"
)

type WalletBuilder struct {
	wallet models.Wallet
}

func NewWalletBuilder() *WalletBuilder {
	return &WalletBuilder{}
}

func (b *WalletBuilder) WithID(id int) *WalletBuilder {
	b.wallet.ID = id
	return b
}

func (b *WalletBuilder) WithUID(uid string) *WalletBuilder {
	b.wallet.UID = uid
	return b
}

func (b *WalletBuilder) WithAmount(amount int) *WalletBuilder {
	b.wallet.Amount = amount
	return b
}

func (b *WalletBuilder) Build() models.Wallet {
	return b.wallet
}
