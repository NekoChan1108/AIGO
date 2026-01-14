package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	// ANSI 颜色代码
	reset  = "\033[0m"  // 重置
	red    = "\033[31m" // 红色
	green  = "\033[32m" // 绿色
	yellow = "\033[33m" // 黄色
	blue   = "\033[34m" // 蓝色
	purple = "\033[35m" // 紫色
	cyan   = "\033[36m" // 青色
	gray   = "\033[37m" // 灰色
	white  = "\033[97m" // 白色
)

// 图标常量
const (
	successIcon = "✅"
	errorIcon   = "❌"
	warnIcon    = "⚠️"
	infoIcon    = "ℹ️"
	debugIcon   = "🔍"
	fatalIcon   = "💥"
	panicIcon   = "🚨"
)

// colorLogger 结构体定义
type colorLogger struct {
	logger *log.Logger
}

// 创建一个新的彩色日志实例
func newColorLogger(out io.Writer) *colorLogger {
	return &colorLogger{
		// |log.Lshortfile 防止出现文件名行数等信息 只带上日期和时间
		logger: log.New(out, "", log.LstdFlags),
	}
}

// 全局日志实例
var defaultLogger = newColorLogger(os.Stdout)

// colorize 返回带颜色的字符串
/*
@description: 带颜色的字符串
@param color 颜色代码
@param text 要格式化的文本
@return string
*/
func colorize(color, text string) string {
	return color + text + reset
}

// getColoredPrefix 返回带颜色和图标的前缀
func getColoredPrefix(color, icon, level string) string {
	return colorize(color, fmt.Sprintf("[%s%s] ", icon, level))
}

// Error 红色错误日志
func (l *colorLogger) Error(v ...any) {
	l.logger.SetPrefix(getColoredPrefix(red, errorIcon, "ERROR"))
	l.logger.Output(1, colorize(red, fmt.Sprint(v...)))
}

// Errorf 带格式化的红色错误日志
func (l *colorLogger) Errorf(format string, v ...any) {
	l.logger.SetPrefix(getColoredPrefix(red, errorIcon, "ERROR"))
	l.logger.Output(1, colorize(red, fmt.Sprintf(format, v...)))
}

// Warn 黄色警告日志
func (l *colorLogger) Warn(v ...any) {
	l.logger.SetPrefix(getColoredPrefix(yellow, warnIcon, "WARN"))
	l.logger.Output(1, colorize(yellow, fmt.Sprint(v...)))
}

// Warnf 带格式化的黄色警告日志
func (l *colorLogger) Warnf(format string, v ...any) {
	l.logger.SetPrefix(getColoredPrefix(yellow, warnIcon, "WARN"))
	l.logger.Output(1, colorize(yellow, fmt.Sprintf(format, v...)))
}

// Info 蓝色信息日志
func (l *colorLogger) Info(v ...any) {
	l.logger.SetPrefix(getColoredPrefix(blue, infoIcon, "INFO"))
	l.logger.Output(1, colorize(blue, fmt.Sprint(v...)))
}

// Infof 带格式化的蓝色信息日志
func (l *colorLogger) Infof(format string, v ...any) {
	l.logger.SetPrefix(getColoredPrefix(blue, infoIcon, "INFO"))
	l.logger.Output(1, colorize(blue, fmt.Sprintf(format, v...)))
}

// Debug 青色调试日志
func (l *colorLogger) Debug(v ...any) {
	l.logger.SetPrefix(getColoredPrefix(cyan, debugIcon, "DEBUG"))
	l.logger.Output(1, colorize(cyan, fmt.Sprint(v...)))
}

// Debugf 带格式化的青色调试日志
func (l *colorLogger) Debugf(format string, v ...any) {
	l.logger.SetPrefix(getColoredPrefix(cyan, debugIcon, "DEBUG"))
	l.logger.Output(1, colorize(cyan, fmt.Sprintf(format, v...)))
}

// Success 绿色成功日志
func (l *colorLogger) Success(v ...any) {
	l.logger.SetPrefix(getColoredPrefix(green, successIcon, "SUCCESS"))
	l.logger.Output(1, colorize(green, fmt.Sprint(v...)))
}

// Successf 带格式化的绿色成功日志
func (l *colorLogger) Successf(format string, v ...any) {
	l.logger.SetPrefix(getColoredPrefix(green, successIcon, "SUCCESS"))
	l.logger.Output(1, colorize(green, fmt.Sprintf(format, v...)))
}

// Fatal 紫色致命错误日志
func (l *colorLogger) Fatal(v ...any) {
	l.logger.SetPrefix(getColoredPrefix(purple, fatalIcon, "FATAL"))
	l.logger.Output(1, colorize(purple, fmt.Sprint(v...)))
	os.Exit(1)
}

// Fatalf 带格式化的紫色致命错误日志
func (l *colorLogger) Fatalf(format string, v ...any) {
	l.logger.SetPrefix(getColoredPrefix(purple, fatalIcon, "FATAL"))
	l.logger.Output(1, colorize(purple, fmt.Sprintf(format, v...)))
	os.Exit(1)
}

// Panic 红色panic日志
func (l *colorLogger) Panic(v ...any) {
	l.logger.SetPrefix(getColoredPrefix(red, panicIcon, "PANIC"))
	l.logger.Output(1, colorize(red, fmt.Sprint(v...)))
	panic(colorize(red, fmt.Sprint(v...)))
}

// Panicf 带格式化的红色panic日志
func (l *colorLogger) Panicf(format string, v ...any) {
	l.logger.SetPrefix(getColoredPrefix(red, panicIcon, "PANIC"))
	l.logger.Output(1, colorize(red, fmt.Sprintf(format, v...)))
	panic(colorize(red, fmt.Sprintf(format, v...)))
}

// 暴露默认实例的方法 - 全局函数

// Error 错误日志
func Error(v ...any) {
	defaultLogger.Error(v...)
}

// Errorf 带格式化的错误日志
func Errorf(format string, v ...any) {
	defaultLogger.Errorf(format, v...)
}

// Warn 警告日志
func Warn(v ...any) {
	defaultLogger.Warn(v...)
}

// Warnf 带格式化的警告日志
func Warnf(format string, v ...any) {
	defaultLogger.Warnf(format, v...)
}

// Info 信息日志
func Info(v ...any) {
	defaultLogger.Info(v...)
}

// Infof 带格式化的信息日志
func Infof(format string, v ...any) {
	defaultLogger.Infof(format, v...)
}

// Debug 调试日志
func Debug(v ...any) {
	defaultLogger.Debug(v...)
}

// Debugf 带格式化的调试日志
func Debugf(format string, v ...any) {
	defaultLogger.Debugf(format, v...)
}

// Success 成功日志
func Success(v ...any) {
	defaultLogger.Success(v...)
}

// Successf 带格式化的成功日志
func Successf(format string, v ...any) {
	defaultLogger.Successf(format, v...)
}

// Fatal 致命错误日志
func Fatal(v ...any) {
	defaultLogger.Fatal(v...)
}

// Fatalf 带格式化的致命错误日志
func Fatalf(format string, v ...any) {
	defaultLogger.Fatalf(format, v...)
}

// Panic panic日志
func Panic(v ...any) {
	defaultLogger.Panic(v...)
}

// Panicf 带格式化的panic日志
func Panicf(format string, v ...any) {
	defaultLogger.Panicf(format, v...)
}
