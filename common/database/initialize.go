package database

import (
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	toolsConfig "github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	toolsDB "github.com/go-admin-team/go-admin-core/tools/database"
	. "github.com/go-admin-team/go-admin-core/tools/gorm/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/global"
)

// Setup 配置数据库
func Setup() {
	for k := range toolsConfig.DatabasesConfig {
		SetupSimpleDatabase(k, toolsConfig.DatabasesConfig[k], false)
	}

	//租户模式，不同的公司，不同数据库，从服务器获取配置
	dynamicDBInit()
}

func SetupSimpleDatabase(host string, c *toolsConfig.Database, isFromTenant bool) {
	if global.Driver == "" {
		global.Driver = c.Driver
	}
	log.Infof("%s => %s", host, pkg.Green(c.Source))
	registers := make([]toolsDB.ResolverConfigure, len(c.Registers))
	for i := range c.Registers {
		registers[i] = toolsDB.NewResolverConfigure(
			c.Registers[i].Sources,
			c.Registers[i].Replicas,
			c.Registers[i].Policy,
			c.Registers[i].Tables)
	}
	resolverConfig := toolsDB.NewConfigure(c.Source, c.MaxIdleConns, c.MaxOpenConns, c.ConnMaxIdleTime, c.ConnMaxLifeTime, registers)
	db, err := resolverConfig.Init(&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: New(
			logger.Config{
				SlowThreshold: time.Second,
				Colorful:      true,
				LogLevel: logger.LogLevel(
					log.DefaultLogger.Options().Level.LevelForGorm()),
			},
		),
	}, opens[c.Driver])

	if err != nil {
		if !isFromTenant {
			log.Fatal(pkg.Red(c.Driver+" connect error :"), err)
		} else {
			log.Info(pkg.Red("租户 "+host+" MYSQL注册失败！！！"+c.Driver+" connect error :"), err)
			//应该发一个消息给管理员
			return
		}
	} else {
		if !isFromTenant {
			log.Info(pkg.Green(c.Driver + " connect success !"))
		} else {
			log.Info(pkg.Green("租户 " + host + " MYSQL注册成功!" + c.Driver + " connect success !"))
		}

	}

	//e := mycasbin.Setup(db, host+"_admin")

	sdk.Runtime.SetDb(host, db)
	//sdk.Runtime.SetCasbin(host, e)
}
