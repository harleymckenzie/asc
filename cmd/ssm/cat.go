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

	param, err := svc.GetParameter(ctx, &ascTypes.GetParameterInput{
		Name:    args[0],
		Decrypt: catDecrypt,
	})
	if err != nil {
		return fmt.Errorf("get parameter: %w", err)
	}

	fmt.Println(aws.ToString(param.Value))
	return nil
}
