package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/yoonhero/ohpotatocoin/utils"
)

const (
	signature     string = "62f931acce9bf5179760e85cc7caf85b1ff085816797ddbd210d86f175e14bd23d07fa6a77914e8d65113b036b962b2de742eb1028136d86cbbbfbe96e69a38a"
	privateKey    string = "307702010104208da51576e4d2078c580b52cf50ee1fb0141a5f0f00e2f022f3efd15e15a46b92a00a06082a8648ce3d030107a14403420004aa12814b0ad1bbb3876e8ac29abb9d0227cd68f98e36d380e5b563b085ed48a3358aa20c98af29ee9e70bba7aa1a13bc5e354d0635ae34351b3e0a793dc02c88"
	hashedMessage string = "1c5863cd55b5a4413fd59f054af57ba3c75c0698b3851d70f99b8de2d5c7338f"
)

func Start() {
	privBytes, err := hex.DecodeString(privateKey)

	utils.HandleErr(err)

	_, err = x509.ParseECPrivateKey(privBytes)
	utils.HandleErr(err)

	sigBytes, err := hex.DecodeString(signature)
	rBytes := sigBytes[:len(sigBytes)/2]
	sBytes := sigBytes[len(sigBytes)/2:]

	var bigR, bigS = big.Int{}, big.Int{}

	bigR.SetBytes(rBytes)
	bigS.SetBytes(sBytes)

	fmt.Println(bigR, bigS)
}
