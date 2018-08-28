package ansible

import (
	"github.com/RSE-Cambridge/data-acc/internal/pkg/registry"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPlugin_GetInventory(t *testing.T) {
	volume := registry.Volume{Name: "1", UUID: "abcdefgh"}
	brickAllocations := []registry.BrickAllocation{
		{Hostname: "dac1", Device: "nvme1n1", AllocatedIndex: 0},
		{Hostname: "dac1", Device: "nvme2n1", AllocatedIndex: 1},
		{Hostname: "dac1", Device: "nvme3n1", AllocatedIndex: 2},
		{Hostname: "dac2", Device: "nvme2n1", AllocatedIndex: 3},
		{Hostname: "dac2", Device: "nvme3n1", AllocatedIndex: 4},
	}
	result := getInventory(volume, brickAllocations)
	expected := `dac-fake:
  children:
    abcdefgh:
      hosts:
        dac1:
          abcdefgh_mgs: nvme0n1
          abcdefgh_mdt: nvme1n1
          abcdefgh_osts: {nvme2n1: 1, nvme3n1: 2}
        dac2:
          abcdefgh_osts: {nvme2n1: 3, nvme3n1: 4}
      vars:
        abcdefgh_mgsnode: dac1
`
	assert.Equal(t, expected, result)
}

func TestPlugin_GetInventory_withNoOstOnOneHost(t *testing.T) {
	volume := registry.Volume{Name: "1", UUID: "abcdefgh"}
	brickAllocations := []registry.BrickAllocation{
		{Hostname: "dac1", Device: "nvme1n1", AllocatedIndex: 0},
		{Hostname: "dac2", Device: "nvme2n1", AllocatedIndex: 1},
	}
	result := getInventory(volume, brickAllocations)
	expected := `dac-fake:
  children:
    abcdefgh:
      hosts:
        dac1:
          abcdefgh_mgs: nvme0n1
          abcdefgh_mdt: nvme1n1
        dac2:
          abcdefgh_osts: {nvme2n1: 1}
      vars:
        abcdefgh_mgsnode: dac1
`
	assert.Equal(t, expected, result)
}

func TestPlugin_GetPlaybook(t *testing.T) {
	volume := registry.Volume{Name: "1", UUID: "abcdefgh"}
	result := getPlaybook(volume)
	assert.Equal(t, `---
- name: Install Lustre
  hosts: abcdefgh
  any_errors_fatal: true
  become: yes
  gather_facts: no
  roles:
    - role: lustre`, result)
}
