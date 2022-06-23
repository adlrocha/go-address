package address_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-address"
)

func TestNaming(t *testing.T) {
	addr1, err := address.NewIDAddress(101)
	require.NoError(t, err)
	addr2, err := address.NewIDAddress(102)
	require.NoError(t, err)
	root := address.RootSubnet
	net1 := address.NewSubnetID(root, addr1)
	net2 := address.NewSubnetID(net1, addr2)

	t.Log("Test actors")
	actor1, err := net1.Actor()
	require.NoError(t, err)
	require.Equal(t, addr1, actor1)
	actor2, err := net2.Actor()
	require.NoError(t, err)
	require.Equal(t, addr2, actor2)
	actorRoot, err := root.Actor()
	require.NoError(t, err)
	require.Equal(t, address.Undef, actorRoot)

	t.Log("Test parents")
	parent1 := net1.Parent()
	require.Equal(t, parent1, root)
	parent2 := net2.Parent()
	require.Equal(t, net1, parent2)
	parentRoot := root.Parent()
	require.Equal(t, address.UndefSubnetID, parentRoot)

}

func TestHAddress(t *testing.T) {
	_, err := address.NewFromBytes([]byte{address.Hierarchical})
	require.Error(t, err)

	_, err = address.NewFromString(string(address.Hierarchical))
	require.Error(t, err)

	id, err := address.NewIDAddress(1000)
	require.NoError(t, err)

	a, err := address.NewHCAddress(address.RootSubnet, id)
	require.NoError(t, err)

	sn, err := a.Subnet()
	require.NoError(t, err)
	require.Equal(t, address.RootSubnet, sn)

	raw, err := a.RawAddr()
	require.NoError(t, err)
	require.Equal(t, id, raw)

	_, err = id.Subnet()
	require.Error(t, err, address.ErrNotHierarchical)
	require.Equal(t, fmt.Sprintf("%v%v%v", address.RootSubnet, address.HCAddressSeparator, id), a.PrettyPrint())
}

func TestSubnetOps(t *testing.T) {
	testParentAndBottomUp(t, "/root/a", "/root/a/b", "/root/a", 2, false)
	testParentAndBottomUp(t, "/root/c/a", "/root/a/b", "/root", 1, true)
	testParentAndBottomUp(t, "/root/c/a/d", "/root/c/a/e", "/root/c/a", 3, true)
	testParentAndBottomUp(t, "/root/c/a", "/root/c/b", "/root/c", 2, true)

	require.Equal(t, address.SubnetID("/root/a/b/c").Down("/root/a"), address.SubnetID("/root/a/b"))
	require.Equal(t, address.SubnetID("/root/a/b/c").Down("/root/a/b"), address.SubnetID("/root/a/b/c"))
	require.Equal(t, address.SubnetID("/root/a").Down("/root/a/b/c"), address.UndefSubnetID)
	require.Equal(t, address.SubnetID("/root/b").Down("/root/a/b/c"), address.UndefSubnetID)
	require.Equal(t, address.SubnetID("/root/b").Down("/root/b"), address.UndefSubnetID)

	require.Equal(t, address.SubnetID("/root/a/b/c").Up("/root/a"), address.SubnetID("/root"))
	require.Equal(t, address.SubnetID("/root").Up("/root/a"), address.UndefSubnetID)
	require.Equal(t, address.SubnetID("/root/a/b/c").Up("/root/a/b/c/d"), address.UndefSubnetID)
	require.Equal(t, address.SubnetID("/root/a/b/c").Up("/root/a/b"), address.SubnetID("/root/a"))
}

func testParentAndBottomUp(t *testing.T, from, to, parent string, exl int, bottomup bool) {
	p, l := address.SubnetID(from).CommonParent(address.SubnetID(to))
	require.Equal(t, p, address.SubnetID(parent))
	require.Equal(t, exl, l)
}
