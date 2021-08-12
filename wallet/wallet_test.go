package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"testing"
)

const (
	testKey     = "307702010104204e373c3aa9c1a578b6a8a8e6b5ac94a80c93f0212f3ff1f369aac2f67fb3b920a00a06082a8648ce3d030107a144034200049ae3e98182094236c5ac66a23c96313369a17a2c16a0f9c30953b889eee897ba370275552dd0e12521c126a095011bd8d55fbaf9223424d75676be4a583fe247"
	testPayload = "00dbf4792e022af73040829f1ba8a3618ce6002068ac0ce66a3b0953e4b37145"
	testSig     = "4ca00232f369f06d290219bf46acbee059990ea070683fb0ccd4a40c04d8752466256a66e8f8d01b2393327826349ad2a3e93866aa96fc04c4d248000d33a160"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w
}

// func TestVerify(t *testing.T) {
// 	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
// 	b, _ := x509.MarshalECPrivateKey(privKey)
// 	t.Logf("%x", b)
// }

func TestWalletSign(t *testing.T) {
	s := WalletSign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{testPayload, true},
		{"00dbf4792e022af73040829f1ba8a3618ce6002068ac0ce66a3b0953e4b37144", false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and testPayload")
		}
	}

}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("RestoreBigInts should return error when payload is not hex.")
	}
}
