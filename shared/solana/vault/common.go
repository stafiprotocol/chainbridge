// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vault

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)


const defaultRPCURL = "http://api.mainnet-beta.solana.com/rpc"

func GetPassword(input string) (string, error) {
	fd := os.Stdin.Fd()
	fmt.Printf(input)
	pass, err := terminal.ReadPassword(int(fd))
	fmt.Println("")
	return string(pass), err
}


func errorCheck(prefix string, err error) {
	if err != nil {
		fmt.Printf("ERROR: %s: %s\n", prefix, err)
		if strings.HasSuffix(err.Error(), "connection refused") && strings.Contains(err.Error(), defaultRPCURL) {
			fmt.Println("Have you selected a valid Solana JSON-RPC endpoint ? You can use the --rpc-url flag or SLNC_GLOBAL_RPC_URL environment variable.")
		}
		os.Exit(1)
	}
}
func MustGetWallet() *Vault {
	vault, err := setupWallet()
	errorCheck("wallet setup", err)
	return vault
}

func setupWallet() (*Vault, error) {
	walletFile := viper.GetString("global-vault-file")
	if _, err := os.Stat(walletFile); err != nil {
		return nil, fmt.Errorf("wallet file %q missing: %w", walletFile, err)
	}

	v, err := NewVaultFromWalletFile(walletFile)
	if err != nil {
		return nil, fmt.Errorf("loading vault: %w", err)
	}

	boxer, err := SecretBoxerForType(v.SecretBoxWrap)
	if err != nil {
		return nil, fmt.Errorf("secret boxer: %w", err)
	}

	if err := v.Open(boxer); err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}

	return v, nil
}

func GetDecryptPassphrase() (string, error) {
	passphrase, err := GetPassword("Enter passphrase to decrypt your solana vault: ")
	if err != nil {
		return "", fmt.Errorf("reading password: %s", err)
	}

	return passphrase, nil
}
func GetEncryptPassphrase() (string, error) {
	passphrase, err := GetPassword("Enter passphrase to encrypt your solana vault: ")
	if err != nil {
		return "", fmt.Errorf("reading password: %s", err)
	}

	passphraseConfirm, err := GetPassword("Confirm passphrase: ")
	if err != nil {
		return "", fmt.Errorf("reading confirmation password: %s", err)
	}

	if passphrase != passphraseConfirm {
		fmt.Println()
		return "", errors.New("passphrase mismatch!")
	}
	return passphrase, nil

}