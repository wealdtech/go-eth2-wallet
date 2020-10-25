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
	"testing"

	"github.com/stretchr/testify/require"
	wallet "github.com/wealdtech/go-eth2-wallet"
	unencrypted "github.com/wealdtech/go-eth2-wallet-encryptor-unencrypted"
)

func TestEncryptor(t *testing.T) {
	// Ensure default encryptor is set.
	require.Equal(t, "keystore", wallet.GetEncryptor())

	// Attempt to set a nil encryptor; should error.
	require.EqualError(t, wallet.UseEncryptor(nil), "no encryptor supplied")

	// Attempt to set a different encryptor.
	require.NoError(t, wallet.UseEncryptor(unencrypted.New()))

	// Confirm the encryptor has been set.
	require.Equal(t, "unencrypted", wallet.GetEncryptor())
}
