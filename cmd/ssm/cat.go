package ssm

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/harleymckenzie/asc/internal/service/ssm"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ssm/types"
	"github.com/harleymckenzie/asc/internal/shared/cmdutil"
	"github.com/spf13/cobra"
)

var catDecrypt bool

func init() {
	newCatFlags(catCmd)
}

var catCmd = &cobra.Command{
	Use:     "cat <parameter-name>",
	Short:   "Print the value of an SSM parameter",
	Hidden:  true,
	Args:    cobra.ExactArgs(1),
	Example: `  asc ssm cat /myapp/prod/config
  asc ssm cat /myapp/prod/secret --decrypt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmdutil.DefaultErrorHandler(CatSSMParameter(cmd, args))
	},
}

func newCatFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&catDecrypt, "decrypt", "d", false, "Decrypt SecureString values")
}

func CatSSMParameter(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	svc, err := cmdutil.CreateService(cmd, ssm.NewSSMService)
	if err != nil {
		return fmt.Errorf("create ssm service: %w", err)
	}

	names, err := resolveGlob(ctx, svc, args[0])
	if err != nil {
		return fmt.Errorf("resolve glob: %w", err)
	}

	if len(names) == 0 {
		return fmt.Errorf("no parameters matching: %s", args[0])
	}

	for _, name := range names {
		param, err := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
			Name:    name,
			Decrypt: catDecrypt,
		})
		if err != nil {
			return fmt.Errorf("get parameter %s: %w", name, err)
		}

		fmt.Println(aws.ToString(param.Value))
	}
	return nil
}
