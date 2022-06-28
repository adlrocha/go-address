package address

import (
	"path"
	"strings"

	"github.com/multiformats/go-varint"
)

var id0, _ = NewIDAddress(0)

const (
	ROOT_STR          = "/root"
	SUBNET_SEPARATOR  = "/"
	UNDEF_STR         = SUBNET_SEPARATOR
	HC_ADDR_SEPARATOR = ":"
)

// RootSubnet is the ID of the root network
var RootSubnet = SubnetID{
	Parent: ROOT_STR,
	Actor:  id0,
}

// UndefSubnetID is the undef ID
var UndefSubnetID = SubnetID{
	Parent: UNDEF_STR,
	Actor:  id0,
}

// SubnetID represents the ID of a subnet
type SubnetID struct {
	Parent string
	Actor  Address
}

// NewSubnetID generates the ID for a subnet from the networkName of its parent.
//
// It takes the parent name and adds the source address of the subnet actor that represents the subnet.
func NewSubnetID(parentName SubnetID, SubnetActorAddr Address) SubnetID {
	return SubnetID{
		Parent: parentName.String(),
		Actor:  SubnetActorAddr,
	}
}

func SubnetIDFromString(str string) (SubnetID, error) {
	switch str {
	case ROOT_STR:
		return RootSubnet, nil
	case UNDEF_STR:
		return UndefSubnetID, nil
	}

	s1 := strings.Split(str, SUBNET_SEPARATOR)
	actor, err := NewFromString(s1[len(s1)-1])
	if err != nil {
		return UndefSubnetID, err
	}
	return SubnetID{
		Parent: strings.Join(s1[:len(s1)-1], SUBNET_SEPARATOR),
		Actor:  actor,
	}, nil
}

// GetParent returns the ID of the parent network.
func (id SubnetID) GetParent() (SubnetID, error) {
	if id == RootSubnet {
		return UndefSubnetID, nil
	}
	return SubnetIDFromString(id.Parent)
}

// GetActor returns the subnet actor for a subnet
//
// Returns the address of the actor that handles the logic for a subnet
// in its parent Subnet.
func (id SubnetID) GetActor() Address {
	if id == RootSubnet {
		return Undef
	}
	return id.Actor
}

func (id SubnetID) CommonParent(other SubnetID) (SubnetID, int) {
	s1 := strings.Split(id.String(), SUBNET_SEPARATOR)
	s2 := strings.Split(other.String(), SUBNET_SEPARATOR)
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}
	out := SUBNET_SEPARATOR
	l := 0
	for i, s := range s2 {
		if s == s1[i] {
			out = path.Join(out, s)
			l = i
		} else {
			sn, err := SubnetIDFromString(out)
			if err != nil {
				return UndefSubnetID, 0
			}
			return sn, l
		}
	}
	sn, err := SubnetIDFromString(out)
	if err != nil {
		return UndefSubnetID, 0
	}
	return sn, l
}

func (id SubnetID) Down(curr SubnetID) SubnetID {
	s1 := strings.Split(id.String(), SUBNET_SEPARATOR)
	s2 := strings.Split(curr.String(), SUBNET_SEPARATOR)
	// curr needs to be contained in id
	if len(s2) >= len(s1) {
		return UndefSubnetID
	}
	_, l := id.CommonParent(curr)
	out := SUBNET_SEPARATOR
	for i := 0; i <= l+1 && i < len(s1); i++ {
		if i < len(s2) && s1[i] != s2[i] {
			// they are not in a common path
			return UndefSubnetID
		}
		out = path.Join(out, s1[i])
	}
	sn, err := SubnetIDFromString(out)
	if err != nil {
		return UndefSubnetID
	}
	return sn
}

func (id SubnetID) Up(curr SubnetID) SubnetID {
	s1 := strings.Split(id.String(), SUBNET_SEPARATOR)
	s2 := strings.Split(curr.String(), SUBNET_SEPARATOR)
	// curr needs to be contained in id
	if len(s2) > len(s1) {
		return UndefSubnetID
	}

	_, l := id.CommonParent(curr)
	out := SUBNET_SEPARATOR
	for i := 0; i < l; i++ {
		if i < len(s1) && s1[i] != s2[i] {
			// they are not in a common path
			return UndefSubnetID
		}
		out = path.Join(out, s1[i])
	}
	sn, err := SubnetIDFromString(out)
	if err != nil {
		return UndefSubnetID
	}
	return sn
}

// String returns the id in string form.
func (id SubnetID) String() string {
	if id == RootSubnet {
		return ROOT_STR
	}
	return strings.Join([]string{id.Parent, id.Actor.String()}, SUBNET_SEPARATOR)
}

// Subnet returns subnet information for an address if any.
func (a Address) Subnet() (SubnetID, error) {
	if a.str[0] != Hierarchical {
		return UndefSubnetID, ErrNotHierarchical
	}
	snSize, _, err := varint.FromUvarint([]byte(a.str[1:2]))
	if err != nil {
		return UndefSubnetID, err
	}
	return SubnetIDFromString(a.str[3 : snSize+3])
}

// RawAddr return the address without subnet context information.
func (a Address) RawAddr() (Address, error) {
	if a.str[0] != Hierarchical {
		return a, nil
	}
	snSize, _, err := varint.FromUvarint([]byte(a.str[1:2]))
	if err != nil {
		return Undef, err
	}
	return NewFromBytes([]byte(a.str[snSize+3:]))
}

func (a Address) PrettyPrint() string {
	if a.str[0] != Hierarchical {
		return a.String()
	}
	sn, _ := a.Subnet()
	raw, _ := a.RawAddr()
	return string(sn.String() + HC_ADDR_SEPARATOR + raw.String())
}
