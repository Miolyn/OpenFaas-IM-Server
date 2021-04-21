package logger

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	oplogging "github.com/op/go-logging"
)

var (
	Logger           *oplogging.Logger
	Module           = "Gimme"
	prefix           = "[Gimme]"
	level            = "DEBUG"
	logDir           = "log"
	defaultFormatter = `%{color:bold} %{time:2006/01/02 15:04:05} %{shortfile} ▶ [%{level:.6s}] %{message}%{color:reset}`
)

func init() {
	Logger = GetLoggerModule(Module)

}
func InitWithMain(name string) {
	GetLoggerModule(Module + "-" + name)
}

func registerStdout(out io.Writer, backends []oplogging.Backend, module string) []oplogging.Backend {

	level, err := oplogging.LogLevel(level)
	if err != nil {
		fmt.Println(err)
	}

	format := oplogging.MustStringFormatter(defaultFormatter)
	backend := oplogging.NewLogBackend(os.Stdout, prefix, 0)

	backendLeveled := oplogging.AddModuleLevel(oplogging.NewBackendFormatter(backend, format))
	backendLeveled.SetLevel(level, module)
	backends = append(backends, backendLeveled)
	return backends
}

func GetLoggerModule(module string) *oplogging.Logger {

	logger := oplogging.MustGetLogger(module)
	var backends []oplogging.Backend

	prefix = fmt.Sprintf("[%s]", module)
	backends = registerStdout(os.Stdout, backends, module)

	backends = append(backends, getFileBackend(module))

	logger.SetBackend(oplogging.SetBackend(backends...))
	return logger
}

func getFileBackend(module string) oplogging.LeveledBackend {
	//判断是否存在该文件夹
	if err := os.MkdirAll(logDir, 0777); err != nil {
		panic(err)
	}
	// 打开一个文件
	file, err := os.OpenFile(path.Join(logDir, module+"_info.log"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	//backend := l.getLogBackend(file, LogLevelMap[l.level])
	//level, err := oplogging.LogLevel(c.Stdout)
	level, err := oplogging.LogLevel(level)
	if err != nil {
		panic(err)
	}
	backend := getLogBackend(file, int(level))
	//logging.SetBackend(backend)
	return backend
}

func getLogBackend(out io.Writer, level int) oplogging.LeveledBackend {
	pattern := defaultFormatter
	pattern = strings.Replace(pattern, "%{color:bold}", "", -1)
	pattern = strings.Replace(pattern, "%{color:reset}", "", -1)
	//if !c.Logfile {
	//	//remove %{logfile} tag
	//	pattern = strings.Replace(pattern, "%{longfile}", "", -1)
	//}

	//pattern = strings.Replace(pattern, "%{longfile}", "", -1)

	backend := oplogging.NewLogBackend(out, prefix, 1)
	format := oplogging.MustStringFormatter(pattern)
	backendFormatter := oplogging.NewBackendFormatter(backend, format)
	backendLeveled := oplogging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(oplogging.Level(level), "")
	return backendLeveled
}
