package main

import (
	"context"
	"log"

	"github.com/coinbase/waas-client-library-go/auth"
	"github.com/coinbase/waas-client-library-go/clients"
	v1clients "github.com/coinbase/waas-client-library-go/clients/v1"
	blockchain "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/blockchain/v1"
	pools "github.com/coinbase/waas-client-library-go/gen/go/coinbase/cloud/pools/v1"
	"google.golang.org/api/iterator"
)

const (
	// apiKeyName is the name of the API Key to use. Fill this out before running the main function.
	apiKeyName = "organizations/bc2ba454-5651-4559-afe9-5e70665bb5e6/apiKeys/da45e081-4c6d-4fe1-8baf-c0bac338546b"

	// privKeyTemplate is the private key of the API Key to use. Fill this out before running the main function.
	privKeyTemplate = "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEILvf0lunBJ8lG+m68ZwVDCgmB3BszddokTKOV6JNF0ehoAoGCCqGSM49\nAwEHoUQDQgAE/Ej2J53GbwwUATZ0QT61HHTnNP9o1NGO7LSweH2rcpHZbiTGTg1k\n4MB0Z4tOjrdesoFUbijsIRCaMFwbi7KSnQ==\n-----END EC PRIVATE KEY-----\n"
)

// An example function to demonstrate how to use the WaaS client libraries.
func main() {
	ctx := context.Background()

	authOpt := clients.WithAPIKey(&auth.APIKey{
		Name:       apiKeyName,
		PrivateKey: privKeyTemplate,
	})

	poolsClient, err := v1clients.NewPoolServiceClient(ctx, authOpt)

	if err != nil {
		log.Fatalf("error instantiating pools client: %v", err)
	}

	log.Printf("creating pool...")

	pool, err := poolsClient.CreatePool(ctx, &pools.CreatePoolRequest{
		Pool: &pools.Pool{
			DisplayName: "My First Pool",
		},
	})

	if err != nil {
		log.Fatalf("error creating pool: %v", err)
	}

	log.Printf("created pool: %v", pool)

	blockchainClient, err := v1clients.NewBlockchainServiceClient(ctx, authOpt)

	if err != nil {
		log.Fatalf("error instantiating blockchain client: %v", err)
	}

	log.Printf("listing networks...")

	networkIter := blockchainClient.ListNetworks(ctx, &blockchain.ListNetworksRequest{})

	for network, err := networkIter.Next(); err == nil; network, err = networkIter.Next() {
		log.Printf("got network: %v", network)
	}

	if err != nil && err != iterator.Done {
		log.Fatalf("error listing networks: %v", err)
	}

	log.Printf("listing first 5 assets on Ethereum Goerli...")

	assetIter := blockchainClient.ListAssets(ctx, &blockchain.ListAssetsRequest{
		Parent: "networks/ethereum-goerli",
	})

	assetCount := 0

	for asset, err := assetIter.Next(); err == nil && assetCount < 5; asset, err = assetIter.Next() {
		log.Printf("got asset: %v", asset)

		assetCount++
	}

	if err != nil && err != iterator.Done {
		log.Fatalf("error listing assets: %v", err)
	}
}
