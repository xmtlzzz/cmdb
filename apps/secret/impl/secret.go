package impl

import (
	"cmdb/apps/secret"
	"context"

	"github.com/infraboard/mcube/v2/types"
)

func (s SecretServiceImpl) CreateSecret(ctx context.Context, request *secret.CreateSecretRequest) (*secret.Secret, error) {
	//TODO implement me
	panic("implement me")
}

func (s SecretServiceImpl) QuerySecret(ctx context.Context, request *secret.QuerySecretRequest) (*types.Set[*secret.Secret], error) {
	//TODO implement me
	panic("implement me")
}

func (s SecretServiceImpl) DescribeSecret(ctx context.Context, request *secret.DescribeSecretRequest) (*secret.Secret, error) {
	//TODO implement me
	panic("implement me")
}

func (s SecretServiceImpl) SyncResource(ctx context.Context, request *secret.SyncResourceRequest, handleFunc secret.SyncResourceHandleFunc) error {
	//TODO implement me
	panic("implement me")
}
