package src

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"

	"github.com/samber/lo"
	"golang.org/x/sys/windows"
	"v2ray.os.executable.file/src/config"
	"v2ray.os.executable.file/src/model"
)

type OperatingSystem interface {
	RunV2ray()
	OpenProxy(proxy []*model.Proxy)
	CloseProxy() error
	Exit()
	AddCommand(cmd *exec.Cmd)
}

// OS 结构体
type Windows struct {
	OsType      string
	CommandList []*exec.Cmd
	Proxy       []*model.Proxy
}

type Mac struct {
	OsType string
}

func (w *Windows) RunV2ray() {

	v2rayPath := AbsPathByRelativePath(config.WindowsV2rayPath)
	configPath := AbsPathByRelativePath(config.ConfigJsonPath)

	v2rayCommand := fmt.Sprintf("%s run -config=%s ", v2rayPath, configPath)
	commands := [][]string{
		{"-Command", fmt.Sprintf("chcp 65001 ; %s", v2rayCommand)},
	}

	command(w.OsType, commands, w)
}

func (w *Windows) OpenProxy(proxy []*model.Proxy) {

	defer func() {
		setupSignalHandler(w)

		if !IsProduction() {
			// time.Sleep(5 * time.Second)
			pid := uint32(os.Getpid())

			// Attach to the current console
			err := windows.GenerateConsoleCtrlEvent(windows.CTRL_BREAK_EVENT, pid)
			if err != nil {
				fmt.Printf("Failed to send signal: %v\n", err)
			}
		}
	}()

	// set proxy
	w.Proxy = proxy

	// 組合 ip:port
	ipWithPort := lo.Map(proxy, func(p *model.Proxy, _ int) string {
		return fmt.Sprintf("%s=%s:%s", p.Protocol, p.Listen, p.Port)
	})
	ip := lo.Map(proxy, func(p *model.Proxy, _ int) string {
		return fmt.Sprintf("%s.*", getPrefix(p.Listen, "."))
	})

	// 組合 ip , port
	listen := strings.Join(ipWithPort, ";")
	listenWildcard := strings.Join(lo.Uniq(ip), ";")
	// 命令
	commands := [][]string{
		{"-Command", fmt.Sprintf(`chcp 65001 ; Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyServer -Value "%s"`, listen)},
		{"-Command", `chcp 65001 ; Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyEnable -Value 1`},
		{"-Command", fmt.Sprintf(`chcp 65001 ; Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyOverride -Value "%s;"`, listenWildcard)},
		// {"-Command", `Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyServer -Value "127.0.0.1:1087"`},
		// {"-Command", `Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyEnable -Value 1"`},
		// {"-Command", `SetItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyOverride -Value "localhost;127.*;10.*;"`},
	}
	command(w.OsType, commands, w)

}

func (w *Windows) CloseProxy() error {
	commands := [][]string{{"-Command", `chcp 65001 ; Set-ItemProperty -Path "HKCU:\Software\Microsoft\Windows\CurrentVersion\Internet Settings" -Name ProxyEnable -Value 0`}}

	// 逐条执行
	err := command(w.OsType, commands, w)
	if err != nil {
		fmt.Println("please close your proxy by yourself")
	}
	return nil
}
func (w *Windows) Exit() {
	// 關閉代理
	w.CloseProxy()

	// 杀死子进程
	for _, cmd := range w.CommandList {
		if cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}

	// 杀死v2ray 进程
	var commands [][]string
	for _, value := range w.Proxy {
		data := []string{"-Command", fmt.Sprintf(`chcp 65001 ; $connection = Get-NetTCPConnection -LocalPort %s -LocalAddress %s | Where-Object { $_.State -ne 'TimeWait' }; if ($connection -and $connection.State -ne 'Idle') { Stop-Process -Id $connection.OwningProcess }`, value.Port, value.Listen)}
		commands = append(commands, data)
	}

	command(w.OsType, commands, w)

}
func (w *Windows) AddCommand(cmd *exec.Cmd) {
	w.CommandList = append(w.CommandList, cmd)
}

func (m *Mac) RunV2ray() {

}

func (m *Mac) OpenProxy(inbound []*model.Proxy) {
}
func (m *Mac) CloseProxy() error {
	return nil
}
func (m *Mac) Exit() {

}

func (m *Mac) AddCommand(cmd *exec.Cmd) {

}

func command(osType string, commands [][]string, system OperatingSystem) error {
	// 要執行的指令
	for _, commandLine := range commands {
		cmd := exec.Command(osType, commandLine...)

		err := cmd.Start()

		if err != nil {
			return err
		}
		system.AddCommand(cmd)
	}

	return nil
}

func GetOS() OperatingSystem {
	// 获取操作系统名称
	os := runtime.GOOS

	// 判斷操作系统
	switch os {
	case "windows":
		return &Windows{
			OsType: "powershell",
		}
	case "darwin":
		return &Mac{
			OsType: "sh",
		}
	default:
		panic(fmt.Sprintf("当前操作系统未知: %s", os))
	}
}

// setupSignalHandler 设置信号处理函数
func setupSignalHandler(system OperatingSystem) {
	// 创建一个通道来接收操作系统信号
	signalChan := make(chan os.Signal, 1)
	// 捕捉特定的信号，这里主要是终止信号
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	// 使用 Goroutine 来监听信号
	go func() {
		sig := <-signalChan
		fmt.Println("接收終止信號:", sig)
		// 调用你要在程序关闭时执行的函数
		system.Exit()
		os.Exit(0)
	}()
}
