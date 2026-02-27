package main

import (
    "flag"
    "fmt"
    "github.com/BurntSushi/toml"
    "github.com/am2b/go-tools/pkg/pathutil"
    "log"
    "math/rand"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
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
        fmt.Println("读取helper.toml,copy指定目录下的一张图片到剪贴板,然后发送'粘贴'消息给hammerspoon")
        fmt.Println("Usage:")
        fmt.Println("copy_image_for_mail [-h]")
        fmt.Println("Options:")
        fmt.Println("-h:Show this help message")
        return
    }

    home, _ := os.UserHomeDir()
    tomlPath := filepath.Join(home, ".config", "mail-helper", "helper.toml")

    type helper struct {
        ResDir string `toml:"res_dir"`
        Images string `toml:"images"`
    }

    var helperIns helper
    if _, err := toml.DecodeFile(tomlPath, &helperIns); err != nil {
        log.Fatal(err)
    }

    imagePath := filepath.Join(helperIns.ResDir, helperIns.Images)
    imagePath = pathutil.ExpandPath(imagePath)

    entries, err := os.ReadDir(imagePath)
    if err != nil {
        log.Fatal(err)
    }

    imageFileNameSlice := make([]string, 0, len(entries))

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }

        imageFileName := strings.ToLower(entry.Name())
        if strings.HasSuffix(imageFileName, ".jpg") ||
            strings.HasSuffix(imageFileName, ".jpeg") ||
            strings.HasSuffix(imageFileName, ".png") {
            imageFileNameSlice = append(imageFileNameSlice, imageFileName)
        }
    }

    if len(imageFileNameSlice) == 0 {
        log.Fatal("没有图片了")
    }

    rand.Seed(time.Now().UnixNano())
    randomIndex := rand.Intn(len(imageFileNameSlice))
    randomImageFileName := imageFileNameSlice[randomIndex]
    randomImageFilePath := filepath.Join(imagePath, randomImageFileName)

    if _, err := os.Stat(randomImageFilePath); os.IsNotExist(err) {
        log.Fatalf("文件不存在:%s\n", randomImageFilePath)
    }

    imageType, err := getAppleScriptImageType(filepath.Ext(randomImageFilePath))
    if err != nil {
        log.Fatal(err)
    }

    escapedPath := escapeForAppleScript(randomImageFilePath)

    script := fmt.Sprintf(`
    set the clipboard to (read (POSIX file "%s") as %s)
    `, escapedPath, imageType)

    cmd := exec.Command("osascript", "-e", script)

    output, err := cmd.CombinedOutput()
    if err != nil {
        log.Fatalf("写入剪贴板失败:%v\n%s\n", err, output)
    }

    os.Remove(randomImageFilePath)

    //传递消息给hammerspoon
    fmt.Println("cmd v")
}
