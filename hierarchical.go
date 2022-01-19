package address

import (
	"path"
	"strings"
)

// Root is the ID of the root network
const RootSubnet = SubnetID("/root")

// Undef is the undef ID
const UndefSubnetID = SubnetID("")

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
