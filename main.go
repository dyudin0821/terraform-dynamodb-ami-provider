package main

import (
	"context"
	"flag"
	"log"

	"github.com/dyudin0821/terraform-provider-dynamodb-ami/amidynamodb"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name amidynamodb

var (
	version string = "0.1.0"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/dyudin0821/terraform-provider-amidynamodb",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), amidynamodb.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}

}
