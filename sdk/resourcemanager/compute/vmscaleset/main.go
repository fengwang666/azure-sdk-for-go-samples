// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"
)

var (
	subscriptionID     string
	location           = "westus2"
	resourceGroupName  = "fengwang-dev"
	vmScaleSetName     = "reimage-test"
	instanceId         = "0"
)

func main() {
	subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	if len(subscriptionID) == 0 {
		log.Fatal("AZURE_SUBSCRIPTION_ID is not set.")
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	var totalLatency int64
	for i := 0; i < 100; i++ {
		log.Printf("Run #%d", i)
		start := time.Now()
		err = reimageVM(ctx, cred)
		if err != nil {
			log.Printf("The average latency of reimage is %d ms", totalLatency/int64(i))
			log.Fatal(err)
		}
		elapsed := time.Since(start).Milliseconds()
		log.Printf("Spent %d ms", elapsed)
		totalLatency += elapsed
	}
	log.Printf("The average latency of reimage is %d ms", totalLatency/100)
}

func reimageVM(ctx context.Context, cred azcore.TokenCredential) error {
	vmssClient, err := armcompute.NewVirtualMachineScaleSetsClient(subscriptionID, cred, nil)
	if err != nil {
		return err
	}
	pollerResp, err := vmssClient.BeginReimage(ctx, resourceGroupName, vmScaleSetName,
		&armcompute.VirtualMachineScaleSetsClientBeginReimageOptions{VMScaleSetReimageInput: &armcompute.VirtualMachineScaleSetReimageParameters{
			InstanceIDs: []*string{
				to.Ptr(instanceId)},
		},
		})
	if err != nil {
		return err
	}

	_, err = pollerResp.PollUntilDone(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
