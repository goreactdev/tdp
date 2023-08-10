package tonconnect

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/abi"
	"github.com/tonkeeper/tongo/liteapi"

	"encoding/base64"
	"encoding/binary"
	"encoding/hex"

	log "github.com/sirupsen/logrus"
)

const (
	tonProofPrefix   = "ton-proof-item-v2/"
	tonConnectPrefix = "ton-connect"
)

func SignatureVerify(pubkey ed25519.PublicKey, message, signature []byte) bool {
	return ed25519.Verify(pubkey, message, signature)
}

func ConvertTonProofMessage(ctx context.Context, tp *TonProof) (*ParsedMessage, error) {
	log := log.WithContext(ctx).WithField("prefix", "ConverTonProofMessage")

	addr, err := tongo.ParseAccountID(tp.Address)
	if err != nil {
		return nil, err
	}

	var parsedMessage ParsedMessage

	sig, err := base64.StdEncoding.DecodeString(tp.Proof.Signature)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	parsedMessage.Workchain = addr.Workchain
	parsedMessage.Address = addr.Address[:]
	parsedMessage.Domain = tp.Proof.Domain
	parsedMessage.Timstamp = tp.Proof.Timestamp
	parsedMessage.Signature = sig
	parsedMessage.Payload = tp.Proof.Payload
	parsedMessage.StateInit = tp.Proof.StateInit
	
	return &parsedMessage, nil
}

func CreateMessage(ctx context.Context, message *ParsedMessage) ([]byte, error) {
	wc := make([]byte, 4)
	binary.BigEndian.PutUint32(wc, uint32(message.Workchain))

	ts := make([]byte, 8)
	binary.LittleEndian.PutUint64(ts, uint64(message.Timstamp))

	dl := make([]byte, 4)
	binary.LittleEndian.PutUint32(dl, message.Domain.LengthBytes)
	m := []byte(tonProofPrefix)
	m = append(m, wc...)
	m = append(m, message.Address...)
	m = append(m, dl...)
	m = append(m, []byte(message.Domain.Value)...)
	m = append(m, ts...)
	m = append(m, []byte(message.Payload)...)
	log.Info(string(m))
	messageHash := sha256.Sum256(m)
	fullMes := []byte{0xff, 0xff}
	fullMes = append(fullMes, []byte(tonConnectPrefix)...)
	fullMes = append(fullMes, messageHash[:]...)
	res := sha256.Sum256(fullMes)
	log.Info(hex.EncodeToString(res[:]))
	return res[:], nil
}

func GetWalletPubKey(ctx context.Context, address tongo.AccountID, net *liteapi.Client) (ed25519.PublicKey, error) {
	_, result, err := abi.GetPublicKey(ctx, net, address)
	if err != nil {
		return nil, err
	}
	if r, ok := result.(abi.GetPublicKeyResult); ok {
		i := big.Int(r.PublicKey)
		b := i.Bytes()
		if len(b) < 24 || len(b) > 32 { //govno kakoe-to
			return nil, fmt.Errorf("invalid publock key")
		}
		return append(make([]byte, 32-len(b)), b...), nil //make padding if first bytes are empty
	}
	return nil, fmt.Errorf("can't get publick key")
}

func CheckProof(ctx context.Context, address tongo.AccountID, net *liteapi.Client, tonProofReq *ParsedMessage, domain string) (bool, error) {

	log := log.WithContext(ctx).WithField("prefix", "CheckProof")
	pubKey, err := GetWalletPubKey(ctx, address, net)
	if err != nil {
		if tonProofReq.StateInit == "" {
			log.Errorf("get wallet address error: %v", err)
			return false, err
		}

		pubKey, err = ParseStateInit(tonProofReq.StateInit)
		if err != nil {
			log.Errorf("parse wallet state init error: %v", err)
			return false, err
		}
	}

	if time.Now().After(time.Unix(tonProofReq.Timstamp, 0).Add(time.Duration(1000) * time.Second)) {
		msgErr := "proof has been expired"
		log.Error(msgErr)
		return false, fmt.Errorf(msgErr)
	}
	

	if tonProofReq.Domain.Value != domain {
		msgErr := fmt.Sprintf("wrong domain: %v", tonProofReq.Domain)
		log.Error(msgErr)
		return false, fmt.Errorf(msgErr)
	}

	mes, err := CreateMessage(ctx, tonProofReq)
	if err != nil {
		log.Errorf("create message error: %v", err)
		return false, err
	}

	return SignatureVerify(pubKey, mes, tonProofReq.Signature), nil
}