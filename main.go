package main

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
)

func KeyMatch(r string, p string) bool {
	return r == p
}

func KeyMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return (bool)(KeyMatch(name1, name2)), nil
}

func main() {
	a, _ := gormadapter.NewAdapter("mysql", "root:123456@tcp(127.0.0.1:3306)/test", true) // Your driver and data source.
	e, _ := casbin.NewEnforcer("./model.conf", a)

	e.AddFunction("my_func", KeyMatchFunc)

	// 加载策略
	e.LoadPolicy()

	// 添加策略，自动在数据库中添加
	// e.AddPolicy("alice", "data2", "read")
	// 根据指定的 index(策略的 v0, v1...) 获取它所有策略
	// filterdPolicy := e.GetFilteredPolicy(0, "alice")
	// fmt.Println("policy：", filterdPolicy)
	// 修改策略
	// e.UpdatePolicy([]string{"alice", "data2", "read"}, []string{"alice", "data666", "read"})

	// 获取用户拥有的角色
	// roles, _ := e.GetRolesForUser("alice")
	// fmt.Println("roles: ", roles) // [data1_admin]

	// 获取用户拥有的隐式权限（也就是根据角色来的权限、和自身的权限）
	// permission, _ := e.GetImplicitPermissionsForUser("alice")
	// fmt.Println("permission: ", permission) // [[data1_admin data1 read] [data1_admin data666 read]]

	// 添加策略组
	// e.AddGroupingPolicy("alice", "data1_admin")

	// 测试 request
	sub := "alice" // the user that wants to access a resource.
	obj := "data1" // the resource that is going to be accessed.
	act := "read"  // the operation that the user performs on the resource.
	ok, err := e.Enforce(sub, obj, act)

	if err != nil {
		fmt.Println("err: ", err)
	}

	if ok {
		fmt.Println("通过")
	} else {
		fmt.Println("未通过")
	}
}
