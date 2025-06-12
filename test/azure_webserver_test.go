package test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestAzureLinuxVMCreation(t *testing.T) {
	t.Parallel()

	subscriptionID := "a61f90e1-d9c8-4c9c-9cd5-0a904a50cf48"
	labelPrefix := "khal0233"
	resourceGroupName := fmt.Sprintf("%s-A05-RG", labelPrefix)
	vmName := fmt.Sprintf("%sA05VM", labelPrefix)
	nicName := fmt.Sprintf("%sA05Nic", labelPrefix)

	terraformOptions := &terraform.Options{
		TerraformDir: "../",
		Vars: map[string]interface{}{
			"labelPrefix": labelPrefix,
		},
	}

	// Deploy infrastructure
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Wait for Azure to register the resources
	time.Sleep(20 * time.Second)

	// ✅ Test 1: Check if VM exists
	vmExists := azure.VirtualMachineExists(t, vmName, resourceGroupName, subscriptionID)
	assert.True(t, vmExists, "Expected VM to exist")

	// ✅ Test 2: Check if NIC exists
	time.Sleep(15 * time.Second)
	nicExists := azure.NetworkInterfaceExists(t, nicName, resourceGroupName, subscriptionID)
	assert.True(t, nicExists, "Expected NIC to exist")

	// ✅ Test 3: Check the OS image info
	vm := azure.GetVirtualMachine(t, vmName, resourceGroupName, subscriptionID)
	image := vm.StorageProfile.ImageReference

	assert.Contains(t, strings.ToLower(*image.Publisher), "canonical", "Expected Canonical as the OS publisher")
	assert.Contains(t, strings.ToLower(*image.Offer), "ubuntu", "Expected Ubuntu as the OS")
}
