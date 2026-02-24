package impl

import (
	"cmdb/apps/resource"
	"cmdb/apps/secret"
	"context"

	"github.com/infraboard/mcube/v2/desense"
	"github.com/infraboard/mcube/v2/ioc/config/cache"
	"github.com/infraboard/mcube/v2/types"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s SecretServiceImpl) CreateSecret(ctx context.Context, request *secret.CreateSecretRequest) (*secret.Secret, error) {
	ins := secret.NewSecret(request)
	if err := ins.EncryptedApiSecret(); err != nil {
		return nil, err
	}
	// 使能upsert等同于gorm的save
	_, err := s.coll.UpdateOne(ctx, bson.M{"id": ins.Id}, bson.M{"$set": ins}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (s SecretServiceImpl) QuerySecret(ctx context.Context, request *secret.QuerySecretRequest) (*types.Set[*secret.Secret], error) {
	ins := secret.NewSecretSet()
	filter := bson.M{}
	cursor, err := s.coll.Find(ctx, filter, options.Find().SetLimit(int64(request.PageSize)).SetSkip(request.Offset))
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		// 构建载体实例
		se := secret.NewSecret(secret.NewCreateSecretRequest())
		if err := cursor.Decode(se); err != nil {
			return nil, err
		}
		ins.Items = append(ins.Items, se)
		ins.Total += int64(len(ins.Items))
	}
	ins.Total -= 1
	// 调用mcube封装的脱敏逻辑实现
	if err := desense.MaskStruct(ins); err != nil {
		return nil, err
	}
	return ins, nil
}

// 增加cache缓存实现
func (s SecretServiceImpl) DescribeSecret(ctx context.Context, request *secret.DescribeSecretRequest) (*secret.Secret, error) {
	se := secret.NewSecret(&secret.CreateSecretRequest{})

	// mcube封装的缓存Getter逻辑：如果能从缓存查到，就直接获取，如果没有就新建一个再保存到缓存中
	// 1. 从缓存中去(内存， 公共的内存服务 Redis)
	// 2. 能获取到，直接返回
	// 3. 不能获取, 选好从本地获取，返回，再把他设置到缓存中去
	// 4. 怎么实现: redis redis get(key)/set(key), obj -> JSON
	// 5. https://github.com/redis/go-redis  get, set
	// CacheGetter --> go-redis --> ObjectFinder
	if err := cache.NewGetter(ctx, func(ctx context.Context, objectId string) (any, error) {
		return s.describeSecret(ctx, request)
	}).
		Get(request.Id, se); err != nil {
		return nil, err
	}

	return se, nil
}

// 实际的describe逻辑
func (s SecretServiceImpl) describeSecret(ctx context.Context, request *secret.DescribeSecretRequest) (*secret.Secret, error) {
	// 初始化实例
	se := secret.NewSecret(&secret.CreateSecretRequest{})
	if err := s.coll.FindOne(ctx, bson.M{"id": request.Id}).Decode(se); err != nil {
		return nil, err
	}
	// 设置秘钥已经被加密
	se.SetIsEncrypted(true)
	// 解密
	if err := se.DecryptedApiSecret(); err != nil {
		return nil, err
	}
	return se, nil
}

func (s SecretServiceImpl) SyncResource(ctx context.Context, request *secret.SyncResourceRequest, handleFunc secret.SyncResourceHandleFunc) error {
	se, _ := s.DescribeSecret(ctx, secret.NewDescribeSecretRequest(request.Id))
	return se.Sync(func(in secret.ResourceResponse) {
		in.Resource.Meta.Namespace = "default"
		in.Resource.Meta.Domain = "default"

		// 调用resource模块的save实现保存
		res, err := resource.GetService().Save(ctx, in.Resource)
		if err != nil {
			in.Success = false
			in.Message = "保存失败"
		} else {
			in.Success = true
			in.Resource = res
		}
		// 传参函数调用
		handleFunc(in)
	})
}
