module github.com/stafiprotocol/chainbridge

go 1.15

require (
	github.com/ChainSafe/log15 v1.0.0
	github.com/decred/base58 v1.0.3
	github.com/ethereum/go-ethereum v1.10.6
	github.com/go-kit/kit v0.9.0 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/itering/scale.go v1.0.47
	github.com/itering/substrate-api-rpc v0.3.5
	github.com/mr-tron/base58 v1.2.0
	github.com/near/borsh-go v0.3.0
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stafiprotocol/go-substrate-rpc-client v1.0.9
	github.com/stafiprotocol/solana-go-sdk v0.3.10
	github.com/stretchr/testify v1.7.0
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

replace github.com/near/borsh-go => github.com/tpkeeper/borsh-go v0.3.1-0.20210918134133-8433ab32a70c
