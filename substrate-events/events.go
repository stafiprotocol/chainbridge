package events

import (
  "github.com/stafiprotocol/go-substrate-rpc-client/types"
)

type Events struct {
  ChainBridge_FungibleTransfer        []EventFungibleTransfer        //nolint:stylecheck,golint
  ChainBridge_NonFungibleTransfer     []EventNonFungibleTransfer     //nolint:stylecheck,golint
  ChainBridge_GenericTransfer         []EventGenericTransfer         //nolint:stylecheck,golint
  ChainBridge_RelayerThresholdChanged []EventRelayerThresholdChanged //nolint:stylecheck,golint
  ChainBridge_ChainWhitelisted        []EventChainWhitelisted        //nolint:stylecheck,golint
  ChainBridge_RelayerAdded            []EventRelayerAdded            //nolint:stylecheck,golint
  ChainBridge_RelayerRemoved          []EventRelayerRemoved          //nolint:stylecheck,golint
  ChainBridge_VoteFor                 []EventVoteFor                 //nolint:stylecheck,golint
  ChainBridge_VoteAgainst             []EventVoteAgainst             //nolint:stylecheck,golint
  ChainBridge_ProposalApproved        []EventProposalApproved        //nolint:stylecheck,golint
  ChainBridge_ProposalRejected        []EventProposalRejected        //nolint:stylecheck,golint
  ChainBridge_ProposalSucceeded       []EventProposalSucceeded       //nolint:stylecheck,golint
  ChainBridge_ProposalFailed          []EventProposalFailed          //nolint:stylecheck,golint
}

type EventFungibleTransfer struct {
  Phase        types.Phase
  Destination  types.U8
  DepositNonce types.U64
  ResourceId   types.Bytes32
  Amount       types.U256
  Recipient    types.Bytes
  Topics       []types.Hash
}

type EventNonFungibleTransfer struct {
  Phase        types.Phase
  Destination  types.U8
  DepositNonce types.U64
  ResourceId   types.Bytes32
  TokenId      types.Bytes
  Recipient    types.Bytes
  Metadata     types.Bytes
  Topics       []types.Hash
}

type EventGenericTransfer struct {
  Phase        types.Phase
  Destination  types.U8
  DepositNonce types.U64
  ResourceId   types.Bytes32
  Metadata     types.Bytes
  Topics       []types.Hash
}

type EventRelayerThresholdChanged struct {
  Phase     types.Phase
  Threshold types.U32
  Topics    []types.Hash
}

type EventChainWhitelisted struct {
  Phase   types.Phase
  ChainId types.U8
  Topics  []types.Hash
}

type EventRelayerAdded struct {
  Phase   types.Phase
  Relayer types.AccountID
  Topics  []types.Hash
}

type EventRelayerRemoved struct {
  Phase   types.Phase
  Relayer types.AccountID
  Topics  []types.Hash
}

type EventVoteFor struct {
  Phase        types.Phase
  SourceId     types.U8
  DepositNonce types.U64
  Voter        types.AccountID
  Topics       []types.Hash
}

type EventVoteAgainst struct {
  Phase        types.Phase
  SourceId     types.U8
  DepositNonce types.U64
  Voter        types.AccountID
  Topics       []types.Hash
}

type EventProposalApproved struct {
  Phase        types.Phase
  SourceId     types.U8
  DepositNonce types.U64
  Topics       []types.Hash
}

type EventProposalRejected struct {
  Phase        types.Phase
  SourceId     types.U8
  DepositNonce types.U64
  Topics       []types.Hash
}

type EventProposalSucceeded struct {
  Phase        types.Phase
  SourceId     types.U8
  DepositNonce types.U64
  Topics       []types.Hash
}

type EventProposalFailed struct {
  Phase        types.Phase
  SourceId     types.U8
  DepositNonce types.U64
  Topics       []types.Hash
}