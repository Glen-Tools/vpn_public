package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"v2ray.os.executable.file/src"
	"v2ray.os.executable.file/src/config"
	"v2ray.os.executable.file/src/model"
)

func main() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("捕獲panic:", r)
		}
	}()

	// 讀取配置文件
	outBound := new(model.Outbound)
	src.GetConfigByType[model.Outbound](config.InboundYmlPath, outBound)
	count := len(outBound.User)
	if count == 0 {
		log.Fatalf("config user.yml is empty")
	}

	// 使用者切換，讀取輸入
	userNumber := 0
	if src.IsProduction() {
		isExit := false
		for !isExit {
			userNumber, isExit = getUserNumber(count)
		}

		if isExit && userNumber == 0 {
			log.Fatal("yml 的 user 為空，請確認填寫設定檔後再重新執行")
		}
	}

	// fmt.Println(userNumber)

	fmt.Println("載入設定檔中....")

	// 寫回json 複製配置
	var v2rayConfig map[string]any
	// var v2rayConfigCopy map[string]any
	src.GetJsonByType[map[string]any](config.ConfigJsonPath, &v2rayConfig)
	// copier.Copy(&v2rayConfigCopy, v2rayConfig)

	user := lo.If(src.IsProduction(), func() int {
		return userNumber - 1
	}).Else(func() int {
		return 2 //user id = 2
	})()

	// 設定 user id , log path 設定
	routing := new(model.Routing)
	src.GetConfigByType[model.Routing](config.RoutingYmlPath, routing)
	host := lo.Map(routing.Host, func(h model.Host, _ int) string {
		return fmt.Sprintf("%s%s", config.RoutingDomainPrefix, h.Url)
	})
	setJsonByKey(v2rayConfig, config.RoutingHostPath, host)
	setJsonByKey(v2rayConfig, config.OutboundsUserIdPath, outBound.User[user].Id)

	accessPath := lo.If(!src.IsProduction(), src.AbsPathByRelativePath(config.LogAccessName)).Else("")

	errorPath := src.AbsPathByRelativePath(config.LogErrorName)
	setJsonByKey(v2rayConfig, config.LogAccessPath, accessPath)
	setJsonByKey(v2rayConfig, config.LogErrorPath, errorPath)

	data := src.InterfaceTOJson(v2rayConfig)
	src.WriteToFile(config.ConfigJsonPath, data, 0644)

	// 讀取設定檔
	// inboundJson := src.GetConfigByInterface(v2rayConfigCopy, config.InboundPath)
	inboundJson := src.GetConfigByInterface(v2rayConfig, config.InboundPath)

	var proxy []*model.Proxy
	err := src.AnyToStructByMapstructure[[]*model.Proxy](inboundJson, &proxy)
	if err != nil {
		log.Fatalf("unable to set proxy into struct: %v", err)
	}

	fmt.Println("開啟proxy....")

	// 回傳需執行的 interface
	localOs := src.GetOS()
	localOs.RunV2ray()
	localOs.OpenProxy(proxy)

	time.Sleep(time.Second * 1)
	fmt.Println("v2ray 與 proxy 正在執行.... ， 想終止或離開請按 Ctrl+C 或直接關閉此窗口.")

	select {}
}

func getUserNumber(count int) (userNumber int, isExit bool) {

	fmt.Printf("切換使用者，請輸入 1 ~ %d 之間的數字，退出請輸入exit \n", count)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if err != nil {
		fmt.Println("讀取輸入時發生錯誤:", err)
		return 0, false
	}

	if input == "exit" {
		return 0, true
	}

	num := cast.ToInt(input)
	if num >= 1 && num <= count {
		return num, true
	}

	fmt.Println("輸入錯誤。")
	return 0, false
}

func setJsonByKey(data map[string]any, path string, value any) {
	err := src.SetNestedField(data, path, value)
	if err != nil {
		log.Fatalf("unable to set user id into struct: %v", err)
	}
}
