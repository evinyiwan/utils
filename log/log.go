package log

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"strings"
	"time"
	"xy-kf-gin/pkg/conf"
)

var (
	log *zap.Logger
)

func Logger() *zap.Logger {
	if log != nil {
		return log
	}
	return nil
}


func LoggerInit(config conf.InitConf, logLevel string) {

	env := config.Env()

	logFilePath := config.LogPath()

	logFileName := time.Now().Format("20060102") + ".log"

	fileName := path.Join(logFilePath, logFileName)

	//日志分割
	hook := &lumberjack.Logger{
		Filename:   fileName, // 日志文件路径，默认 os.TempDir()
		MaxSize:    100,      // 每个日志文件保存10M，默认 100M
		MaxBackups: 30,       // 保留30个备份，默认不限
		MaxAge:     7,        // 保留7天，默认不限
		Compress:   true,     // 是否压缩，默认不压缩
	}
	write := zapcore.AddSync(hook)

	// 设置日志级别
	var level zapcore.Level
	switch strings.ToLower(logLevel) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	default:
		level = zapcore.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     encodeTime,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	//开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	//开启文件及行号
	development := zap.Development()
	//设置初始化字段，添加一个服务器名称
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))

	var (
		encode zapcore.Encoder
	)
	// dev模式时错误信息打印到控制台
	if env == "dev" {
		encode = zapcore.NewConsoleEncoder(encoderConfig)
		write = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)) // 打印到控制台
	}else{
		encode = zapcore.NewJSONEncoder(encoderConfig)
	}

	//构造
	core := zapcore.NewCore(
		encode,
		write,	 // 打印到文件
		atomicLevel,
	)
	log = zap.New(core, caller, development)
	log.Info("Logger init success")

}

//格式化日志时间
func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder){
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
