package address

import (
	"path"
	"strings"

	"github.com/multiformats/go-varint"
)

var id0, _ = NewIDAddress(0)

const (
	RootStr         = "/root"
	SubnetSeparator = "/"
	UndefStr        = SubnetSeparator
	HCAddrSeparator = ":"
)

// RootSubnet is the ID of the root network
var RootSubnet = SubnetID{
	Parent: RootStr,
	Actor:  id0,
}

// UndefSubnetID is the undef ID
var UndefSubnetID = SubnetID{
	Parent: UndefStr,
	Actor:  id0,
}

// SubnetID represents the ID of a subnet
type SubnetID struct {
	Parent string
	Actor  Address
}

func (id SubnetID) Key() string {
	return id.String()
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
	case RootStr:
		return RootSubnet, nil
	case UndefStr:
		return UndefSubnetID, nil
	}

	s1 := strings.Split(str, SubnetSeparator)
	actor, err := NewFromString(s1[len(s1)-1])
	if err != nil {
		return UndefSubnetID, err
	}
	return SubnetID{
		Parent: strings.Join(s1[:len(s1)-1], SubnetSeparator),
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
	s1 := strings.Split(id.String(), SubnetSeparator)
	s2 := strings.Split(other.String(), SubnetSeparator)
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}
	out := SubnetSeparator
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
	s1 := strings.Split(id.String(), SubnetSeparator)
	s2 := strings.Split(curr.String(), SubnetSeparator)
	// curr needs to be contained in id
	if len(s2) >= len(s1) {
		return UndefSubnetID
	}
	_, l := id.CommonParent(curr)
	out := SubnetSeparator
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
	s1 := strings.Split(id.String(), SubnetSeparator)
	s2 := strings.Split(curr.String(), SubnetSeparator)
	// curr needs to be contained in id
	if len(s2) > len(s1) {
		return UndefSubnetID
	}

	_, l := id.CommonParent(curr)
	out := SubnetSeparator
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
		return RootStr
	}
	return strings.Join([]string{id.Parent, id.Actor.String()}, SubnetSeparator)
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
	return string(sn.String() + HCAddrSeparator + raw.String())
}
