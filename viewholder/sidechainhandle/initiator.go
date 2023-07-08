package sidechainhandle

import (
	cltms "crypto"
	"encoding/hex"
	"fmt"
	"viewholder/chain"
	"viewholder/config"
	"viewholder/store"

	"github.com/rs/xid"
)

var (
	LBP = store.NewLevelDBProvider(config.KEY_STORE_PATH)
)

func InitCLTMS(configpath string) error {
	config, err := chain.LoadConfig(configpath)
	if err != nil {
		return err
	}
	nodes := len(config.Singers)

	storeService := store.NewStoreService(LBP)

	for i := 0; i < nodes; i++ {
		keys, err := cltms.KeyGen()
		if err != nil {
			return err
		}
		serilize, err := keys.Serilize()
		if err != nil {
			return err
		}

		storeService.StoreKeyValue([]byte(config.Singers[i].Domain), serilize)

	}
	return nil
}

func InitChainId(configPath, viewholder string) error {
	keyGen, err := cltms.KeyGen()
	if err != nil {
		return err
	}

	serilize, err := keyGen.Serilize()
	if err != nil {
		return err
	}

	encodeToString := hex.EncodeToString(serilize)

	storeService := store.NewStoreService(LBP)

	chainid := xid.New()
	storeService.StoreKeyValue([]byte(config.LOCAL_CHAIN_IDENTIDY_STORE_KEY), chainid.Bytes())
	storeService.StoreKeyValue([]byte(config.LOCAL_CHAIN_CERT_STORE_KEY), serilize)

	fmt.Println("s_chain_identity:", chainid.String())
	fmt.Println("s_chain_viewholder_ip_with_port:", viewholder)
	fmt.Println("s_chain_cert:", encodeToString)

	return nil
}

func InitFinished() bool {
	storeService := store.NewStoreService(LBP)
	return storeService.StoreKeyValue([]byte(config.LOCAL_CHAIN_INITED_FLAG), []byte("1"))
}

func InitedCheck() bool {
	storeService := store.NewStoreService(LBP)
	return storeService.KeyExists([]byte(config.LOCAL_CHAIN_INITED_FLAG))
}

func GetLocalChainId() string {
	storeService := store.NewStoreService(LBP)
	queryByKey := storeService.QueryByKey([]byte(config.LOCAL_CHAIN_IDENTIDY_STORE_KEY))
	if len(queryByKey) == 0 {
		return ""
	}

	fromBytes, err := xid.FromBytes(queryByKey)
	if err != nil {
		return ""
	}

	return fromBytes.String()
}

func GetChainSign(message []byte) ([]byte, error) {
	storeservice := store.NewStoreService(LBP)
	bt := storeservice.QueryByKey([]byte(config.LOCAL_CHAIN_CERT_STORE_KEY))
	//fmt.Println("bt=========================>", bt)
	key := cltms.Keys{}
	ok, err := key.DeSerilize(bt)
	if err != nil || !ok {
		return []byte{}, err
	}
	return key.Sign(message)
}

func GetChainPubkey() ([]byte, error) {
	storeservice := store.NewStoreService(LBP)
	bt := storeservice.QueryByKey([]byte(config.LOCAL_CHAIN_CERT_STORE_KEY))
	//fmt.Println("bt:", bt)
	key := cltms.Keys{}
	ok, err := key.DeSerilize(bt)
	if err != nil || !ok {
		return []byte{}, err
	}
	return key.PublicKey.Bytes(), nil
}
