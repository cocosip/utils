package buildinfo

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
)

// BuildInfo holds all relevant application build and version information.
type BuildInfo struct {
	Version     string `json:"version"`     // 应用程序版本，遵循语义化版本规范 (e.g., "v1.0.0")
	BuildCommit string `json:"buildCommit"` // Git 提交哈希 (e.g., "abcdef123456")
	BuildTime   string `json:"buildTime"`   // 构建时间 (UTC, e.g., "2023-10-27T10:00:00Z")
	DevMode     bool   `json:"devMode"`     // 是否处于开发模式
}

var (
	// 这些变量将在编译时通过 -ldflags 注入。
	// 例如: go build -ldflags "-X github.com/cocosip/utils/buildinfo.ldVersion=v1.0.0"
	ldVersion     = "DEVELOPMENT Version"
	ldBuildCommit = "DEVELOPMENT Build"
	ldBuildTime   = "UNKNOWN"
	ldDevMode     = "true"
)

var info BuildInfo // 存储解析后的构建信息

func init() {
	info = BuildInfo{
		Version:     ldVersion,
		BuildCommit: ldBuildCommit,
		BuildTime:   ldBuildTime,
		DevMode:     parseBool(ldDevMode),
	}
}

// parseBool 是一个辅助函数，用于安全地解析布尔字符串。
// 如果解析失败，它将记录错误并返回 false。
func parseBool(s string) bool {
	b, err := strconv.ParseBool(s)
	if err != nil {
		slog.Error("Failed to parse boolean string", "string", s, "error", err)
		return false // 默认返回 false
	}
	return b
}

// Get 返回填充了构建信息的 BuildInfo 结构体实例。
func Get() BuildInfo {
	return info
}

// String 返回 BuildInfo 的人类可读字符串表示形式。
func (bi BuildInfo) String() string {
	return fmt.Sprintf("Version: %s\nBuild: %s\nBuilt At: %s\nDev Mode: %t",
		bi.Version, bi.BuildCommit, bi.BuildTime, bi.DevMode)
}

// JSON 返回 BuildInfo 的 JSON 字符串表示形式。
// 如果 Marshal 失败，它将返回一个空字符串和错误。
func (bi BuildInfo) JSON() (string, error) {
	data, err := json.MarshalIndent(bi, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal BuildInfo to JSON: %w", err)
	}
	return string(data), nil
}
