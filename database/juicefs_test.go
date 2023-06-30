package database

import (
	"fmt"
	"testing"

	"github.com/gobackup/gobackup/config"
	"github.com/longbridgeapp/assert"
	"github.com/spf13/viper"
)

func TestJuicefs_init(t *testing.T) {
	viper := viper.New()
	// viper.Set("src_url", "redis://:cvVAfYGseGhjZ8eh@192.168.3.37:6379/0")
	viper.Set("src_url", "redis://:123456@192.168.3.72:16379/14")
	viper.Set("backup_dir", "./MixedAI/docker_packages/opencv")
	viper.Set("update_force", false)
	viper.Set("includes", []string{})
	viper.Set("excludes", []string{".trash/*"})
	viper.Set("threads", 10)

	base := newBase(
		config.ModelConfig{
			DumpPath: "/data/backups",
		},
		// Creating a new base object.
		config.SubConfig{
			Type:  "juicefs",
			Name:  "juicefs-test",
			Viper: viper,
		},
	)

	db := &Juicefs{
		Base: base,
	}

	err := db.init()
	assert.NoError(t, err)
	// db.stage = "src"
	// script := db.build()
	// fmt.Print(script)
	fmt.Printf("perform error:%v", db.perform())
	// assert.Equal(t, script, "juicefs --host 1.2.3.4 --port 1234 -u user1 -ppass1 --ignore-table=my_db.aa --ignore-table=my_db.bb --a1 --a2 --a3 my_db foo bar --result-file=/data/backups/mysql/mysql1/my_db.sql")
}
