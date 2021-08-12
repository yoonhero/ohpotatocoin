package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/yoonhero/ohpotatocoin/utils"
)

const (
	fileName string = "ohpotatocoin.wallet"
)

// wallet struct
type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

// // return file exists or not
// func hasWalletFile() bool {
// 	_, err := os.Stat(fileName)
// 	return !os.IsNotExist(err)

// }

// create random private key
func CreatePrivKey() *ecdsa.PrivateKey {
	// https://m.blog.naver.com/aepkoreanet/221178375642
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return priKey
}

// // save the key
// func persistKey(key *ecdsa.PrivateKey) {
// 	bytes, err := x509.MarshalECPrivateKey(key)
// 	utils.HandleErr(err)
// 	err = os.WriteFile(fileName, bytes, 0644)
// 	utils.HandleErr(err)
// }

// parse the key
// func restoreKey() (key *ecdsa.PrivateKey) {
// 	keyAsBytes, err := os.ReadFile(fileName)
// 	utils.HandleErr(err)
// 	key, err = x509.ParseECPrivateKey(keyAsBytes)
// 	utils.HandleErr(err)
// 	return
// }

// bytes to hex decimal
func encodeBigInts(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

// address from key
func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
}

func WalletSign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

// sign the signature
func Sign(payload string, keyAsBytes []byte) string {
	privkey, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, privkey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

// restore ints to byte
func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	// decode payload
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	utils.HandleErr(err)

	// divide bytes
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]

	bigA, bigB := big.Int{}, big.Int{}

	// set bytes
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)

	return &bigA, &bigB, nil
}

// Verify the signature
func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err)

	// make same publickey
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	// decode the payload
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	// Verify public key and payload and address
	ok := ecdsa.Verify(&publicKey, payloadBytes, r, s)
	return ok
}

// func Wallet() *wallet {
// 	if w == nil {
// 		w = &wallet{}
// 		// has a wallet already
// 		if hasWalletFile() {
// 			// yes -> restore from file
// 			w.privateKey = restoreKey()
// 		} else {
// 			// no -> create prv key, save to file
// 			key := CreatePrivKey()
// 			persistKey(key)
// 			w.privateKey = key
// 		}
// 		w.Address = aFromK(w.privateKey)
// 	}
// 	return w
// }

func RestApiWallet(key []byte) *wallet {
	var wall *wallet
	wall = &wallet{}
	restoredKey := restapiRestoreKey(key)
	wall.privateKey = restoredKey
	wall.Address = aFromK(wall.privateKey)
	return wall
}

// parse the key
func restapiRestoreKey(keyAsBytes []byte) *ecdsa.PrivateKey {
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return key
}

func RestApiCreatePrivKey() (address string, bytes []byte) {
	// https://m.blog.naver.com/aepkoreanet/221178375642
	priKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	bytes, err = x509.MarshalECPrivateKey(priKey)
	utils.HandleErr(err)
	address = aFromK(priKey)
	return address, bytes
}
