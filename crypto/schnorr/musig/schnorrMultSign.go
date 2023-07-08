package musig

import (
	"bytes"
	"fmt"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/suites"
	"reflect"
)

type Signer struct {
	Key    *Key
	Nonces []*Key
}

func createSigner(noncesNum int, cert []byte) Signer {

	var nonces []*Key
	for j := 0; j < noncesNum; j++ {
		nonce := NewKey(cert)
		nonces = append(nonces, nonce)
	}

	return Signer{
		Key:    NewKey([]byte("cert")),
		Nonces: nonces,
	}
}

func DealPrivateKey(privateKeys [][]byte) (Signature, []kyber.Point) {
	var curve = suites.MustFind("Ed25519")

	count := len(privateKeys)

	noncesNum := count

	msg := []byte("msg-test")

	var signers []Signer
	var pubKeys []kyber.Point
	var publicNonces [][]kyber.Point

	for i := 0; i < count; i++ {
		signer := createSigner(noncesNum, privateKeys[i])
		signers = append(signers, signer)
		pubKeys = append(pubKeys, signer.Key.Pub)
		var pubNonces []kyber.Point
		for _, nonce := range signer.Nonces {
			pubNonces = append(pubNonces, nonce.Pub)
		}
		publicNonces = append(publicNonces, pubNonces)
	}

	var Rvalues []kyber.Point
	for j := 0; j < noncesNum; j++ {
		Rj := curve.Point().Null()
		for i := 0; i < len(publicNonces); i++ {
			Rj = curve.Point().Add(publicNonces[i][j], Rj)
		}
		Rvalues = append(Rvalues, Rj)
	}

	R := ComputeR(msg, Rvalues, pubKeys...)

	var sigs []kyber.Scalar
	for _, s := range signers {
		sigs = append(sigs, SignMulti(msg, s.Key, s.Nonces, R, Rvalues, pubKeys...))
	}
	sig := Signature{
		R: R,
		S: AggregateSignatures(sigs...),
	}
	return sig, pubKeys
}

func SignMsg(msg []byte, sig *Signature, pubKeys []kyber.Point) {
	ok := VerifySignature(msg, sig, pubKeys...)
	fmt.Println(ok)
}

func Each(certs [][]byte) {
	sig, pubKeys := DealPrivateKey(certs)
	sigByte, pubKeyByte := toBytes(sig, pubKeys)
	fmt.Println("type:", reflect.TypeOf(sigByte))
	separator := byte(':')
	keys := bytesTo2DBytes(sigByte, separator)
	fmt.Println("type:", reflect.TypeOf(keys), "keys:", keys)
	crypto, points := toCrypto(keys[0], pubKeyByte)
	//SignMsg([]byte("msg-test"), &sig, pubKeys)
	SignMsg([]byte("msg-test"), crypto, points)
}

func bytesTo2DBytes(data []byte, sep byte) [][]byte {
	splitData := bytes.Split(data, []byte{sep})
	result := make([][]byte, len(splitData))
	for i, b := range splitData {
		result[i] = make([]byte, len(b))
		copy(result[i], b)
	}
	return result
}

func toCrypto(sigByte []byte, pubKeyByte []byte) (*Signature, []kyber.Point) {
	signature, err := DecodeSignature(sigByte)
	if err != nil {
		fmt.Println("DecodeSignature err!")
	}
	p := make([]kyber.Point, 0)

	for i := 0; i < 2; i++ {
		t := curve.Point()
		err = t.UnmarshalBinary(pubKeyByte[i*32 : (i+1)*32])
		p = append(p, t)
		if err != nil {
			fmt.Println("UnmarshalBinary err!")
		}
	}
	return signature, p
}

func toBytes(sig Signature, keys []kyber.Point) ([]byte, []byte) {
	pubKeysByte := make([]byte, 0)

	for i := 0; i < len(keys); i++ {
		pubKeysByte = append(pubKeysByte, encodePoint(keys[i])...)
	}

	signatureByte := append(encodePoint(sig.R), encodeScalar(sig.S)...)

	return signatureByte, pubKeysByte
}
