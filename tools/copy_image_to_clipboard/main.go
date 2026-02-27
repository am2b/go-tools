package main

import (
    "flag"
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

func escapeForAppleScript(path string) string {
    path = strings.ReplaceAll(path, "\\", "\\\\")
    path = strings.ReplaceAll(path, "\"", "\\\"")
    return path
}

func getAppleScriptImageType(ext string) (string, error) {
    switch strings.ToLower(ext) {
    case ".png":
        return "«class PNGf»", nil
    case ".jpg", ".jpeg":
        return "«class JPEG»", nil
    default:
        return "", fmt.Errorf("不支持的图片格式:%s", ext)
    }
}

func main() {
    // 定义-h选项
    helpFlag := flag.Bool("h", false, "Show usage information")

    // 解析命令行参数
    flag.Parse()

    // 如果-h被设置,打印帮助信息
    if *helpFlag {
        fmt.Println("copy指定的一张图片到剪贴板")
        fmt.Println("Usage:")
        fmt.Println("Usage: copy_image_to_clipboard <image-path>")
        fmt.Println("Options:")
        fmt.Println("-h:Show this help message")
        return
    }

    // 检查路径参数
    if len(flag.Args()) < 1 {
        fmt.Fprintf(os.Stderr, "用法:%s <image-path>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }

    imagePath := flag.Arg(0)
    if _, err := os.Stat(imagePath); os.IsNotExist(err) {
        fmt.Fprintf(os.Stderr, "文件不存在:%s\n", imagePath)
        os.Exit(1)
    }

    absPath, err := filepath.Abs(imagePath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "获取绝对路径失败:%v\n", err)
        os.Exit(1)
    }

    imageType, err := getAppleScriptImageType(filepath.Ext(absPath))
    if err != nil {
        fmt.Fprintf(os.Stderr, "%v\n", err)
        os.Exit(1)
    }

    escapedPath := escapeForAppleScript(absPath)

    script := fmt.Sprintf(`
    set the clipboard to (read (POSIX file "%s") as %s)
    `, escapedPath, imageType)

    cmd := exec.Command("osascript", "-e", script)

    output, err := cmd.CombinedOutput()
    if err != nil {
        fmt.Fprintf(os.Stderr, "写入剪贴板失败:%v\n%s\n", err, output)
        os.Exit(1)
    }

    fmt.Println("图片已复制到剪贴板:", absPath)
}
