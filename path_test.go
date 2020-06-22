// Copyright © 2020 Weald Technology Trading
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package wallet_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	wallet "github.com/wealdtech/go-eth2-wallet"
	"gotest.tools/assert"
)

func TestWalletAndAccountNames(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		err         error
		walletName  string
		accountName string
	}{
		{
			name: "Nil",
			err:  errors.New("invalid account format"),
		},
		{
			name: "EmptyPath",
			path: "",
			err:  errors.New("invalid account format"),
		},
		{
			name: "SeparatorOnly",
			path: "/",
			err:  errors.New("invalid account format"),
		},
		{
			name: "AccountOnly",
			path: "/Account",
			err:  errors.New("invalid account format"),
		},
		{
			name:       "WalletOnly",
			path:       "Wallet",
			walletName: "Wallet",
		},
		{
			name:       "WalletOnlyTrailingSeparator",
			path:       "Wallet/",
			walletName: "Wallet",
		},
		{
			name:        "DoubleSeparator",
			path:        "Wallet//",
			walletName:  "Wallet",
			accountName: "/",
		},
		{
			name:        "WalletAndAccount",
			path:        "Wallet/Account",
			walletName:  "Wallet",
			accountName: "Account",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			walletName, accountName, err := wallet.WalletAndAccountNames(test.path)
			if test.err != nil {
				require.NotNil(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				require.Nil(t, err)
				assert.Equal(t, test.walletName, walletName)
				assert.Equal(t, test.accountName, accountName)
			}
		})
	}
}
