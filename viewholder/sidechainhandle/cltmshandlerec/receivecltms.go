package sidechain

import (
	"crypto/md5"
	"crypto/schnorr/musig"
	"fmt"
	cltms "git.labpc.bluarry.top/bluarry/CLTMS/crypto"
	"git.labpc.bluarry.top/bluarry/viewholder/sidechain/events"
	"viewholder/chain"

	"go.dedis.ch/kyber/v3/suites"
	"time"
)

func CrossInAggr(msg *events.CrossLanguageRecvedEvent, errch chan<- error) {
	logger.Infof("Starting Crossin Aggr")
	configpath := viper.GetString("configpath")
	config, err := chain.LoadConfig(configpath)
	if err != nil {
		errch <- err
		return
	}

	chainid := GetLocalChainid()
	if chainid != msg.Req.Recverid {
		errch <- fmt.Errorf("get the wrong cl")
		return
	}
	logger.Infof("get localchainid : %v", chainid)

	logger.Infof("Start verify the chain's sign ing...")

	ok, err := verifyChainSign(msg.Req.Sendercerts, msg.Req.Messages, msg.Req.Sendersign)
	if !ok || err != nil {
		errch <- fmt.Errorf("verify not pass or verify key error: %v", err)
		return
	}
	logger.Infof("verify the chain's sign successed")
	logger.Infof("Start verify the chain's CLMTS ing...")

	err := cltmsAggrVerify(msg.Req.Messages, msg.Req.Apk, msg.Req.Vdfsign, msg.Req.Sigmas)
	if err != nil {
		errch <- fmt.Errorf("verfy not pass while verify sigma: %v", err)
		return
	}
	logger.Infof("verify the chain's CLMTS successed")

	ok, err := chaincode.UserChaincode_CrossIn(configpath, string(msg.Req.Messages), msg.Req.Senderid)
	if !ok || err != nil {
		errch <- fmt.Errorf("transfer the state into chain error: %v", err)
		return
	}
	println("end:", time.Now().UnixNano())
	return
}

func MD5Byte(data []byte) string {
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func cltmsAggrVerify(messages []byte, apk []byte, vdfsign []byte, sigmas [][]byte) error {

	cltmS := cltms.NewCLTMS(0, 200)
	ok := cltmS.VDFVer([]byte{1}, vdfsign)
	if !ok {
		return fmt.Errorf("VDF Verify error")
	}

	signByte := sigmas[0]
	crypto, points := bytesToCrypto(signByte, apk, 4)
	//fmt.Println("messages:", MD5Byte(messages), "crypto")
	ok = musig.VerifySignature([]byte(MD5Byte(messages)), crypto, points...)
	if !ok {
		return fmt.Errorf("Aggre Verify error")
	}
	return nil
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

var curve = suites.MustFind("Ed25519")
