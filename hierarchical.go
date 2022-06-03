package address

import (
	"path"
	"strings"

	"github.com/multiformats/go-varint"
)

var id0, _ = NewIDAddress(0)

const ROOT_STR = "/root"
const UNDEF_STR = "/"

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

// NewSubnetID generates the ID for a subnet from the networkName
// of its parent.
//
// It takes the parent name and adds the source address of the subnet
// actor that represents the subnet.
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

	s1 := strings.Split(str, "/")
	actor, err := NewFromString(s1[len(s1)-1])
	if err != nil {
		return UndefSubnetID, err
	}
	return SubnetID{
		Parent: strings.Join(s1[:len(s1)-1], "/"),
		Actor:  actor,
	}, nil
}

// GetParent returns the ID of the parent network.
func (id SubnetID) GetParent() (SubnetID, error) {
	if id.Parent == ROOT_STR {
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
	s1 := strings.Split(id.String(), "/")
	s2 := strings.Split(other.String(), "/")
	if len(s1) < len(s2) {
		s1, s2 = s2, s1
	}
	out := "/"
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
	s1 := strings.Split(id.String(), "/")
	s2 := strings.Split(curr.String(), "/")
	// curr needs to be contained in id
	if len(s2) >= len(s1) {
		return UndefSubnetID
	}
	_, l := id.CommonParent(curr)
	out := "/"
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
	s1 := strings.Split(id.String(), "/")
	s2 := strings.Split(curr.String(), "/")
	// curr needs to be contained in id
	if len(s2) > len(s1) {
		return UndefSubnetID
	}

	_, l := id.CommonParent(curr)
	out := "/"
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

// String returns the id in string form
func (id SubnetID) String() string {
	return strings.Join([]string{id.Parent, id.Actor.String()}, "/")
}

// returns useful payload from hierarchical address
func (a Address) hierarchical_payload() (string, error) {
	size, _, err := varint.FromUvarint(a.Bytes()[1:2])
	if err != nil {
		return "", err
	}
	// hierarchical addresses have a fixed size. We prefix a single byte
	// varint to know the total size of the address
	return a.str[2 : size+2], nil
}

// Subnet returns subnet information for an address if any.
func (a Address) Subnet() (SubnetID, error) {
	if a.str[0] != Hierarchical {
		return UndefSubnetID, ErrNotHierarchical
	}
	pl, err := a.hierarchical_payload()
	if err != nil {
		return UndefSubnetID, err
	}
	return SubnetIDFromString(strings.Split(pl, "::")[0])
}

// RawAddr return the address without subnet context information
func (a Address) RawAddr() (Address, error) {
	if a.str[0] != Hierarchical {
		return a, nil
	}
	pl, err := a.hierarchical_payload()
	if err != nil {
		return Undef, err
	}
	return NewFromString(strings.Split(pl, "::")[1])
}

func (a Address) PrettyPrint() string {
	if a.str[0] != Hierarchical {
		return a.String()
	}
	pl, err := a.hierarchical_payload()
	if err != nil {
		return UNDEF_STR
	}
	return string(pl)
}
