package log_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/dmzlingyin/utils/log"
)

func TestLog(t *testing.T) {
	printTestData()
	log.SetLevel(log.LevelError)
	fmt.Println("-------------------------")
	printTestData()
}

func TestNewLogger(t *testing.T) {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	l := log.New(log.LevelDebug, file)
	log.SetLogger(l)
	printTestData()
}

func BenchmarkLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		printTestData()
	}
}

func printTestData() {
	log.Debug("Hello World")
	log.Debugf("This is a test %s", "debug")
	log.Info("Hello World")
	log.Infof("This is a test %s", "info")
	log.Warn("Hello World")
	log.Warnf("This is a test %s", "warn")
	log.Error("Hello World")
	log.Errorf("This is a test %s", "error")
}
