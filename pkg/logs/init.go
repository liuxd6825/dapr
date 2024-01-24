package logs

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"time"
)

func NewWriter(saveFile string, saveDays int, rotationHour int) io.Writer {
	logFile := saveFile + ".%Y-%m-%d-%H.log"
	// 配置日志每隔 1 小时轮转一个新文件，保留最近 30 天的日志文件，多余的自动清理掉。
	writer, _ := rotatelogs.New(
		logFile,
		rotatelogs.WithLinkName(saveFile),
		rotatelogs.WithMaxAge(time.Duration(24*saveDays)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(rotationHour)*time.Hour),
	)
	return writer
}

func SetOutput() {

}
