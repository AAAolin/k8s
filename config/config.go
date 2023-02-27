package config

import "time"

const (
	ListenAddr = "0.0.0.0:9091"
	KubeConfig = "D:\\GoProject\\k8s\\config\\config"
	LogsLimit  = 2000

	Username = "admin"
	Password = "123456"

	// 数据库配置
	DBUser     = "root"
	DBPassword = "123456"
	DBType     = "mysql"
	DBIP       = "192.168.1.16"
	DBPort     = 3306
	DBName     = "k8s"
	// 连接池配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
)

type Student struct {
	PersonValue Person
	Metha       meth
	age         int
}
type Person struct {
	name string
}

type meth interface {
	get()
	post()
}

func test() {
	var s Student
	s.PersonValue.name = "aaa"
	s.Metha.get()
}
