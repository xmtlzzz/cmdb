package impl

import (
	"cmdb/apps/resource"
	"context"

	"github.com/infraboard/mcube/v2/exception"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (s *ResourceServiceImpl) Search(ctx context.Context, r *resource.SearchRequestSet) (*resource.ResourceSet, error) {
	set := resource.NewResourceSet()
	filter := bson.M{}

	// regex正则匹配name的内容
	if r.Keywords != "" {
		filter["name"] = bson.M{"$regex": r.Keywords, "$options": "im"}
	}
	// 匹配type字段内容
	if r.Type != nil {
		filter["type"] = *r.Type
	}
	// 多个tags遍历
	if r.Tags != nil {
		for k, v := range r.Tags {
			filter[k] = v
		}
	}
	// 查找逻辑
	resp, err := s.coll.Find(ctx, filter, options.Find().SetLimit(r.PageSize).SetSkip(r.SetSkip()))
	if err != nil {
		return nil, err
	}
	for resp.Next(ctx) {
		res := &resource.Resource{}
		if err := resp.Decode(res); err != nil {
			return nil, err
		}
		set.Items = append(set.Items, res)
	}
	return set, nil
}

// 将数据保存到mongodb
func (s *ResourceServiceImpl) Save(ctx context.Context, r *resource.Resource) (*resource.Resource, error) {
	if err := r.Validate(); err != nil {
		return nil, err
	}
	_, err := s.coll.InsertOne(ctx, r)
	if err != nil {
		exception.NewBadRequest("写入错误，resource.go 26")
	}
	return r, nil
}

func (s *ResourceServiceImpl) DeleteResource(context.Context, *resource.DeleteResourceRequest) error {
	return nil
}
