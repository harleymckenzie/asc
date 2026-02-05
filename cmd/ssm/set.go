package ssm

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var (
	setType        string
	setDescription string
	setOverwrite   bool
	setStdin       bool
)

func init() {
	newSetFlags(setCmd)
}

var setCmd = &cobra.Command{
	Use:   "set <parameter-name> [value]",
	Short: "Create or update an SSM parameter",
	Long: `Create or update an SSM parameter.

If the parameter already exists, you will be prompted for confirmation
unless --overwrite is specified.

For sensitive values, use --stdin to read from stdin (avoids shell history).`,
	GroupID: "actions",
	Args:    cobra.RangeArgs(1, 2),
	Example: `  asc ssm set /myapp/prod/db-host "localhost"
  asc ssm set /myapp/prod/db-pass "secret" --type SecureString
  asc ssm set /myapp/prod/config "value" --description "App config"
  asc ssm set /myapp/prod/secret --type SecureString --stdin
  echo "secret" | asc ssm set /myapp/prod/secret --type SecureString --stdin
  asc ssm set /myapp/prod/key "newvalue" --overwrite`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(SetSSMParameter(cmd, args))
	},
}

func newSetFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&setType, "type", "t", "String", "Parameter type: String, StringList, or SecureString")
	cmd.Flags().StringVarP(&setDescription, "description", "d", "", "Parameter description")
	cmd.Flags().BoolVarP(&setOverwrite, "overwrite", "o", false, "Overwrite existing parameter without confirmation")
	cmd.Flags().BoolVar(&setStdin, "stdin", false, "Read value from stdin")
}

func SetSSMParameter(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	paramName := args[0]

	// Get value from args or stdin
	var value string
	if setStdin {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("read from stdin: %w", err)
		}
		value = strings.TrimSuffix(string(data), "\n")
	} else if len(args) >= 2 {
		value = args[1]
	} else {
		return fmt.Errorf("value required: provide as argument or use --stdin")
	}

	// Validate type
	validTypes := map[string]bool{"String": true, "StringList": true, "SecureString": true}
	if !validTypes[setType] {
		return fmt.Errorf("invalid type %q: must be String, StringList, or SecureString", setType)
	}

	// Check if parameter already exists
	exists := false
	existingParam, err := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
		Name:    paramName,
		Decrypt: false,
	})
	if err == nil && existingParam != nil {
		exists = true
	}

	// Prompt for confirmation if exists and not overwriting
	if exists && !setOverwrite {
		fmt.Printf("Parameter %s already exists (type: %s).\n", paramName, existingParam.Type)
		fmt.Print("Overwrite? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read confirmation: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	// Put the parameter
	err = svc.PutParameter(ctx, &ascTypes.PutParameterInput{
		Name:        paramName,
		Value:       value,
		Type:        setType,
		Description: setDescription,
		Overwrite:   exists, // Only set overwrite if it exists
	})
	if err != nil {
		return fmt.Errorf("put parameter: %w", err)
	}

	if exists {
		fmt.Printf("Updated: %s\n", paramName)
	} else {
		fmt.Printf("Created: %s\n", paramName)
	}

	return nil
}
