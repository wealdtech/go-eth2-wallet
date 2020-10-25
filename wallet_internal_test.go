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

package wallet

import (
	"testing"

	"github.com/stretchr/testify/require"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
	scratch "github.com/wealdtech/go-eth2-wallet-store-scratch"
	e2wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

func TestWalletFromBytes(t *testing.T) {
	tests := []struct {
		name      string
		store     e2wtypes.Store
		encryptor e2wtypes.Encryptor
		input     []byte
		err       string
	}{
		{
			name: "Nil",
			err:  "no store specified",
		},
		{
			name:      "StoreNil",
			input:     []byte(`{"name":"ND test","type":"non-deterministic","uuid":"e45d4f2c-00e9-44ee-98b3-ea12d4d937a9","version":1}`),
			encryptor: keystorev4.New(),
			err:       "no store specified",
		},
		{
			name:  "EncryptorNil",
			input: []byte(`{"name":"ND test","type":"non-deterministic","uuid":"e45d4f2c-00e9-44ee-98b3-ea12d4d937a9","version":1}`),
			store: scratch.New(),
			err:   "no encryptor specified",
		},
		{
			name:      "DataMissing",
			store:     scratch.New(),
			encryptor: keystorev4.New(),
			err:       "unexpected end of JSON input",
		},
		{
			name:      "DataBad",
			input:     []byte(`x`),
			store:     scratch.New(),
			encryptor: keystorev4.New(),
			err:       "invalid character 'x' looking for beginning of value",
		},
		{
			name:      "TypeUnknown",
			input:     []byte(`{"name":"ND test","type":"unknown","uuid":"e45d4f2c-00e9-44ee-98b3-ea12d4d937a9","version":1}`),
			store:     scratch.New(),
			encryptor: keystorev4.New(),
			err:       "unsupported wallet type \"unknown\"",
		},
		{
			name:      "HDGood",
			input:     []byte(`{"crypto":{"checksum":{"function":"sha256","message":"b13cf3ccd0924f5611a323d30ebbb2259ed155f8e77e3299f4a42cae79b104c4","params":{}},"cipher":{"function":"aes-128-ctr","message":"7082de0a9e179364de9f1973449056b6a297f905a3da5e6700d93e0636d85a8bc2114a196ba16e93db156330372c29c5535faa5473898b44b5563915f089ad41","params":{"iv":"32e2dd6b1566b79b642c298efcbe5d3c"}},"kdf":{"function":"pbkdf2","message":"","params":{"c":16,"dklen":32,"prf":"hmac-sha256","salt":"29e5f413252246f8ab6aff735d5ca807ba6cc76251aa650efae2aefc65ac7a91"}}},"name":"HD test","nextaccount":0,"type":"hierarchical deterministic","uuid":"2d67faca-a781-4ec6-aec8-5d7f520f55a9","version":1}`),
			store:     scratch.New(),
			encryptor: keystorev4.New(),
		},
		{
			name:      "NDGood",
			input:     []byte(`{"name":"ND test","type":"non-deterministic","uuid":"e45d4f2c-00e9-44ee-98b3-ea12d4d937a9","version":1}`),
			store:     scratch.New(),
			encryptor: keystorev4.New(),
		},
		{
			name:      "DistributedGood",
			input:     []byte(`{"name":"Distributed test","type":"distributed","uuid":"2aafba72-d748-498f-8f3f-eae3dc9c36c1","version":1}`),
			store:     scratch.New(),
			encryptor: keystorev4.New(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wallet, err := walletFromBytes(test.input, test.store, test.encryptor)
			if test.err == "" {
				require.NoError(t, err)
				require.NotNil(t, wallet)
			} else {
				require.EqualError(t, err, test.err)
			}
		})
	}
}
