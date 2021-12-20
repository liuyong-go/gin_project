package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/liuyong-go/gin_project/libs/ydefer"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerStruct struct {
	Development  bool         //开发，线上环境
	InfoPath     string       `yaml:"infoPath"`
	ErrorPath    string       `yaml:"errorPath"`
	EncodeConfig EncodeConfig `yaml:"encodeConfig"`
	RotationLogs RotateLogs   `yaml:"rotationLogs"`
	Async        bool         `yaml:"async"`
	Buffer       LogBuffer    `yaml:"buffer"`
}
type RotateLogs struct {
	MaxAge       time.Duration `yaml:"maxAge"`       // 保存小时数
	RotationTime time.Duration `yaml:"rotationTime"` //切割频率 小时记录
}
type EncodeConfig struct {
	MessageKey string `yaml:"messageKey"`
	LevelKey   string `yaml:"levelKey"`
	TimeKey    string `yaml:"timeKey"`
	CallerKey  string `yaml:"callerKey"`
}
type LogBuffer struct {
	BufferSize    int           `yaml:"bufferSize"`
	FlushInterval time.Duration `yaml:"flushInterval"`
}

var logger *zap.SugaredLogger

func InitLogger(conf LoggerStruct) {
	// config := zapcore.EncoderConfig{}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:   conf.EncodeConfig.MessageKey,
		LevelKey:     conf.EncodeConfig.LevelKey,  //结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:      conf.EncodeConfig.TimeKey,   //结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		CallerKey:    conf.EncodeConfig.CallerKey, //结构化（json）输出：打印日志的文件对应的Key
		EncodeLevel:  zapcore.CapitalLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,  //采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}, //输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		}, //
	}
	//自定义日志级别：自定义Info级别
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel && lvl >= zapcore.InfoLevel
	})
	//自定义日志级别：自定义Warn级别
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	var wsInfo zapcore.WriteSyncer
	var wsWarn zapcore.WriteSyncer
	var core zapcore.Core
	if conf.Development == true {
		wsInfo = os.Stdout
	} else {
		wsInfo = zapcore.AddSync(getWriter(conf.RotationLogs, conf.InfoPath))
		wsWarn = zapcore.AddSync(getWriter(conf.RotationLogs, conf.ErrorPath))
	}
	if conf.Async == true {
		var closeInfo, closeWarn CloseFunc
		wsInfo, closeInfo = Buffer(wsInfo, conf.Buffer.BufferSize, conf.Buffer.FlushInterval*time.Second)
		wsWarn, closeWarn = Buffer(wsWarn, conf.Buffer.BufferSize, conf.Buffer.FlushInterval*time.Second)
		ydefer.Register(closeInfo)
		ydefer.Register(closeWarn)
	}
	if conf.Development == true {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), wsInfo, infoLevel), //同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), wsInfo, infoLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), wsWarn, warnLevel),
		)
	}
	//实现多个输出

	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.WarnLevel)).Sugar()
}
func getWriter(rl RotateLogs, filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+".%Y%m%d%H",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*rl.MaxAge),             // 保存30天
		rotatelogs.WithRotationTime(time.Hour*rl.RotationTime), //切割频率 24小时
	)
	if err != nil {
		panic(err)
	}
	return hook
}
func Panic(ctx context.Context, data ...interface{}) {
	data = withTrace(ctx, data...)
	logger.Panic(data...)
}

func Info(ctx context.Context, data ...interface{}) {
	data = withTrace(ctx, data...)
	logger.Info(data...)
}
func withTrace(ctx context.Context, data ...interface{}) (returnData []interface{}) {
	_, file, line, _ := runtime.Caller(2)
	returnData = append(data, fmt.Sprintf(" ,path: file:%s,line :%d", file, line))
	ltrace := ctx.Value("trace")
	if ltrace != nil {
		traceStr := fmt.Sprintf(", traceID:%s,spanID:%s", ltrace.(*Trace).TraceID, ltrace.(*Trace).SpanID)
		returnData = append(returnData, traceStr)
	}

	return
}

func Warn(ctx context.Context, data ...interface{}) {
	data = withTrace(ctx, data...)
	logger.Warn(data...)
}
