package netapp_test

import (
	"reflect"
	"testing"

	"github.com/pepabo/go-netapp/netapp"
)

func TestAggregate_ListSuccess(t *testing.T) {
	c, teardown := createTestClientWithFixtures(t)
	defer teardown()

	tru := new(bool)
	*tru = true
	f := new(bool)
	*f = false
	options := &netapp.AggrOptions{
		MaxRecords: 20,
		Query: &netapp.AggrInfo{
			AggrRaidAttributes: &netapp.AggrRaidAttributes{
				IsRootAggregate: f,
			},
		},
	}
	call, _, err := c.Aggregate.List(options)
	checkResponseSuccess(&call.Results.ResultBase, err, t)

	info := call.Results.AggrAttributes[0]
	tests := []struct {
		name string
		got  interface{}
		want interface{}
	}{
		{"Interface Name", info.AggregateName, "aggr0_root_cluster01_01"},
		{"Cluster Name", info.AggrOwnershipAttributes.Cluster, "test-cluster01"},
		{"Node Name", info.AggrOwnershipAttributes.HomeName, "test-cluster01-01"},
		{"Root Aggregate", info.AggrRaidAttributes.IsRootAggregate, tru},
		{"Size", info.AggrSpaceAttributes.SizeTotal, 394454966272},
		{"Total Files", info.AggrInodeAttributes.FilesTotal, 31149},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !reflect.DeepEqual(tt.got, tt.want) {
				t.Errorf("Aggregate.GetList() got = %+v, want %+v", tt.got, tt.want)
			}
		})
	}
}

func TestAggregate_CreateSuccess(t *testing.T) {

	c, teardown := createTestClientWithFixtures(t)
	defer teardown()

	createOpts := &netapp.AggrCreateInfo{
		Aggreagate:     "test_2",
		ActualDiskType: "VMDISK",
		Nodes:          []string{"test-node-2"},
		DiskCount:      1,
	}
	call, _, err := c.Aggregate.Create(createOpts)

	checkResponseSuccess(&call.Results.SingleResultBase, err, t)
}

func TestAggregate_CreateFailure(t *testing.T) {

	c, teardown := createTestClientWithFixtures(t)
	defer teardown()

	createOpts := &netapp.AggrCreateInfo{
		Aggreagate:     "test_2",
		ActualDiskType: "VMDISK",
		Nodes:          []string{"test-node-2"},
		DiskCount:      1,
	}
	call, _, err := c.Aggregate.Create(createOpts)

	checkResponseFailure(&call.Results.SingleResultBase, err, t)

	testFailureResult(17, "An aggregate already uses test_2 as name", &call.Results, t)
}
