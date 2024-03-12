package status

import (
	"errors"
	"github.com/cossim/hipush/config"
	"github.com/cossim/hipush/pkg/store"
	"github.com/thoas/stats"
)

// Stats provide response time, status code count, etc.
var Stats *stats.Stats

// StatStorage implements the storage interface
var StatStorage *StateStorage

func InitAppStatus(cfg *config.Config) error {
	var s store.Store

	switch cfg.Storage.Type {
	case "memory":
		s = store.NewMemoryStore()
	case "file":
		s = store.NewFileStore(cfg.Storage.Path)
	default:
		//logx.LogError.Error("storage error: can't find storage driver")
		return errors.New("can't find storage driver")
	}

	StatStorage = NewStateStorage(s)
	return StatStorage.Init()
}

// App is status structure
type App struct {
	Version        string        `json:"version"`
	BusyWorkers    int           `json:"busy_workers"`
	SuccessTasks   int           `json:"success_tasks"`
	FailureTasks   int           `json:"failure_tasks"`
	SubmittedTasks int           `json:"submitted_tasks"`
	TotalCount     int64         `json:"total_count"`
	Ios            IosStatus     `json:"ios"`
	Android        AndroidStatus `json:"android"`
	Huawei         HuaweiStatus  `json:"huawei"`
}

// AndroidStatus is android structure
type AndroidStatus struct {
	PushSuccess int64 `json:"push_success"`
	PushError   int64 `json:"push_error"`
}

// IosStatus is iOS structure
type IosStatus struct {
	PushSuccess int64 `json:"push_success"`
	PushError   int64 `json:"push_error"`
}

// HuaweiStatus is huawei structure
type HuaweiStatus struct {
	PushSuccess int64 `json:"push_success"`
	PushError   int64 `json:"push_error"`
}

//// InitAppStatus for initialize app status
//func InitAppStatus(conf *config.Config) (*StateStorage, error) {
//	//logx.LogAccess.Info("Init App Status Engine as ", conf.Stat.Engine)
//
//	var store store2.Store
//	//nolint:goconst
//	switch "" {
//	case "memory":
//		store = store2.NewMemoryStore()
//	//case "redis":
//	//	store = redis.New(conf)
//	//case "boltdb":
//	//	store = boltdb.New(conf)
//	//case "buntdb":
//	//	store = buntdb.New(conf)
//	//case "leveldb":
//	//	store = leveldb.New(conf)
//	//case "badger":
//	//	store = badger.New(conf)
//	default:
//		log.Printf("storage error: can't find storage driver")
//		return nil, errors.New("can't find storage driver")
//	}
//
//	StatStorage := NewStateStorage(store)
//
//	if err := StatStorage.Init(); err != nil {
//		log.Printf("storage error: " + err.Error())
//		return nil, err
//	}
//
//	return StatStorage, nil
//}
