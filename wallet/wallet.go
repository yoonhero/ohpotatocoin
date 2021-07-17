package wallet

import (
	"crypto/ecdsa"
	"os"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

func hasWalletFile() bool {
	_, err := os.Stat("nomadcoin.wallet")
	return !os.IsNotExist(err)

}

var w *wallet

func Wallet() *wallet {
	// has a wallet already?
	if w == nil {
		if hasWalletFile() {
			// yes -> restore from file

		}
		// no -> create prv key, save to file

	}
	return w
}
