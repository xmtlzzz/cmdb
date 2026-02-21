package impl_test

import (
	"cmdb/apps/resource"
	_ "cmdb/apps/resource/impl"
	"cmdb/test"
	"testing"
)

func init() {
	test.SetUp()
}

func TestResourceServiceImpl_Save(t *testing.T) {
	resp, err := svc.Save(ctx, &resource.Resource{
		Spec: &resource.Spec{Vendor: 1, Name: "test"},
		Meta: &resource.Meta{Id: "1", Domain: "SH", Namespace: "test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestResourceServiceImpl_Search(t *testing.T) {
	sr := resource.NewSearchRequestSet()
	//sr.Keywords = "test"
	resp, err := svc.Search(ctx, sr)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
