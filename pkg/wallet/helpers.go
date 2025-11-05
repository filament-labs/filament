package wallet

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/filecoin-project/go-address"
)

type DerivedAddrs struct {
	F1  address.Address
	F4  address.Address
	Hex string
}

// DeriveFromECDSA derives f1, f4 and 0x variants from an ECDSA private key.
func DeriveFromECDSA(priv *ecdsa.PrivateKey) (*DerivedAddrs, error) {
	pubBytes := ethcrypto.FromECDSAPub(&priv.PublicKey) // 65 bytes uncompressed
	f1, err := address.NewSecp256k1Address(pubBytes)
	if err != nil {
		return nil, fmt.Errorf("derive f1: %w", err)
	}

	ethAddr := ethcrypto.PubkeyToAddress(priv.PublicKey)
	ethBytes := ethAddr.Bytes()
	f4, err := address.NewDelegatedAddress(uint64(address.Delegated), ethBytes)
	if err != nil {
		return nil, fmt.Errorf("derive f4: %w", err)
	}

	return &DerivedAddrs{F1: f1, F4: f4, Hex: "0x" + hex.EncodeToString(ethBytes)}, nil
}
