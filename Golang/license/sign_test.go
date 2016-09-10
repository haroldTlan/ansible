package main

import (
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"fmt"
	"io"
	"math/big"
	"testing"
)

//<config>
//	<server ipaddr="127.0.0.1" port="8080"/>
//  <license>
//    abc
//	</license>
//</config>
func fromBase10(base10 string) *big.Int {
	i, ok := new(big.Int).SetString(base10, 10)
	if !ok {
		panic("bad number: " + base10)
	}
	return i
}

func TestSign(t *testing.T) {
	sn := "c62965591712ef485341872c6295ebfffeaa8a3dce74d0e20dffe8b9e3ee0f65"
	loadPrivKey()
	/*
		signature, err := sign(sn)
		if err != nil {
			t.Fatalf("sign error: %v", err)
		}
	*/

	signature := "7da1b6ed81f9d67a7036486e4d33ecd664495ec7fb1b2a9cd7f740aedb2c6ccb36c409352db16f9edcca8953721315632132dbe1853ac3ca2cad6e870fca99ea1b2ae65c9914511474974ba718e8e53c74f5f61f0036a6d82aa5a2b8b9a4c474eba53a8f6811ecacc2c436c457a42d1a8f093f84e80df12f2d1574101b1b3f34"
	sigbuf := make([]byte, len(signature))
	fmt.Sscanf(signature, "%x", &sigbuf)

	hash := md5.New()
	io.WriteString(hash, sn)
	hashed := hash.Sum(nil)

	var h crypto.Hash
	pubKey := &rsa.PublicKey{
		N: fromBase10("126038038516492034489881010707522756455005310820723628794048567491219653586876002712941473403005276243429681350407059668213363248724006391092540187693872519570891047411229657493659432418029829008660673664620025809544514419347167680091518538641680141780633312725341167771832755283446081256635145120586638842379"),
		E: 65537}

	err := rsa.VerifyPKCS1v15(pubKey, h, hashed, sigbuf)
	if err != nil {
		t.Fatalf("verify error: %v", err)
	}
}

func TestByteCompare(t *testing.T) {
	b := make([]byte, 1024)
	fmt.Sscanf("616263", "%x", &b)
	fmt.Printf("abc: %v\n", string(b))
}
