package address

import (
	"fmt"
	"path"
	"strings"
)

var id0, _ = NewIDAddress(0)

const (
	ROOT_STR          = "/root"
	SUBNET_SEPARATOR  = "/"
	UNDEF_STR         = SUBNET_SEPARATOR
	HC_ADDR_SEPARATOR = ":"
	HC_ADDR_END       = byte(',')
	HA_ROOT_LEN       = 5
	HA_LEVEL_LEN      = 23
	RAW_ADDR_LEN      = 66
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

// String returns the id in string form
func (id SubnetID) String() string {
	return strings.Join([]string{id.Parent, id.Actor.String()}, SUBNET_SEPARATOR)
}

// Levels returns the number of levels in the current subnetID
func (id SubnetID) Levels() int {
	return len(strings.Split(id.String(), SUBNET_SEPARATOR)) - 1

}

// returns useful payload from hierarchical address
func (a Address) parse_hierarchical() ([]string, error) {
	str := string(a.Payload())
	str = strings.Split(str, string(HC_ADDR_END))[0]
	out := strings.Split(str, HC_ADDR_SEPARATOR)
	if len(out) != 2 {
		return nil, fmt.Errorf("error parsing hierarchical address")
	}
	return out, nil
}

// Subnet returns subnet information for an address if any.
func (a Address) Subnet() (SubnetID, error) {
	if a.str[0] != Hierarchical {
		return UndefSubnetID, ErrNotHierarchical
	}
	pl, err := a.parse_hierarchical()
	if err != nil {
		return UndefSubnetID, err
	}
	return SubnetIDFromString(pl[0])
}

// RawAddr return the address without subnet context information
func (a Address) RawAddr() (Address, error) {
	if a.str[0] != Hierarchical {
		return a, nil
	}
	pl, err := a.parse_hierarchical()
	if err != nil {
		return Undef, err
	}
	return decode_raw_str(pl[1])
}

func (a Address) PrettyPrint() string {
	if a.str[0] != Hierarchical {
		return a.String()
	}
	pl, err := a.parse_hierarchical()
	if err != nil {
		return UNDEF_STR
	}
	raw, err := decode_raw_str(pl[1])
	if err != nil {
		return UNDEF_STR
	}
	return string(pl[0] + HC_ADDR_SEPARATOR + raw.String())
}
