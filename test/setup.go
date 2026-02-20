package test

import (
	"github.com/infraboard/mcube/v2/ioc"
)

// 各个功能测试初始化配置文件路径
func SetUp() {
	ioc.DevelopmentSetupWithPath("C:\\Users\\Administrator\\Desktop\\code\\Go\\cmdb\\etc\\application.toml")
}
