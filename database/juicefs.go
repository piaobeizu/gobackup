package database

import (
	"fmt"
	"strings"

	"github.com/gobackup/gobackup/helper"
	"github.com/gobackup/gobackup/logger"
)

// MySQL database
//
// type: juicefs
// host: 127.0.0.1
// port: 3306
// socket:
// database:
// username: root
// password:
// args:
type Juicefs struct {
	Base
	srcUrl string
	dstUrl string
	// cacheSize   string
	// cacheDir    string
	backupDir   string
	forceUpdate bool
	includes    []string
	excludes    []string
	threads     int

	stage    string
	srcMount string
}

func (db *Juicefs) init() (err error) {
	viper := db.viper

	db.srcUrl = viper.GetString("src_url")
	db.backupDir = viper.GetString("backup_dir")
	db.forceUpdate = viper.GetBool("update_force")
	db.includes = viper.GetStringSlice("includes")
	db.excludes = viper.GetStringSlice("excludes")
	db.threads = viper.GetInt("threads")

	return nil
}

func (db *Juicefs) build() string {
	dumpArgs := []string{}
	// if db.stage == "dst" {
	// 		if !strings.HasPrefix(db.dstUrl, "file://") {
	// 			db.dstMount = fmt.Sprintf("/jfs/sync/dst-%s", db.name)
	// 			if len(db.srcUrl) > 0 {
	// 				dumpArgs = append(dumpArgs, db.dstUrl)
	// 			}
	// 			dumpArgs = append(dumpArgs, db.dstMount)
	// 			return "juicefs mount -d" + " " + strings.Join(dumpArgs, " ")
	// 		} else {
	// 			db.dstMount = db.dstUrl[7:]
	// 			return "echo 'you dont need to mount if used local path'"
	// 		}
	// 	}
	// if db.stage == "udst" {
	// 	if !strings.HasPrefix(db.dstUrl, "file://") {
	// 		db.dstMount = fmt.Sprintf("/jfs/sync/dst-%s", db.name)
	// 		dumpArgs = append(dumpArgs, db.dstMount)
	// 		return "juicefs umount " + strings.Join(dumpArgs, " ")
	// 	} else {
	// 		return "echo 'you dont need to umount if used local path'"
	// 	}
	// }
	if db.stage == "src" {
		db.srcMount = fmt.Sprintf("/jfs/sync/src-%s", db.name)
		if len(db.srcUrl) > 0 {
			dumpArgs = append(dumpArgs, db.srcUrl)
		}
		dumpArgs = append(dumpArgs, db.srcMount)
		return "juicefs mount -d" + " " + strings.Join(dumpArgs, " ")
	} else if db.stage == "usrc" {
		db.srcMount = fmt.Sprintf("/jfs/sync/src-%s", db.name)
		dumpArgs = append(dumpArgs, db.srcMount)
		return "juicefs umount " + strings.Join(dumpArgs, " ")
	} else {
		if db.forceUpdate {
			dumpArgs = append(dumpArgs, "--force-update")
		}
		if db.threads > 0 {
			dumpArgs = append(dumpArgs, "--threads", fmt.Sprintf("%d", db.threads))
		}
		if len(db.includes) > 0 {
			for _, include := range db.includes {
				dumpArgs = append(dumpArgs, "--include", include)
			}
		}
		if len(db.excludes) > 0 {
			for _, exclude := range db.excludes {
				dumpArgs = append(dumpArgs, "--exclude", exclude)
			}
		}
		dumpArgs = append(dumpArgs, fmt.Sprintf("%s/%s", db.srcMount, db.backupDir), fmt.Sprintf("%s/%s", db.dumpPath, db.backupDir))
	}

	return "juicefs sync " + " " + strings.Join(dumpArgs, " ")
}

func (db *Juicefs) perform() error {
	logger := logger.Tag("Juicefs")

	logger.Info("-> juicefs mount source data...")
	db.stage = "src"
	_, err := helper.Exec(db.build())
	if err != nil {
		return fmt.Errorf("-> Mount %s error: %s", db.srcMount, err)
	}

	// logger.Info("-> juicefs mount destinate data...")
	// db.stage = "dst"
	// _, err = helper.Exec(db.build())
	// if err != nil {
	// 	return fmt.Errorf("-> Mount %s error: %s", db.dstMount, err)
	// }

	logger.Info("-> juicefs sync data...")
	db.stage = "sync"
	_, err = helper.Exec(db.build())
	if err != nil {
		return fmt.Errorf("-> sync error: %s", err)
	}

	logger.Info("-> juicefs umount source data...")
	db.stage = "usrc"
	_, err = helper.Exec(db.build())
	if err != nil {
		return fmt.Errorf("-> Umount %s error: %s", db.srcMount, err)
	}

	// logger.Info("-> juicefs umount destinate data...")
	// db.stage = "udst"
	// _, err = helper.Exec(db.build())
	// if err != nil {
	// 	return fmt.Errorf("-> Umount %s error: %s", db.dstMount, err)
	// }
	logger.Info("dump path:", db.dumpPath)
	return nil
}
