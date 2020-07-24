# go-eth2-wallet

[![Tag](https://img.shields.io/github/tag/wealdtech/go-eth2-wallet.svg)](https://github.com/wealdtech/go-eth2-wallet/releases/)
[![License](https://img.shields.io/github/license/wealdtech/go-eth2-wallet.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/wealdtech/go-eth2-wallet?status.svg)](https://godoc.org/github.com/wealdtech/go-eth2-wallet)
[![Travis CI](https://img.shields.io/travis/wealdtech/go-eth2-wallet.svg)](https://travis-ci.org/wealdtech/go-eth2-wallet)
[![codecov.io](https://img.shields.io/codecov/c/github/wealdtech/go-eth2-wallet.svg)](https://codecov.io/github/wealdtech/go-eth2-wallet)

Go library to provide access to Ethereum 2 wallets with advanced features.

** Please note that this library uses standards that are not yet final, and as such may result in changes that alter public and private keys.  Do not use this library for production use just yet **

## Table of Contents

- [Install](#install)
- [Usage](#usage)
- [Maintainers](#maintainers)
- [Contribute](#contribute)
- [License](#license)

## Install

`go-eth2-wallet` is a standard Go module which can be installed with:

```sh
go get github.com/wealdtech/go-eth2-wallet
```

## Usage

Before using `go-eth2-wallet` it is important to understand the concepts behind it.

An *account* is the container for private keys.  Each account is named, and contains information about the derivation of the private key: either the key itself, or information that can be passed to a third-party system to identify the key (_e.g._ for a hardware wallet).  An account may or may not have a passphrase depending on the information stored within it.

A *wallet* is the container for accounts.  Each wallet is named, and contains information about how to create new accounts: one wallet might create accounts by generating random keys, another might use a seed phrase to create accounts by a deterministic method.  A wallet may or may not have a passphrase depending on the information stored within it.

A *store* is a storage system for wallets and accounts.  The store allows access to wallets and accounts regardless of where they are.  Stores can be local, for example on the filesystem, or remote, for example on Amazon's S3 storage.  Stores have their own configuration, which may or may not include account IDs and passwords.  Stores can be encrypted with a passphrase, in which case all data written to the store is encrypted prior to being written; this increases security where there are concerns about others accessing the data.

An overview of the architecture as laid out above is shown below, showing an application interacting with multiple local and remote stores:

![Overview](images/overview.svg)

And to recap: an application can access one or more stores.  Each store contains a number of wallets, and each wallet contains a number of accounts.

The Ethereum 2 wallet is designed to be highly modular: additional stores and wallets can be added, allowing for the ability to upgrade features as required.  This module hides the complexity that comes with that modularity and provides a simple unified API to handle all common (and not-so-common) wallet operations.

### Stores

The following stores are available:

   - [filesystem](https://github.com/wealdtech/go-eth2-wallet-store-filesystem): this stores wallets and accounts on the local filesystem.
   - [s3](https://github.com/wealdtech/go-eth2-wallet-store-s3): this stores wallets and accounts on Amazon S3.
   - [scratch](https://github.com/wealdtech/go-eth2-wallet-store-scratch): this stores wallets and accounts in memory.

Please refer to the documentation for each store to understand its functionality and available options.

#### Wallets

The following wallet types are available:

  - [nd](https://github.com/wealdtech/go-eth2-wallet-nd): this is a traditional non-deterministic wallet where private keys are generated randomly and have no relationship to each other.
  - [hd](https://github.com/wealdtech/go-eth2-wallet-hd): this is a hierarchical deterministic wallet where private keys are generated based on a seed phrase and path.
  - [distributed](https://github.com/wealdtech/go-eth2-wallet-distributed): this is a wallet whose accounts form part of a distributed composite.


Please refer to the documentation for each wallet type to understand its functionality and available options.

### Examples

#### Opening an existing wallet from the default store

```go
import (
    e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {
    wallet, err := e2wallet.OpenWallet("my wallet")
    if err != nil {
        panic(err)
    }

    ...
}
```

#### Creating a new wallet on the default store

```go
import (
    e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {
    wallet, err := e2wallet.CreateWallet("my wallet")
    if err != nil {
        panic(err)
    }

    ...
}
```

#### Creating a new wallet on the Amazon S3 store

```go
import (
    e2wallet "github.com/wealdtech/go-eth2-wallet"
    s3 "github.com/wealdtech/go-eth2-wallet-store-s3"
)

func main() {
    s3Store, err := s3.New(s3.WithPassphrase([]byte("store secret")))
    if err != nil {
        panic(err)
    }
    err = e2wallet.UseStore(s3Store)
    if err != nil {
        panic(err)
    }
    wallet, err := e2wallet.CreateWallet("my wallet")

    ...
}
```

### List all wallets, and all accounts in each wallet

```go
import (
    "fmt"

    e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {
    for wallet := range e2wallet.Wallets() {
        fmt.Printf("Found wallet %s\n", wallet.Name())
        for account := range wallet.Accounts() {
            fmt.Printf("Wallet %s has account %s\n", wallet.Name(), account.Name())
        }
    }
}
```

### Creating an account in an existing wallet

```go
import (
    e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {
    wallet, err := e2wallet.OpenWallet("my wallet")
    if err != nil {
        panic(err)
    }

    err = wallet.Unlock([]byte("wallet passphrase"))
    if err != nil {
        panic(err)
    }
    // Always immediately defer locking the wallet to ensure it does not remain unlocked outside of the function.
    defer wallet.Lock()

    account, err := wallet.CreateAccount("primary account", []byte("secret passphrase"))
    if err != nil {
        panic(err)
    }
    // Wallet should be locked as soon as unlocked operations have finished; it is safe to explicitly call wallet.Lock() as well
    // as defer it as per above.
    wallet.Lock()

    ...
}
```

### Signing data and verifying signatures

```go
import (
    "errors"

    e2wallet "github.com/wealdtech/go-eth2-wallet"
)

func main() {
    wallet, err := e2wallet.OpenWallet("my wallet")
    if err != nil {
        panic(err)
    }
    account, err := wallet.AccountByName("primary account")
    if err != nil {
        panic(err)
    }

    err = account.Unlock([]byte("my secret passphrase"))
    if err != nil {
        panic(err)
    }
    // Always immediately defer locking the wallet to ensure it does not remain unlocked outside of the function.
    defer account.Lock()

    signature, err := account.Sign([]byte("some data to sign"))
    if err != nil {
        panic(err)
    }
    // Wallet should be locked as soon as unlocked operations have finished; it is safe to explicitly call wallet.Lock() as well
    // as defer it as per above.
    account.Lock()
    
    verified := signature.Verify([]byte("some data to sign"), account.PublicKey())
    if !verified {
        panic(errors.New("failed to verify signature"))
    }

    ...
}
```

## Maintainers

Jim McDonald: [@mcdee](https://github.com/mcdee).

## Contribute

Contributions welcome. Please check out [the issues](https://github.com/wealdtech/go-eth2-wallet/issues).

## License

[Apache-2.0](LICENSE) © 2019 Weald Technology Trading Ltd
