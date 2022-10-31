package ignitecmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/services/network"
)

// NewNetworkRequestParamChange creates a new command to send change param request
func NewNetworkRequestParamChange() *cobra.Command {
	c := &cobra.Command{
		Use:   "param [launch-id] [module-name] [param-name] [value (json, string, number) as bytes]",
		Short: "Send request to change param",
		RunE:  networkRequestParamChange,
		Args:  cobra.ExactArgs(4),
	}

	flagSetClearCache(c)
	c.Flags().AddFlagSet(flagNetworkFrom())
	c.Flags().AddFlagSet(flagSetHome())
	c.Flags().AddFlagSet(flagSetKeyringBackend())
	c.Flags().AddFlagSet(flagSetKeyringDir())
	return c
}

func networkRequestParamChange(cmd *cobra.Command, args []string) error {
	session := cliui.New(cliui.StartSpinner())
	defer session.End()

	nb, err := newNetworkBuilder(cmd, CollectEvents(session.EventBus()))
	if err != nil {
		return err
	}

	// parse launch ID
	launchID, err := network.ParseID(args[0])
	if err != nil {
		return err
	}

	// parse module name
	moduleName := args[1]

	// parse param name
	paramName := args[2]

	// parse value
	value := []byte(args[3])

	n, err := nb.Network()
	if err != nil {
		return err
	}

	return n.SendParamChangeRequest(
		cmd.Context(),
		launchID,
		moduleName,
		paramName,
		value,
	)
}
