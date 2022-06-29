package address_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/go-address"
)

func TestNaming(t *testing.T) {
	address.CurrentNetwork = address.Mainnet
	addr1, err := address.NewIDAddress(101)
	require.NoError(t, err)
	addr2, err := address.NewIDAddress(102)
	require.NoError(t, err)
	root := address.RootSubnet
	net1 := address.NewSubnetID(root, addr1)
	net2 := address.NewSubnetID(net1, addr2)

	t.Log("Test actors")
	actor1 := net1.GetActor()
	require.Equal(t, actor1, addr1)
	actor2 := net2.GetActor()
	require.NoError(t, err)
	require.Equal(t, actor2, addr2)
	actorRoot := root.GetActor()
	require.NoError(t, err)
	require.Equal(t, address.Undef, actorRoot)

	t.Log("Test parents")
	parent1, err := net1.GetParent()
	require.NoError(t, err)
	require.Equal(t, root, parent1)
	parent2, err := net2.GetParent()
	require.NoError(t, err)
	require.Equal(t, parent2, net1)
	parentRoot, err := root.GetParent()
	require.NoError(t, err)
	require.Equal(t, parentRoot, address.UndefSubnetID)
}

func TestCborMarshal(t *testing.T) {
	addr1, err := address.NewIDAddress(101)
	require.NoError(t, err)
	root := address.RootSubnet
	net1 := address.NewSubnetID(root, addr1)

	var buf bytes.Buffer
	err = net1.MarshalCBOR(&buf)
	require.NoError(t, err)
	net2 := address.SubnetID{}
	err = net2.UnmarshalCBOR(&buf)
	require.NoError(t, err)
	require.Equal(t, net1, net2)
}

func TestHAddress(t *testing.T) {
	address.CurrentNetwork = address.Mainnet
	id, _ := address.NewIDAddress(1000)
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
	require.Equal(t, a.PrettyPrint(), "/root:f01000")

}

func TestRustInterop(t *testing.T) {
	// This string address was generated from the Rust implementation.
	_, err := address.NewFromString("f4bqys64tpn52c6zrqgeydamidvvmn62lofvhjd2ugzca6sof2j2ubwok6cj4xxbfzz4yuxfkgobpihhd2thlanmsh3w2ptld2gqkn2aoph33q")
	require.NoError(t, err)
}

func TestSubnetOps(t *testing.T) {
	address.CurrentNetwork = address.Mainnet
	testParentAndBottomUp(t, "/root/f01", "/root/f01/f02", "/root/f01", 2, false)
	testParentAndBottomUp(t, "/root/f03/f01", "/root/f01/f02", "/root", 1, true)
	testParentAndBottomUp(t, "/root/f03/f01/f04", "/root/f03/f01/f05", "/root/f03/f01", 3, true)
	testParentAndBottomUp(t, "/root/f03/f01", "/root/f03/f02", "/root/f03", 2, true)

	testDownOrUp(t, "/root/f01/f02/f03", "/root/f01", "/root/f01/f02", true)
	testDownOrUp(t, "/root/f01/f02/f03", "/root/f01/f02", "/root/f01/f02/f03", true)
	testDownOrUp(t, "/root/f02", "/root/f01/f02/f03", address.UndefSubnetID.String(), true)
	testDownOrUp(t, "/root/f02", "/root/f02", address.UndefSubnetID.String(), true)

	testDownOrUp(t, "/root/f01/f02/f03", "/root/f01", "/root", false)
	testDownOrUp(t, "/root", "/root/f01", address.UndefSubnetID.String(), false)
	testDownOrUp(t, "/root/f01/f02/f03", "/root/f01/f02/f03/d", address.UndefSubnetID.String(), false)
	testDownOrUp(t, "/root/f01/f02/f03", "/root/f01/f02", "/root/f01", false)
}

func testDownOrUp(t *testing.T, from, to, expected string, down bool) {
	sn, _ := address.SubnetIDFromString(from)
	arg, _ := address.SubnetIDFromString(to)
	ex, _ := address.SubnetIDFromString(expected)
	if down {
		require.Equal(t, sn.Down(arg), ex)
	} else {
		require.Equal(t, sn.Up(arg), ex)
	}
}

func testParentAndBottomUp(t *testing.T, from, to, parent string, exl int, bottomup bool) {
	sfrom, err := address.SubnetIDFromString(from)
	require.NoError(t, err)
	sto, err := address.SubnetIDFromString(to)
	require.NoError(t, err)
	p, l := sfrom.CommonParent(sto)
	sparent, err := address.SubnetIDFromString(parent)
	require.NoError(t, err)
	require.Equal(t, p, sparent)
	require.Equal(t, exl, l)
}
