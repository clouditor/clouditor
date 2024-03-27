package discovery

import (
	"fmt"

	"clouditor.io/clouditor/v2/api/discovery"
	"clouditor.io/clouditor/v2/internal/config"
	"clouditor.io/clouditor/v2/internal/launcher"
	service_discovery "clouditor.io/clouditor/v2/service/discovery"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fmt.Println("init of command")
}

func NewDiscoveryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "discovery",
		Short: "Starts a server which contains the Clouditor Discovery Service",
		Long:  "The Clouditor Discovery service discovers things",
		RunE: func(cmd *cobra.Command, args []string) error {
			l, err := launcher.NewLauncher("Discovery",
				service_discovery.DefaultServiceSpec(),
			)
			if err != nil {
				return err
			}

			return l.Launch()
		},
	}

	cmd.Flags().String(config.AssessmentURLFlag, config.DefaultAssessmentURL, "Specifies the Assessment URL")
	cmd.Flags().String(config.CloudServiceIDFlag, discovery.DefaultCloudServiceID, "Specifies the Cloud Service ID")
	cmd.Flags().Bool(config.DiscoveryAutoStartFlag, config.DefaultDiscoveryAutoStart, "Automatically start the discovery when engine starts")
	cmd.Flags().StringSliceP(config.DiscoveryProviderFlag, "p", []string{}, "Providers to discover, separated by comma")
	cmd.Flags().String(config.DiscoveryResourceGroupFlag, config.DefaultDiscoveryResourceGroup, "Limit the scope of the discovery to a resource group (currently only used in the Azure discoverer")

	_ = viper.BindPFlag(config.AssessmentURLFlag, cmd.Flags().Lookup(config.AssessmentURLFlag))
	_ = viper.BindPFlag(config.CloudServiceIDFlag, cmd.Flags().Lookup(config.CloudServiceIDFlag))
	_ = viper.BindPFlag(config.DiscoveryAutoStartFlag, cmd.Flags().Lookup(config.DiscoveryAutoStartFlag))
	_ = viper.BindPFlag(config.DiscoveryProviderFlag, cmd.Flags().Lookup(config.DiscoveryProviderFlag))
	_ = viper.BindPFlag(config.DiscoveryResourceGroupFlag, cmd.Flags().Lookup(config.DiscoveryResourceGroupFlag))

	return cmd
}
