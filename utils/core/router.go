// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package core

import (
	"fmt"
	"sync"

	log "github.com/ChainSafe/log15"
	"github.com/stafiprotocol/chainbridge/utils/msg"
)

// Writer consumes a message and makes the requried on-chain interactions.
type Writer interface {
	ResolveMessage(message msg.Message) bool
}

// Router forwards messages from their source to their destination
type Router struct {
	registry map[msg.ChainId]Writer
	lock     *sync.RWMutex
	log      log.Logger
}

func NewRouter(log log.Logger) *Router {
	return &Router{
		registry: make(map[msg.ChainId]Writer),
		lock:     &sync.RWMutex{},
		log:      log,
	}
}

// Send passes a message to the destination Writer if it exists
func (r *Router) Send(msg msg.Message) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.log.Trace("Routing message", "src", msg.Source, "dest", msg.Destination, "nonce", msg.DepositNonce, "rId", msg.ResourceId.Hex())
	w := r.registry[msg.Destination]
	if w == nil {
		return fmt.Errorf("unknown destination chainId: %d", msg.Destination)
	}

	go w.ResolveMessage(msg)
	return nil
}

func (r *Router) SupportChainId(chainId msg.ChainId) bool {
	r.lock.Lock()
	defer r.lock.Unlock()

	w := r.registry[chainId]
	if w == nil {
		return false
	} else {
		return true
	}
}

// Listen registers a Writer with a ChainId which Router.Send can then use to propagate messages
func (r *Router) Listen(id msg.ChainId, w Writer) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.log.Debug("Registering new chain in router", "id", id)
	r.registry[id] = w
}
