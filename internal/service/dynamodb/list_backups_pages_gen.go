// Code generated by "internal/generate/listpages/main.go -ListOps=ListBackups -InputPaginator=ExclusiveStartBackupArn -OutputPaginator=LastEvaluatedBackupArn -- list_backups_pages_gen.go"; DO NOT EDIT.

package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func listBackupsPages(ctx context.Context, conn *dynamodb.Client, input *dynamodb.ListBackupsInput, fn func(*dynamodb.ListBackupsOutput, bool) bool) error {
	for {
		output, err := conn.ListBackups(ctx, input)
		if err != nil {
			return err
		}

		lastPage := aws.ToString(output.LastEvaluatedBackupArn) == ""
		if !fn(output, lastPage) || lastPage {
			break
		}

		input.ExclusiveStartBackupArn = output.LastEvaluatedBackupArn
	}
	return nil
}
