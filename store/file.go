package store

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultPath  = "/etc/hipush/data.json"
	saveInterval = 5 * time.Second // 延迟写入的时间间隔
	bufferSize   = 100             // 缓冲区大小
)

type FileStore struct {
	mutex      sync.Mutex
	path       string
	data       map[string]int64
	buffer     map[string]int64
	bufferSize int
	saveTicker *time.Ticker
}

func NewFileStore(path string) *FileStore {
	if path == "" {
		path = defaultPath
	}
	return &FileStore{
		path:       path,
		data:       make(map[string]int64),
		buffer:     make(map[string]int64),
		bufferSize: bufferSize,
	}
}

func (fs *FileStore) Init() error {
	dir := filepath.Dir(fs.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 检查文件是否存在，如果存在则加载数据到内存中
	if _, err := os.Stat(fs.path); !os.IsNotExist(err) {
		if err := fs.loadFromFile(); err != nil {
			return err
		}
	} else {
		if err := fs.saveToFile(); err != nil {
			return err
		}
	}
	// 启动定时保存任务
	fs.saveTicker = time.NewTicker(saveInterval)
	go fs.periodicSave()
	return nil
}

func (fs *FileStore) Get(key string) int64 {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	return fs.data[key]
}

func (fs *FileStore) Set(key string, value int64) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fs.buffer[key] = value
}

func (fs *FileStore) Add(key string, value int64) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	fs.buffer[key] += value
}

func (fs *FileStore) Del(key string) {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	delete(fs.buffer, key)
}

func (fs *FileStore) Close() error {
	// 停止定时保存任务
	fs.saveTicker.Stop()
	// 执行最后一次保存
	return fs.saveToFile()
}

func (fs *FileStore) loadFromFile() error {
	data, err := ioutil.ReadFile(fs.path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &fs.data)
}

func (fs *FileStore) saveToFile() error {
	fs.mutex.Lock()
	defer fs.mutex.Unlock()

	if len(fs.buffer) == 0 {
		return nil
	}

	// 从文件中读取已有的数据
	existingData := make(map[string]int64)
	if _, err := os.Stat(fs.path); !os.IsNotExist(err) {
		data, err := ioutil.ReadFile(fs.path)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(data, &existingData); err != nil {
			return err
		}
	}

	// 将缓冲区的数据合并到内存中
	for key, value := range fs.buffer {
		if oldValue, ok := existingData[key]; ok {
			// 如果文件中已有相应的键值对，则对其进行加1操作
			existingData[key] = oldValue + value
		} else {
			// 如果文件中没有相应的键值对，则直接写入缓冲区的值
			existingData[key] = value
		}
	}
	fs.buffer = make(map[string]int64)

	// 将合并后的数据写入文件
	data, err := json.Marshal(existingData)
	if err != nil {
		return err
	}
	fmt.Println("periodic save data to file data => ", string(data))

	dir := filepath.Dir(fs.path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(fs.path, data, 0644); err != nil {
		return err
	}

	// 重新加载文件数据到内存中
	return fs.loadFromFile()
}

func (fs *FileStore) periodicSave() {
	for range fs.saveTicker.C {
		if err := fs.saveToFile(); err != nil {
			log.Printf("save data to file error: %v", err)
		}
	}
}
