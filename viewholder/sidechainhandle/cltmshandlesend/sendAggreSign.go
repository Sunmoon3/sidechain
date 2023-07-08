package cltmshandlesend

import (
	"bytes"
	"crypto/md5"
	"crypto/schnorr/musig"
	"encoding/pem"
	"fmt"
	"viewholder/config"
	"viewholder/sidechainhandle/events"

	inithandle "viewholder/sidechainhandle"

	"go.dedis.ch/kyber/v3/suites"
	"io/ioutil"
	"strings"
	"time"
	pb "viewholder/sidechainhandle/listener/cross/proto"
)

func MD5Byte(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func CrossOutAggre(message *events.CrossTransferEvent, eventch chan<- interface{}, errch chan<- error) {
	log.Info("Aggre... ")
	payload := message.GetPayload()
	//log.Infof("payload:", payload)

	splitPayLoad := strings.Split(payload, ":")

	if len(splitPayLoad) != 2 {
		errch <- fmt.Errorf("error to do split on %v\n", payload)
		return
	}

	objectByte := []byte(splitPayLoad[1])

	log.Infof("starting generating Aggre proof....")
	fmt.Println("objectByte:", MD5Byte(objectByte))
	start := time.Now()
	// signByte, pubKeysByte := aggreByPrivateKey(objectByte, config.AggreSignCount)
	signByte, pubKeysByte := aggreByPrivateKey([]byte(MD5Byte(objectByte)), config.AggreSignCount)

	duration := time.Now().Sub(start)
	fmt.Printf("aggre computation finished, time spent %s", duration.String())

	log.Infof("starting generating VDF proof....")

	vdfOut, err := crypto.VDFSig([]byte{1}, config.VDFD)

	if err != nil {
		log.Errorf("VDFSig err", err)
	}

	//verfiy := assist.VDFVer(objectByte, vdfOut, config.VDFD)
	//fmt.Println("VDF verify:", verfiy)

	txId := message.GettxId()
	scId := splitPayLoad[0]

	log.Infof("starting generating crosslanguage with proof....")
	aggreReq := generateCrossLanguageAggre(txId, scId, pubKeysByte, signByte, objectByte, vdfOut, errch)
	log.Infof("generate crosslanguage with proof success")

	clgev := events.NewCrosslanguageGeneratedEvent(aggreReq)

	eventch <- clgev

	return
}

func aggreByPrivateKey(objectByte []byte, count int) ([]byte, []byte) {

	keys := getCountPrivateKey(count)

	signature, pubKeys := DealPrivateKey(keys, objectByte)

	signByte, pubKeysByte := cryptoToBytes(signature, pubKeys)

	return signByte, pubKeysByte

	//crypto, points := bytesToCrypto(signByte, pubKeysByte, count)

	//VerifyMsg(objectByte, signature, pubKeys)
	//VerifyMsg(objectByte, crypto, points)
}

func bytesToCrypto(signByte []byte, pubKeysByte []byte, count int) (*musig.Signature, []kyber.Point) {
	signature, err := decodeSignature(signByte)
	if err != nil {
		fmt.Println("DecodeSignature err!")
	}
	p := make([]kyber.Point, 0)

	for i := 0; i < count; i++ {
		t := curve.Point()
		err = t.UnmarshalBinary(pubKeysByte[i*32 : (i+1)*32])
		p = append(p, t)
		if err != nil {
			fmt.Println("UnmarshalBinary err!")
		}
	}
	return signature, p
}

func decodeSignature(sig []byte) (*musig.Signature, error) {
	p := curve.Point()
	err := p.UnmarshalBinary(sig[:32])
	if err != nil {
		return nil, err
	}
	s := curve.Scalar().SetBytes(sig[32:])
	return &musig.Signature{
		R: p,
		S: s,
	}, nil
}

func cryptoToBytes(signature *musig.Signature, pubKeys []kyber.Point) ([]byte, []byte) {
	signatureByte := append(encodePoint(signature.R), encodeScalar(signature.S)...)
	pubKeysByte := make([]byte, 0)
	for i := 0; i < len(pubKeys); i++ {
		pubKeysByte = append(pubKeysByte, encodePoint(pubKeys[i])...)
	}
	return signatureByte, pubKeysByte
}

func VerifyMsg(msg []byte, sig *musig.Signature, pubKeys []kyber.Point) {
	ok := musig.VerifySignature(msg, sig, pubKeys...)
	fmt.Println(ok)
}

var curve = suites.MustFind("Ed25519")

func encodePoint(p kyber.Point) []byte {
	b, _ := p.MarshalBinary()
	return b
}

func encodeScalar(s kyber.Scalar) []byte {
	b, _ := s.MarshalBinary()
	return b
}

func DealPrivateKey(privateKeys [][]byte, objectByte []byte) (*musig.Signature, []kyber.Point) {

	var signers []Signer
	var pubKeys []kyber.Point
	var publicNonces [][]kyber.Point
	var Rvalues []kyber.Point
	var sigs []kyber.Scalar

	count := len(privateKeys)
	noncesNum := count
	msg := objectByte

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

	for j := 0; j < noncesNum; j++ {
		Rj := curve.Point().Null()
		for i := 0; i < len(publicNonces); i++ {
			Rj = curve.Point().Add(publicNonces[i][j], Rj)
		}
		Rvalues = append(Rvalues, Rj)
	}

	R := musig.ComputeR(msg, Rvalues, pubKeys...)

	for _, s := range signers {
		sigs = append(sigs, musig.SignMulti(msg, s.Key, s.Nonces, R, Rvalues, pubKeys...))
	}

	sig := &musig.Signature{
		R: R,
		S: musig.AggregateSignatures(sigs...),
	}
	return sig, pubKeys
}

type Signer struct {
	Key    *musig.Key
	Nonces []*musig.Key
}

func createSigner(noncesNum int, cert []byte) Signer {
	var nonces []*musig.Key
	for j := 0; j < noncesNum; j++ {
		nonce := musig.NewKey(cert)
		nonces = append(nonces, nonce)
	}
	return Signer{
		Key:    musig.NewKey(cert),
		Nonces: nonces,
	}
}

func getCountPrivateKey(count int) [][]byte {
	paths := Path()
	var certs [][]byte
	for i := 0; i < count; i++ {
		keyPem, err := ioutil.ReadFile(paths[i])
		if err != nil {
			fmt.Println("ioutil.ReadFile:", err)
		}
		cert := PrivateKeyFromPEM(keyPem)
		certs = append(certs, cert)
	}
	return certs
}

func PrivateKeyFromPEM(raw []byte) []byte {
	var err error
	if len(raw) <= 0 {
		fmt.Println("PEM is nil", err)
		return nil
	}
	block, _ := pem.Decode(raw)
	plain := block.Bytes
	return plain
}

func generateCrossLanguageAggre(txId string, scId string, publicKey []byte, sig []byte, objectByte []byte, vdfOut []byte, errCh chan<- error) *pb.CrossLanguageRequest {

	localChainId := inithandle.GetLocalChainId()

	chainSign, err := inithandle.GetChainSign(objectByte)
	if err != nil {
		log.Errorf("inithandle.GetChainSign err \n")
		errCh <- err
		return nil
	}

	chainPubKey, err := inithandle.GetChainPubkey()
	if err != nil {
		log.Errorf("inithandle.GetChainPubKey err \n")
		errCh <- err
		return nil
	}

	separator := byte(':')
	keys := bytesTo2DBytes(sig, separator)

	return &pb.CrossLanguageRequest{
		Txid:        txId,
		Chainid:     localChainId,
		Senderid:    localChainId,
		Recverid:    scId,
		Rw:          1,
		Sendercerts: chainPubKey,
		Sendersign:  chainSign,
		Apk:         publicKey,
		Messages:    objectByte,
		Sigmas:      keys,
		Vdfsign:     vdfOut,
	}
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
