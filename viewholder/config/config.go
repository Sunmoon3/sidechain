package config

import "path/filepath"

var (
	NODE_CMD          = "nodeCmd"
	LISTENER_CHAIN    = "listener->chain"
	CHAIN_DISPATCHER  = "chain-dispatcher"
	CROSS_DISPATCHER  = "cross-dispatcher"
	EVENT_PROCESSOR   = "event-processor"
	CROSS_LISTENER    = "cross-listener"
	SIDE_CHAIN_SERVER = "side-chain-server"
	CLTMS_HANDLE_SEND = "cltms-handle-send"
	CLTMS_HANDLE_REC  = "cltms-handle-rec"
	TPS_CMD           = "tps-cmd"

	NODE_DEFAULT_IP   = "127.0.0.1"
	NODE_DEFAULT_PORT = "7890"

	// OPPOSITE_IP = "192.168.182.130"
	OPPOSITE_IP = "127.0.0.1"

	// SIDE_CHAIN_ADDRE = "192.168.182.128:7890"
	SIDE_CHAIN_ADDRE = "10.170.58.122:7890"

	SIDE_CHAIN_ID = "ci6m7mqfeksn8mj6lhc0"

	KEY_STORE_PATH = "storekeys"

	LOCAL_CHAIN_IDENTIDY_STORE_KEY = "cn.bluarry.localchain.identity.store"
	LOCAL_CHAIN_CERT_STORE_KEY     = "cn.bluarry.localchain.cert.store"
	LOCAL_CHAIN_INITED_FLAG        = "cn.bluarry.localchain.inited"

	USER_CHAINCODE_PATENT = "patent"
	USER_CHAINCODE_TPS    = "tps"

	CONFIG_PATH       = "/root/work/chainmaker/viewholder_mult/config.yaml"
	SDK_CONFIG_PATH   = "/root/work/chainmaker/viewholder_mult/sdk_config.yml"
	SDK_CONFIG_PATH_2 = "/root/work/chainmaker/viewholder_mult/sdk_config_2.yml"

	CLIENT_DEFAULT_PEM    = filepath.Join("certs", "client", "tls", "cert.pem")
	CLIENT_DEFAULT_KEY    = filepath.Join("certs", "client", "tls", "private_key")
	CLIENT_DEFAULT_CA_PEM = filepath.Join("certs", "client", "tls", "ca.pem")

	SERVER_DEFAULT_PEM    = filepath.Join("certs", "server", "tls", "server.crt")
	SERVER_DEFAULT_KEY    = filepath.Join("certs", "server", "tls", "server.key")
	SERVER_DEFAULT_CA_PEM = filepath.Join("certs", "server", "tls", "ca.crt")

	MultSignCount = 3

	AggreSignCount = 4

	VDFD = 200
)
