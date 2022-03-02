module github.com/stafiprotocol/chainbridge

go 1.15

require (
	github.com/ChainSafe/log15 v1.0.0
	github.com/JFJun/go-substrate-crypto v1.0.1
	github.com/cosmos/cosmos-sdk v0.45.1
	github.com/cosmos/ibc-go v1.2.6
	github.com/decred/base58 v1.0.3
	github.com/ethereum/go-ethereum v1.10.15
	github.com/gorilla/websocket v1.4.2
	github.com/itering/scale.go v1.0.47
	github.com/itering/substrate-api-rpc v0.3.5
	github.com/mr-tron/base58 v1.2.0
	github.com/near/borsh-go v0.3.0
	github.com/shopspring/decimal v1.2.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.8.1
	github.com/stafihub/rtoken-relay-core/common v0.0.0-20220221025731-3b187c6383df
	github.com/stafihub/stafihub v0.0.0-20220302002328-f830f7d061a2
	github.com/stafiprotocol/go-substrate-rpc-client v1.1.3
	github.com/stafiprotocol/solana-go-sdk v1.2.0
	github.com/stretchr/testify v1.7.0
	github.com/tendermint/tendermint v0.34.14
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20220126234351-aa10faf2a1f8
)

replace (
	github.com/99designs/keyring => github.com/cosmos/keyring v1.1.6
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/keybase/go-keychain => github.com/99designs/go-keychain v0.0.0-20191008050251-8e49817e8af4
	github.com/near/borsh-go => github.com/tpkeeper/borsh-go v0.3.1-0.20210918134133-8433ab32a70c
)
