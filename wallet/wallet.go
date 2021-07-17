package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"

	"github.com/yoonhero/ohpotatocoin/utils"
)

const (
	fileName string = "ohpotatocoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)

}

func createPrivKey() *ecdsa.PrivateKey {
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return priKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func Wallet() *wallet {
	// has a wallet already?
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			// yes -> restore from file
		} else {
			// no -> create prv key, save to file
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}

	}
	return w
}
