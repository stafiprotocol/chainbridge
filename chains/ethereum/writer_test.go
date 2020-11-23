// Copyright 2020 Stafi Protocol
// SPDX-License-Identifier: LGPL-3.0-only

package ethereum

import (
	"testing"
)

func TestWriter_start_stop(t *testing.T) {
	conn := newLocalConnection(t, aliceTestConfig)
	defer conn.Close()

	stop := make(chan int)
	writer := NewWriter(conn, aliceTestConfig, TestLogger, stop, nil, nil)

	err := writer.start()
	if err != nil {
		t.Fatal(err)
	}

	// Initiate shutdown
	close(stop)
}
