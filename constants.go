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
)

// UndefAddressString is the string used to represent an empty address when encoded to a string.
var UndefAddressString = "<empty>"

// PayloadHashLength defines the hash length taken over addresses using the Actor and SECP256K1 protocols.
const PayloadHashLength = 20

// ChecksumHashLength defines the hash length used for calculating address checksums.
const ChecksumHashLength = 4

// MaxAddressStringLength is the max length of an address encoded as a string
// it includes the network prefix, protocol, and bls publickey
// NOTE: To accommodate consensus hierarchies of up to 6 levels in
// hierarchical addresses we add an additional length buffer.
// For the MVP we'll leave it like this, but in the future we may want to
// support constant-length IDs for subnets, to allow us to set
// this MaxLength accurately without worrying about overflows.
const MaxAddressStringLength = 2 + 84 + (4*6 + 6)
const HierarchicalLength = MaxAddressStringLength - 32

// BlsPublicKeyBytes is the length of a BLS public key
const BlsPublicKeyBytes = 48

// BlsPrivateKeyBytes is the length of a BLS private key
const BlsPrivateKeyBytes = 32

var payloadHashConfig = &blake2b.Config{Size: PayloadHashLength}
var checksumHashConfig = &blake2b.Config{Size: ChecksumHashLength}

const encodeStd = "abcdefghijklmnopqrstuvwxyz234567"

// AddressEncoding defines the base32 config used for address encoding and decoding.
var AddressEncoding = base32.NewEncoding(encodeStd)
