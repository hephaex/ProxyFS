package main

import (
	"fmt"
	"os"

	pfsagentConfig "github.com/swiftstack/ProxyFS/pfsagentConfig"
)

// var (
// MyConfig a var to hold config params
// MyConfig *Config
// configPath = defaultConfigPath
// )

const (
	defaultConfigPath string = "/etc/pfsagent/"
)

//
// // Config the pfsagentd config
// type Config struct {
// 	FUSEVolumeName               string  `ini:"FUSEVolumeName"`
// 	FUSEMountPointPath           string  `ini:"FUSEMountPointPath"`
// 	FUSEUnMountRetryDelay        string  `ini:"FUSEUnMountRetryDelay"`
// 	FUSEUnMountRetryCap          int     `ini:"FUSEUnMountRetryCap"`
// 	SwiftAuthURL                 string  `ini:"SwiftAuthURL"`
// 	SwiftAuthUser                string  `ini:"SwiftAuthUser"`
// 	SwiftAuthKey                 string  `ini:"SwiftAuthKey"`
// 	SwiftAccountName             string  `ini:"SwiftAccountName"`
// 	SwiftTimeout                 string  `ini:"SwiftTimeout"`
// 	SwiftRetryLimit              int     `ini:"SwiftRetryLimit"`
// 	SwiftRetryDelay              string  `ini:"SwiftRetryDelay"`
// 	SwiftRetryExpBackoff         float64 `ini:"SwiftRetryExpBackoff"`
// 	SwiftConnectionPoolSize      int     `ini:"SwiftConnectionPoolSize"`
// 	FetchExtentsFromFileOffset   int     `ini:"FetchExtentsFromFileOffset"`
// 	FetchExtentsBeforeFileOffset int     `ini:"FetchExtentsBeforeFileOffset"`
// 	ReadCacheLineSize            int     `ini:"ReadCacheLineSize"`
// 	ReadCacheLineCount           int     `ini:"ReadCacheLineCount"`
// 	SharedFileLimit              int     `ini:"SharedFileLimit"`
// 	ExclusiveFileLimit           int     `ini:"ExclusiveFileLimit"`
// 	DirtyFileLimit               int     `ini:"DirtyFileLimit"`
// 	MaxFlushSize                 int     `ini:"MaxFlushSize"`
// 	MaxFlushTime                 string  `ini:"MaxFlushTime"`
// 	LogFilePath                  string  `ini:"LogFilePath"`
// 	LogToConsole                 string  `ini:"LogToConsole"`
// 	TraceEnabled                 string  `ini:"TraceEnabled"`
// 	HTTPServerIPAddr             string  `ini:"HTTPServerIPAddr"`
// 	HTTPServerTCPPort            int     `ini:"HTTPServerTCPPort"`
// 	ReadDirPlusEnabled           bool    `ini:"ReadDirPlusEnabled"`
// 	XAttrEnabled                 bool    `ini:"XAttrEnabled"`
// 	EntryDuration                string  `ini:"EntryDuration"`
// 	AttrDuration                 string  `ini:"AttrDuration"`
// 	AttrBlockSize                int     `ini:"AttrBlockSize"`
// 	ReaddirMaxEntries            int     `ini:"ReaddirMaxEntries"`
// 	FUSEMaxBackground            int     `ini:"FUSEMaxBackground"`
// 	FUSECongestionThreshhold     int     `ini:"FUSECongestionThreshhold"`
// 	FUSEMaxWrite                 int     `ini:"FUSEMaxWrite"`
// }

func main() {
	if len(os.Args) == 1 {
		loadError := pfsagentConfig.LoadConfig("")
		if nil != loadError {
			fmt.Println("Failed loading config. Error:", loadError)
			os.Exit(1)
		}
		pfsagentConfig.RunStateMachine()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "firstRun":
		err := firstTimeRun()
		if nil != err {
			fmt.Println(err)
			os.Exit(2)
		}
		os.Exit(0)
	case "configPath":
		if len(os.Args) > 1 {
			pfsagentConfig.ConfigPath = os.Args[2]
			loadError := pfsagentConfig.LoadConfig("")
			if nil != loadError {
				fmt.Println("Failed loading config. Error:", loadError)
				os.Exit(1)
			}
			pfsagentConfig.RunStateMachine()
			os.Exit(0)
		}
		fmt.Print(usageText)
		os.Exit(2)
	case "version":
		fmt.Print(version)
		os.Exit(0)
	case "usage":
		fmt.Print(usageText)
		os.Exit(0)
	default:
		fmt.Print(usageText)
		os.Exit(2)
	}
}

func firstTimeRun() error {
	return pfsagentConfig.FirstTimeRun()
}
