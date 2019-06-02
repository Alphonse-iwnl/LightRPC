package utils

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func initLogger(){
	InitNewLogger("/Users/evan/golang/src/logs", "Debug", "unit_test")
}

func TestStat(t *testing.T) {
	initLogger()
	testState := NewStateEntry("statelog test")
	EndStateEntry(testState, 200)
	time.Sleep(time.Second * 30)
}

func TestLog(t *testing.T) {
	initLogger()
	LOG.Debug("test debug log")
	LOG.Debugf("test debug log, param:%s", "test param")
	LOG.Info("test info log")
	LOG.Infof("test info log, param:%s", "test param")
	LOG.Error("test error log")
	LOG.Errorf("test error log, param:%s", "test param")
	LOG.Warn("test warn log")
	LOG.Warnf("test warn log, param:%s", "test param")
	//LOG.Fatalf("test fatal log, param:%s","test param")
}

func TestBenchmarkStatLog(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < 100000; i++ {
				testState := NewStateEntry("statelog test")
				LOG.Info("log/state lib benchmark test.")
				EndStateEntry(testState, 200)
			}
		}()
	}
	wg.Wait()
	fmt.Printf("stat benchmark done.")
}

func TestTomlConfig(t *testing.T) {
	// filePath := ""
	// 注意路径
	initLogger()
	LoadTomlConfig(nil, "/Users/evan/golang/src/rpc/LightBackend/utils/unit_test_toml.toml")
	fmt.Println("Log info:",DefaultToml.Log)
	fmt.Println("Server info:",DefaultToml.Server)
	fmt.Println("Remote info",DefaultToml.ServerClient)
}

func TestProjectPath(t *testing.T) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	fmt.Println(path[:index])
}
