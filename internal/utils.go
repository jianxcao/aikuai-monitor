package aikuaimonitor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sort"

	"github.com/fatih/color"
	"github.com/gosuri/uilive"
)

func ClearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	case "darwin", "linux":
		fmt.Printf("\r")
		os.Stdout.Write([]byte{0x1B, 0x5B, 0x33, 0x3B, 0x4A, 0x1B, 0x5B, 0x48, 0x1B, 0x5B, 0x32, 0x4A})
		// cmd := exec.Command("clear")
		// cmd.Stdout = os.Stdout
		// cmd.Run()
	default:
		fmt.Println("Unsupported platform")
	}
}

func FormatSize(size float64) string {
	units := []string{"B", "KB", "MB", "GB", "TB"}
	unitIndex := 0

	for size >= 1024 && unitIndex < len(units)-1 {
		size /= 1024
		unitIndex++
	}
	return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

func SortMapKeys(data interface{}) []string {
	keys := make([]string, 0)
	switch res := data.(type) {
	case map[string][][]string:
		for key := range res {
			keys = append(keys, key)
		}
	case map[string]map[string][][]string:
		for key := range res {
			keys = append(keys, key)
		}
	}
	sort.Strings(keys)
	return keys
}

func PrintGrid(data [][]string, writer *uilive.Writer) string {
	// 计算每列的最大宽度
	colWidths := make([]int, len(data[0]))
	for _, row := range data {
		for colIndex, cell := range row {
			if len(cell) > colWidths[colIndex] {
				colWidths[colIndex] = len(cell)
			}
		}
	}

	headerStyle := color.New(color.FgWhite, color.Bold)
	defaultStyle := color.New(color.FgGreen)
	var result = ""
	// 输出表格
	for rowIndex, row := range data {
		var rowData = ""
		var style *color.Color
		// 根据列索引选择样式
		switch rowIndex {
		case 0: // 表头
			style = headerStyle
		default: // ID 列
			style = defaultStyle
		}
		for colIndex, cell := range row {
			// 根据列宽度进行格式化和对齐
			format := fmt.Sprintf("%%-%ds", colWidths[colIndex])
			cell = fmt.Sprintf(format, cell)
			rowData += cell + "  "
		}
		result += style.SprintlnFunc()(rowData)
	}
	return result
}

func PrintData(res map[string]map[string][][]string, writer *uilive.Writer) {
	keys := SortMapKeys(res)
	var str = ""
	for _, key := range keys {
		data := res[key]
		str += color.New(color.FgWhite, color.Bold).SprintlnFunc()("主机名: " + key)
		gkeys := SortMapKeys(data)
		for _, gKey := range gkeys {
			str += PrintGrid(data[gKey], writer)
			str += color.New(color.FgBlue).SprintlnFunc()("--------------------------------------------")
		}
	}
	fmt.Fprint(writer, str)
}

// 浅层合并map
func MergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	mergedMap := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			mergedMap[k] = v
		}
	}
	return mergedMap
}
