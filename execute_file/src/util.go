package src

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"v2ray.os.executable.file/src/config"
)

func GetPath(path []string) string {
	// 获取当前可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("获取执行文件路径错误:  %s", err)
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(append([]string{exeDir}, path...)...)
}

func AbsPathByRelativePath(filePath string) string {
	relativePath := strings.Split(filePath, "/")
	return GetPath(relativePath)
}

func GetConfigByType[T any](filePath string, config *T) {

	// 读取 YAML 文件
	absPath := AbsPathByRelativePath(filePath)
	viper.SetConfigFile(absPath)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("fatal error config file: %v \n", err)
	}

	if err := viper.Unmarshal(config); err != nil {
		log.Fatalf("unable to decode into struct: %v", err)
	}
}

func GetJsonByType[T any](filePath string, config *T) {

	// 读取 YAML 文件
	absPath := AbsPathByRelativePath(filePath)

	data, err := os.ReadFile(absPath)
	if err != nil {
		log.Fatalf("unable to read config file: %v", err)
	}

	if err := json.Unmarshal(data, config); err != nil {
		log.Fatalf("fatal error config file: %v \n", err)
	}
}

func InterfaceTOJson(data any) []byte {

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(fmt.Errorf("error marshaling config to JSON: %v", err))
	}
	return jsonData
}

func WriteToFile(filePath string, data []byte, role fs.FileMode) {

	// 將 JSON 寫回到文件中
	absPath := AbsPathByRelativePath(filePath)
	if err := os.WriteFile(absPath, data, role); err != nil {
		panic(fmt.Errorf("error writing JSON to file: %v", err))
	}

}

func splitKeyAndIndexBystring(key string) (arrayKey string, arrayIndex int) {
	index := strings.Index(key, "[")
	if index == -1 {
		return key, 0
	}

	return key[:index], cast.ToInt(key[index+1 : len(key)-1])
}

// SetNestedField 設置嵌套 map 中的值
func SetNestedField(data map[string]any, path string, value interface{}) error {
	keys := strings.Split(path, ".")
	m := data
	for i, key := range keys {
		if i == len(keys)-1 {
			m[key] = value
			return nil
		}

		key, index := splitKeyAndIndexBystring(key)

		// 處理切片的情況
		switch m[key].(type) {
		case map[string]interface{}:
			m = m[key].(map[string]interface{})
		case []interface{}:

			mArray, ok := m[key].([]interface{})
			if !ok {
				return fmt.Errorf("key %s is not an array", key)
			}

			if index >= len(mArray) {
				return fmt.Errorf("array index out of bounds: %d", index)
			}

			m = mArray[index].(map[string]interface{})
		default:
			return fmt.Errorf("unexpected type for key %s", key)
		}
	}
	return nil
}

func GetConfigByInterface(data map[string]any, path string) any {
	keys := strings.Split(path, ".")
	m := data
	for i, key := range keys {
		if i == len(keys)-1 {
			return m[key] // 返回指標m
		}

		key, index := splitKeyAndIndexBystring(key)

		// 處理切片的情況
		switch m[key].(type) {
		case map[string]interface{}:
			m = m[key].(map[string]interface{})
		case []interface{}:

			mArray, ok := m[key].([]interface{})
			if !ok {
				return nil
			}
			m = mArray[index].(map[string]interface{})
		default:
			return nil
		}
	}
	return nil
}

func AnyToStructByMapstructure[T any](data interface{}, t *T) error {
	return mapstructure.Decode(data, t)

}

func IsProduction() bool {
	return !config.Debug
}

func getPrefix(ip string, prefix string) string {
	// 找到第一個 '.' 的位置
	index := strings.Index(ip, prefix)
	if index == -1 {
		// 如果找不到 '.'，返回空字串
		return ""
	}
	// 擷取從開始到第一個 '.' 的部分
	return ip[:index]

}
