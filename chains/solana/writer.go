package solana

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/chains"
	"github.com/stafiprotocol/chainbridge/utils/msg"
	solClient "github.com/stafiprotocol/solana-go-sdk/client"
	"github.com/stafiprotocol/solana-go-sdk/common"
	"github.com/stafiprotocol/solana-go-sdk/tokenprog"
)

var msgLimit = 4096

// write to solana
type writer struct {
	conn            *Connection
	minterProgramId common.PublicKey
	mintManager     common.PublicKey
	router          chains.Router
	log             log15.Logger
	sysErr          chan<- error
	msgChan         chan msg.Message
	stop            <-chan int
}

func NewWriter(conn *Connection, minterProgramId, mintManager common.PublicKey, log log15.Logger, stop <-chan int, sysErr chan<- error) *writer {
	return &writer{
		conn:            conn,
		minterProgramId: minterProgramId,
		mintManager:     mintManager,
		log:             log,
		msgChan:         make(chan msg.Message, msgLimit),
		stop:            stop,
		sysErr:          sysErr,
	}
}

func (w *writer) setRouter(r chains.Router) {
	w.router = r
}

func (w *writer) ResolveMessage(m msg.Message) bool {
	w.log.Info("Attempting to resolve message", "type", m.Type, "src", m.Source, "dst", m.Destination, "nonce", m.DepositNonce, "rId", m.ResourceId.Hex())
	w.log.Info("ResolveMessage: size of msgChan", "size", len(w.msgChan))
	w.msgChan <- m
	return true
}

// resolve msg from other chains
func (w *writer) processMessage(m msg.Message) (processOk bool) {
	switch m.Type {
	case msg.FungibleTransfer:
		poolClient := w.conn.poolClient
		rpcClient := w.conn.GetQueryClient()
		bigAmt := new(big.Int).SetBytes(m.Payload[0].([]byte))
		recipient := m.Payload[1].([]byte)
		toAccount := common.PublicKeyFromBytes(recipient)
		//check toAccount
		var toAccountInfo *tokenprog.TokenAccount
		var err error
		var retry = 0
		for {
			if retry > retryLimit {
				w.log.Error("GetTokenAccountInfo failed, will skip this recipient",
					"token account address", toAccount.ToBase58(),
					"err", err)
				return true
			}
			toAccountInfo, err = rpcClient.GetTokenAccountInfo(context.Background(), toAccount.ToBase58())
			if err != nil {
				// should skip if no account data
				if strings.Contains(err.Error(), "data length not match") {
					w.log.Warn("GetTokenAccountInfo failed, will skip",
						"token account address", toAccount.ToBase58(),
						"err", err)
					return true
				}
				// return false if retry limit
				w.log.Warn("GetTokenAccountInfo failed, will retry...",
					"token account address", toAccount.ToBase58(),
					"err", err)
				retry++
				time.Sleep(waitTime)
				continue
			}
			break
		}

		willUseProposalAccount, seed := GetProposalAccountPubkey(
			poolClient.ProposalBaseAccount.PublicKey,
			poolClient.BridgeProgramId,
			uint8(m.Source),
			uint8(m.Destination),
			uint64(m.DepositNonce),
		)

		//get gridgeAccount info
		bridgeAccount, err := rpcClient.GetBridgeAccountInfo(context.Background(), poolClient.BridgeAccountPubkey.ToBase58())
		if err != nil {
			w.log.Error("GetBridgeAccountInfo err",
				"bridge account address", poolClient.BridgeAccountPubkey.ToBase58(),
				"err", err)
			return false
		}
		var willUseMintAccount common.PublicKey
		if mint, exist := bridgeAccount.ResourceIdToMint[m.ResourceId]; !exist {
			w.log.Error("bridge resourceidToMint not exist",
				"resourceIdToMint", bridgeAccount.ResourceIdToMint,
				"resourceId", m.ResourceId)
			return false
		} else {
			willUseMintAccount = mint
		}

		//toAccount mint should equal willUseMintAccount
		if toAccountInfo.Mint != willUseMintAccount {
			w.log.Warn("TokenAccountInfo's mint account not equal, will skip this fungibleTransfer",
				"token account address", toAccount.ToBase58(),
				"mintAccount in tokenAccount", toAccountInfo.Mint.ToBase58(),
				"mintAccount in bridgeAccount", willUseMintAccount.ToBase58())
			return true
		}

		//check and create proposal is not exist
		_, err = rpcClient.GetMintProposalInfo(context.Background(), willUseProposalAccount.ToBase58())
		if err != nil && err == solClient.ErrAccountNotFound {
			sendOk := w.createProposalAccount(
				rpcClient,
				poolClient,
				toAccount,
				willUseProposalAccount,
				seed,
				m.ResourceId,
				bigAmt.Uint64(),
				"FungibleTransfer",
			)
			if !sendOk {
				return false
			}
		}
		if err != nil && err != solClient.ErrAccountNotFound {
			w.log.Error("GetMintProposalInfo err",
				"proposal account address", willUseProposalAccount.ToBase58(),
				"err", err)
			return false
		}

		//check proposal account is created
		create := w.waitingForProposalAccountCreate(rpcClient, willUseProposalAccount.ToBase58(), "FungibleTransfer")
		if !create {
			return false
		}
		w.log.Info("FungibleTransfer proposalAccount has create", "proposalAccount", willUseProposalAccount.ToBase58())

		valid := w.CheckProposalAccount(willUseProposalAccount, willUseMintAccount, toAccount, bigAmt.Uint64())
		if !valid {
			w.log.Info("FungibleTransfer CheckProposalAccount failed", "proposalAccount", willUseProposalAccount.ToBase58())
			return false
		}
		//if has exe just return
		isExe := w.IsProposalExe(willUseProposalAccount)
		if isExe {
			w.log.Info("FungibleTransfer proposalAccount has execute", "proposalAccount", willUseProposalAccount.ToBase58())
			return true
		}
		//approve proposal
		send := w.approveProposal(
			rpcClient,
			poolClient,
			willUseProposalAccount,
			willUseMintAccount,
			toAccount,
			w.mintManager,
			w.minterProgramId,
			"FungibleTransfer",
		)
		if !send {
			return false
		}

		//check proposal exe result
		exe := w.waitingForProposalExe(rpcClient, willUseProposalAccount.ToBase58(), "FungibleTransfer")
		if !exe {
			return false
		}
		w.log.Info("FungibleTransfer proposalAccount has execute", "proposalAccount", willUseProposalAccount.ToBase58())
		return true
	default:
		w.log.Warn("message type unsupported", "type", m.Type)
		return false
	}
}

func (w *writer) start() error {
	w.log.Debug("Starting solana writer...")
	go func() {
		for {
			select {
			case <-w.stop:
				close(w.msgChan)
				w.log.Info("solana writer stopped")
				return
			case msg := <-w.msgChan:
				result := w.processMessage(msg)
				w.log.Info("processMessage", "result", result)
				if !result {
					w.sysErr <- fmt.Errorf("processMessage failed")
				}
			}
		}
	}()

	return nil
}
