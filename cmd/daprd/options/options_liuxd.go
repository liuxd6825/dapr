package options

import (
	"flag"
	"path/filepath"
	"strings"
)

// Flag type. Allows passing a flag multiple times to get a slice of strings.
// It implements the flag.Value interface.
type stringSliceFlag []string

// String formats the flag value.
func (f stringSliceFlag) String() string {
	return strings.Join(f, ",")
}

// Set the flag value.
func (f *stringSliceFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	*f = append(*f, value)
	return nil
}

// initOptions 1:解决路径参数取不到的问题
func initOptions(opts *Options, args []string) {
	// 1:解决路径参数取不到的问题 flag包的bug
	opts.LogFile = getArgValue(opts.LogFile, args, "-log-file")
	opts.ComponentsPath = getArgValue(opts.ComponentsPath, args, "-components-path")
	opts.ResourcesPath = getArgSlice(opts.ResourcesPath, args, "-resources-path")
	opts.Config = getArgSlice(opts.Config, args, "-config")

	// 2:将路径转为绝对路径
	opts.LogFile = getAbsPath(opts.LogFile)
	opts.ComponentsPath = getAbsPath(opts.ComponentsPath)
	opts.Config = sliceAbsPath(opts.Config)
	opts.ResourcesPath = sliceAbsPath(opts.ResourcesPath)
}

func sliceAbsPath(list stringSliceFlag) stringSliceFlag {
	var res stringSliceFlag
	for _, value := range list {
		_ = res.Set(getAbsPath(value))
	}
	return res
}

func getArgSlice(defVal stringSliceFlag, args []string, argName string) stringSliceFlag {
	if len(defVal) == 0 {
		var s stringSliceFlag
		val := getArgValue("", args, argName)
		if val != "" {
			_ = s.Set(val)
		}
		return s
	}
	return defVal
}

func getAbsPath(pathOrFile string) string {
	s := strings.Trim(pathOrFile, " ")
	if s != "" {
		s, _ = filepath.Abs(s)
	}
	return s
}

func getArgValue(defValue string, args []string, argName string) string {
	if defValue != "" {
		return defValue
	}
	argName = strings.ToLower(argName)
	for i, s := range args {
		if s == argName || s == argName+"=" {
			return args[i+1]
		} else if strings.Contains(s, "=") {
			idx := strings.Index(s, "=")
			name := s[:idx]
			if name == argName {
				return strings.Trim(s[idx+1:], "")
			}
		}
	}
	return ""
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
