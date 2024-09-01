// Copyright 2019 - 2024 Weald Technology Trading.
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
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-ecodec"
	distributed "github.com/wealdtech/go-eth2-wallet-distributed"
	hd "github.com/wealdtech/go-eth2-wallet-hd/v2"
	keystore "github.com/wealdtech/go-eth2-wallet-keystore"
	nd "github.com/wealdtech/go-eth2-wallet-nd/v2"
	wtypes "github.com/wealdtech/go-eth2-wallet-types/v2"
)

const (
	nonDeterministicWallet          = "non-deterministic"
	hierarchicalDeterministicWallet = "hierarchical deterministic"
	distributedWallet               = "distributed"
	keystoreWallet                  = "keystore"
)

// walletOptions are the optons used when opening and creating wallets.
type walletOptions struct {
	store      wtypes.Store
	encryptor  wtypes.Encryptor
	walletType string
	passphrase []byte
	seed       []byte
}

// Option gives options to OpenWallet and CreateWallet.
type Option interface {
	apply(opts *walletOptions)
}

type optionFunc func(*walletOptions)

func (f optionFunc) apply(o *walletOptions) {
	f(o)
}

// WithStore sets the store for the wallet.
func WithStore(store wtypes.Store) Option {
	return optionFunc(func(o *walletOptions) {
		o.store = store
	})
}

// WithEncryptor sets the encryptor for the wallet.
func WithEncryptor(encryptor wtypes.Encryptor) Option {
	return optionFunc(func(o *walletOptions) {
		o.encryptor = encryptor
	})
}

// WithPassphrase sets the passphrase for the wallet.
func WithPassphrase(passphrase []byte) Option {
	return optionFunc(func(o *walletOptions) {
		o.passphrase = passphrase
	})
}

// WithType sets the type for the wallet.
func WithType(walletType string) Option {
	return optionFunc(func(o *walletOptions) {
		o.walletType = walletType
	})
}

// WithSeed sets the seed for a hierarchical deterministic wallet.
func WithSeed(seed []byte) Option {
	return optionFunc(func(o *walletOptions) {
		o.seed = seed
	})
}

// ImportWallet imports a wallet from its encrypted export.
func ImportWallet(encryptedData []byte, passphrase []byte) (wtypes.Wallet, error) {
	type walletExt struct {
		Wallet *walletInfo `json:"wallet"`
	}

	data, err := ecodec.Decrypt(encryptedData, passphrase)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt wallet")
	}

	ext := &walletExt{}
	err = json.Unmarshal(data, ext)
	if err != nil {
		return nil, errors.Wrap(err, "failed to import wallet")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var wallet wtypes.Wallet
	switch ext.Wallet.Type {
	case "nd", nonDeterministicWallet:
		wallet, err = nd.Import(ctx, encryptedData, passphrase, store, encryptor)
	case "hd", hierarchicalDeterministicWallet:
		wallet, err = hd.Import(ctx, encryptedData, passphrase, store, encryptor)
	case distributedWallet:
		wallet, err = distributed.Import(ctx, encryptedData, passphrase, store, encryptor)
	case keystoreWallet:
		wallet, err = keystore.Import(ctx, encryptedData, passphrase, store, encryptor)
	default:
		return nil, fmt.Errorf("unsupported wallet type %q", ext.Wallet.Type)
	}

	return wallet, err
}

// OpenWallet opens an existing wallet.
// If the wallet does not exist an error is returned.
func OpenWallet(name string, opts ...Option) (wtypes.Wallet, error) {
	options := walletOptions{
		store:     store,
		encryptor: encryptor,
	}
	for _, o := range opts {
		if opts != nil {
			o.apply(&options)
		}
	}
	if options.store == nil {
		return nil, errors.New("no store specified")
	}
	if options.encryptor == nil {
		return nil, errors.New("no encryptor specified")
	}

	data, err := options.store.RetrieveWallet(name)
	if err != nil {
		return nil, err
	}

	return walletFromBytes(data, options.store, options.encryptor)
}

// CreateWallet creates a wallet.
// If the wallet already exists an error is returned.
func CreateWallet(name string, opts ...Option) (wtypes.Wallet, error) {
	options := walletOptions{
		store:      store,
		encryptor:  encryptor,
		passphrase: nil,
		walletType: "nd",
		seed:       nil,
	}
	for _, o := range opts {
		if o != nil {
			o.apply(&options)
		}
	}

	if err := checkOptions(&options); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	switch options.walletType {
	case "nd", nonDeterministicWallet:
		return nd.CreateWallet(ctx, name, options.store, options.encryptor)
	case "hd", hierarchicalDeterministicWallet:
		return hd.CreateWallet(ctx, name, options.passphrase, options.store, options.encryptor, options.seed)
	case distributedWallet:
		return distributed.CreateWallet(ctx, name, options.store, options.encryptor)
	case keystoreWallet:
		return keystore.CreateWallet(ctx, name, options.store, options.encryptor)
	default:
		return nil, fmt.Errorf("unhandled wallet type %q", options.walletType)
	}
}

func checkOptions(options *walletOptions) error {
	if options.store == nil {
		return errors.New("no store specified")
	}
	if options.encryptor == nil {
		return errors.New("no encryptor specified")
	}
	if (options.walletType == "hd" || options.walletType == hierarchicalDeterministicWallet) && options.seed == nil {
		return errors.New("no seed specified")
	}

	return nil
}

type walletInfo struct {
	ID   uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

// Wallets provides information on the available wallets.
func Wallets(opts ...Option) <-chan wtypes.Wallet {
	ch := make(chan wtypes.Wallet, 1024)

	options := walletOptions{
		store:     store,
		encryptor: encryptor,
	}
	for _, o := range opts {
		if opts != nil {
			o.apply(&options)
		}
	}
	if options.store == nil {
		return ch
	}
	if options.encryptor == nil {
		return ch
	}

	go func() {
		for data := range options.store.RetrieveWallets() {
			wallet, err := walletFromBytes(data, options.store, options.encryptor)
			if err == nil {
				ch <- wallet
			}
		}
		close(ch)
	}()

	return ch
}

func walletFromBytes(data []byte, store wtypes.Store, encryptor wtypes.Encryptor) (wtypes.Wallet, error) {
	if store == nil {
		return nil, errors.New("no store specified")
	}
	if encryptor == nil {
		return nil, errors.New("no encryptor specified")
	}

	info := &walletInfo{}
	err := json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var wallet wtypes.Wallet
	switch info.Type {
	case "nd", nonDeterministicWallet:
		wallet, err = nd.DeserializeWallet(ctx, data, store, encryptor)
	case "hd", hierarchicalDeterministicWallet:
		wallet, err = hd.DeserializeWallet(ctx, data, store, encryptor)
	case distributedWallet:
		wallet, err = distributed.DeserializeWallet(ctx, data, store, encryptor)
	case keystoreWallet:
		wallet, err = keystore.DeserializeWallet(ctx, data, store, encryptor)
	default:
		return nil, fmt.Errorf("unsupported wallet type %q", info.Type)
	}

	return wallet, err
}
