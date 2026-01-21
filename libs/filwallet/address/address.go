package address

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/go-address"
)

const (
	evmNamespace uint64 = 10
)

// Parse takes a raw string and returns an Address struct.
func Parse(raw string) (Address, error) {
	if len(raw) < 3 {
		return Address{}, errors.New("address too short")
	}

	var t Type
	switch raw[:2] {
	case "f1":
		t = TypeF1
	case "f3":
		t = TypeF3
	case "f4":
		t = TypeF4
	case "0x":
		t = Type0X
	default:
		return Address{}, fmt.Errorf("unknown prefix")
	}
	return Address{Type: t, Value: raw}, nil
}

// DeriveAddressesFromPrivateKey returns all standard Filecoin address formats
func DeriveAddressesFromPrivateKey(privKey *ecdsa.PrivateKey) ([]Address, error) {
	// f1 address (legacy secp256k1)
	pubBytes := crypto.CompressPubkey(&privKey.PublicKey)
	f1Addr, err := address.NewSecp256k1Address(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("derive f1 address: %w", err)
	}

	// f4 address (delegated/eth address)
	ethAddr := crypto.PubkeyToAddress(privKey.PublicKey) // Ethereum address
	f4Addr, err := address.NewDelegatedAddress(evmNamespace, ethAddr.Bytes())
	if err != nil {
		return nil, fmt.Errorf("create f4 address: %w", err)
	}

	return []Address{
		{Type: TypeF1, Value: f1Addr.String()},
		{Type: TypeF4, Value: f4Addr.String()},
		{Type: Type0X, Value: ethAddr.Hex()},
	}, nil
}
