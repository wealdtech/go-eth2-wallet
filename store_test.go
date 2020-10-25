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
	scratch "github.com/wealdtech/go-eth2-wallet-store-scratch"
)

func TestStore(t *testing.T) {
	// Ensure default store is set.
	require.Equal(t, "filesystem", wallet.GetStore())

	// Attempt to set a nil store; should error.
	require.EqualError(t, wallet.UseStore(nil), "no store supplied")

	// Attempt to set a different store.
	require.NoError(t, wallet.UseStore(scratch.New()))

	// Confirm the store has been set.
	require.Equal(t, "scratch", wallet.GetStore())

	// Attempt to set different stores.
	require.NoError(t, wallet.SetStore("filesystem", nil))
	require.Equal(t, "filesystem", wallet.GetStore())
	require.NoError(t, wallet.SetStore("s3", nil))
	require.Equal(t, "s3", wallet.GetStore())
	require.EqualError(t, wallet.SetStore("unknown", nil), "unknown wallet store \"unknown\"")
}
