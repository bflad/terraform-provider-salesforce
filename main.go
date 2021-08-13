package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	tf5server "github.com/hashicorp/terraform-plugin-go/tfprotov5/server"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-provider-salesforce/internal/provider"
	"github.com/hashicorp/terraform-provider-salesforce/internal/providerdynamic"
	"github.com/hashicorp/terraform-provider-salesforce/internal/providerframework"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "dev"

	// goreleaser can also pass the specific commit if you want
	// commit  string = ""
)

func main() {
	if os.Getenv("LEGACY_MODE") != "" {
		var debugMode bool
		flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
		flag.Parse()

		opts := &plugin.ServeOpts{ProviderFunc: provider.New}

		if debugMode {
			err := plugin.Debug(context.Background(), "registry.terraform.io/hashicorp/salesforce", opts)
			if err != nil {
				log.Fatal(err.Error())
			}
			return
		}

		// rudimentary experimental mode toggle
		// TODO: use plugin-mux
		if os.Getenv("SALESFORCE_DYNAMIC_MODE") != "" {
			tf5server.Serve("registry.terraform.io/hashicorp/salesforce", providerdynamic.New)
		} else {
			plugin.Serve(opts)
		}
	} else {
		ctx := context.Background()

		tfsdk.Serve(ctx, providerframework.New, tfsdk.ServeOpts{
			Name: "registry.terraform.io/hashicorp/salesforce",
		})
	}
}
