package address

import (
	"encoding/base32"
	"errors"

	"github.com/minio/blake2b-simd"
)

func init() {

	var err error

	TestAddress, err = NewActorAddress([]byte("satoshi"))
	if err != nil {
		panic(err)
	}

	TestAddress2, err = NewActorAddress([]byte("nakamoto"))
	if err != nil {
		panic(err)
	}
}

var (
	// TestAddress is an account with some initial funds in it.
	TestAddress Address
	// TestAddress2 is an account with some initial funds in it.
	TestAddress2 Address
)

var (
	// ErrUnknownNetwork is returned when encountering an unknown network in an address.
	ErrUnknownNetwork = errors.New("unknown address network")

	// ErrUnknownProtocol is returned when encountering an unknown protocol in an address.
	ErrUnknownProtocol = errors.New("unknown address protocol")
	// ErrInvalidPayload is returned when encountering an invalid address payload.
	ErrInvalidPayload = errors.New("invalid address payload")
	// ErrInvalidLength is returned when encountering an address of invalid length.
	ErrInvalidLength = errors.New("invalid address length")
	// ErrInvalidChecksum is returned when encountering an invalid address checksum.
	ErrInvalidChecksum = errors.New("invalid address checksum")
	// ErrNotHierarchical is returned when trying to access info only available in hierarchical addresses
	ErrNotHierarchical = errors.New("not hierarchical address")
	// ErrInvalidEncoding is returned when encountering a non-standard encoding of an address.
	ErrInvalidEncoding = errors.New("invalid encoding")
)

// UndefAddressString is the string used to represent an empty address when encoded to a string.
var UndefAddressString = "<empty>"

// MaxInt64StringLength defines the maximum length of `int64` as a string.
const MaxInt64StringLength = 19

// PayloadHashLength defines the hash length taken over addresses using the Actor and SECP256K1 protocols.
const PayloadHashLength = 20

// ChecksumHashLength defines the hash length used for calculating address checksums.
const ChecksumHashLength = 4

// MaxAddressStringLength is the max length of an address encoded as a string
// it includes the network prefix, protocol, and bls publickey (see spec)
// (142 bytes HA payload * 1.6 overhead base32 + 6 bytes checkpoint)
const MaxAddressStringLength = 232
const HierarchicalLength = 142

// BlsPublicKeyBytes is the length of a BLS public key
const BlsPublicKeyBytes = 48

// BlsPrivateKeyBytes is the length of a BLS private key
const BlsPrivateKeyBytes = 32

var payloadHashConfig = &blake2b.Config{Size: PayloadHashLength}
var checksumHashConfig = &blake2b.Config{Size: ChecksumHashLength}

const encodeStd = "abcdefghijklmnopqrstuvwxyz234567"

// AddressEncoding defines the base32 config used for address encoding and decoding.
var AddressEncoding = base32.NewEncoding(encodeStd)
