package address

import (
	"path"
	"strings"
)

// Root is the ID of the root network
const RootSubnet = SubnetID("/root")

// Undef is the undef ID
const UndefSubnetID = SubnetID("/")

// SubnetID represents the ID of a subnet
type SubnetID string

// NewSubnetID generates the ID for a subnet from the networkName
// of its parent.
//
// It takes the parent name and adds the source address of the subnet
// actor that represents the subnet.
func NewSubnetID(parentName SubnetID, SubnetActorAddr Address) SubnetID {
	return SubnetID(path.Join(parentName.String(), SubnetActorAddr.String()))
}

// Parent returns the ID of the parent network.
func (id SubnetID) Parent() SubnetID {
	if id == RootSubnet {
		return UndefSubnetID
	}
	return SubnetID(path.Dir(string(id)))
}

// Actor returns the subnet actor for a subnet
//
// Returns the address of the actor that handles the logic for a subnet
// in its parent Subnet.
func (id SubnetID) Actor() (Address, error) {
	if id == RootSubnet {
		return Undef, nil
	}
	_, saddr := path.Split(string(id))
	return NewFromString(saddr)
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
			return SubnetID(out), l
		}
	}
	return SubnetID(out), l
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
	return SubnetID(out)
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
	return SubnetID(out)
}

// String returns the id in string form
func (id SubnetID) String() string {
	return string(id)
}

// Subnet returns subnet information for an address if any.
func (a Address) Subnet() (SubnetID, error) {
	if a.str[0] != Hierarchical {
		return UndefSubnetID, ErrNotHierarchical
	}
	return SubnetID(strings.Split(string(a.str[1:]), "::")[0]), nil
}

// RawAddr return the address without subnet context information
func (a Address) RawAddr() (Address, error) {
	if a.str[0] != Hierarchical {
		return a, nil
	}
	return NewFromString(strings.Split(string(a.str[1:]), "::")[1])
}

func (a Address) PrettyPrint() string {
	if a.str[0] != Hierarchical {
		return a.String()
	}
	return string(string(a.str[1:]))
}
