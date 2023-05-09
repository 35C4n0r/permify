package cmd

import (
	"fmt"
	"net/url"

	"github.com/gookit/color"
	"github.com/spf13/cobra"

	"github.com/Permify/permify/pkg/cmd/flags"
	cov "github.com/Permify/permify/pkg/development/coverage"
	"github.com/Permify/permify/pkg/development/file"
)

// NewCoverageCommand - creates a new coverage command
func NewCoverageCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "coverage <file>",
		Short: "coverage analysis of authorization model and assertions",
		RunE:  coverage(),
		Args:  cobra.ExactArgs(1),
	}

	// register flags for validation
	flags.RegisterValidationFlags(command)

	return command
}

// coverage - coverage analysis of authorization model and assertions
func coverage() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// parse the url from the first argument
		u, err := url.Parse(args[0])
		if err != nil {
			return err
		}

		// create a new decoder from the url
		decoder, err := file.NewDecoderFromURL(u)
		if err != nil {
			return err
		}

		// create a new shape
		s := &file.Shape{}

		// decode the schema from the decoder
		err = decoder.Decode(s)
		if err != nil {
			return err
		}

		color.Notice.Println("initiating coverage analysis... 🚀")

		DisplayCoverageInfo(cov.Run(*s))

		return nil
	}
}

// DisplayCoverageInfo - Display the schema coverage information
func DisplayCoverageInfo(schemaCoverageInfo cov.SchemaCoverageInfo) {
	color.Notice.Println("schema coverage information:")

	for _, entityCoverageInfo := range schemaCoverageInfo.EntityCoverageInfo {
		color.Notice.Printf("entity: %s\n", entityCoverageInfo.EntityName)

		fmt.Printf("  uncovered relationships:\n")

		for _, value := range entityCoverageInfo.UncoveredRelationships {
			fmt.Printf("    - %v\n", value)
		}

		fmt.Printf("  uncovered assertions:\n")

		for key, value := range entityCoverageInfo.UncoveredAssertions {
			fmt.Printf("    %s:\n", key)
			for _, v := range value {
				fmt.Printf("    	%v\n", v)
			}
		}

		fmt.Printf("  coverage relationships percentage:")

		if entityCoverageInfo.CoverageRelationshipsPercent <= 50 {
			color.Danger.Printf(" %d%%\n", entityCoverageInfo.CoverageRelationshipsPercent)
		} else {
			color.Success.Printf(" %d%%\n", entityCoverageInfo.CoverageRelationshipsPercent)
		}

		fmt.Printf("  coverage assertions percentage: \n")

		for key, value := range entityCoverageInfo.CoverageAssertionsPercent {
			fmt.Printf("    %s:", key)
			if value <= 50 {
				color.Danger.Printf(" %d%%\n", value)
			} else {
				color.Success.Printf(" %d%%\n", value)
			}
		}
	}
}
