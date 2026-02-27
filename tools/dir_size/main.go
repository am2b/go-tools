package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

func usage() {
    fmt.Println("计算参数所给定目录的大小,默认为当前目录")
    fmt.Println("Usage:")
    fmt.Println("dir_size [directory_path]")
    fmt.Println("Options:")
    fmt.Println("-h:Show this help message")
}

func VerifyDir(dir string) (string, error) {
    //转换为绝对路径
    dir, err := filepath.Abs(dir)
    if err != nil {
        return "", err
    }

    //stat
    info, err := os.Stat(dir)
    if err != nil {
        if os.IsNotExist(err) {
            return "", errors.New("directory does not exist")
        }
        return "", err
    }

    if !info.IsDir() {
        return "", errors.New("path is not a directory")
    }

    return dir, nil
}

func getFileSize(filePath string, wg *sync.WaitGroup, ch chan int64) {
    defer wg.Done()
    fi, err := os.Stat(filePath)
    if err != nil {
        ch <- 0
        return
    }
    if !fi.IsDir() {
        ch <- fi.Size()
    } else {
        ch <- 0
    }
}

func getDirSize(path string) (int64, error) {
    var totalSize int64
    var wg sync.WaitGroup
    ch := make(chan int64)

    err := filepath.Walk(path, func(file string, fi os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        wg.Add(1)
        go getFileSize(file, &wg, ch)
        return nil
    })

    if err != nil {
        return 0, err
    }

    go func() {
        wg.Wait()
        close(ch)
    }()

    for size := range ch {
        totalSize += size
    }

    return totalSize, nil
}

func convertSize(sizeInBytes int64) string {
    // MiB = 1024 进制,MB = 1000 进制
    units := []string{"B", "KiB", "MiB", "GiB", "TiB"}
    var unitIndex int
    size := float64(sizeInBytes)

    for size >= 1024.0 && unitIndex < len(units)-1 {
        size /= 1024.0
        unitIndex++
    }

    return fmt.Sprintf("%.2f %s", size, units[unitIndex])
}

func main() {
    // 定义-h选项
    helpFlag := flag.Bool("h", false, "Show usage information")

    // 解析命令行参数
    flag.Parse()

    // 如果-h被设置,打印帮助信息
    if *helpFlag {
        usage()
        return
    }

    // 获取目录路径
    path := "."
    if len(flag.Args()) > 0 {
        path = flag.Arg(0)
    }

    // 处理路径中的 ~，将其替换为用户的主目录
    homeDir, err := os.UserHomeDir()
    if err != nil {
        fmt.Println("Error getting home directory:", err)
        return
    }

    // 如果路径中包含 ~，替换为实际的用户主目录
    if strings.HasPrefix(path, "~") {
        path = strings.Replace(path, "~", homeDir, 1)
    }

    // 验证目录
    path, err = VerifyDir(path)
    if err != nil {
        panic(err)
    }

    // 计算目录大小
    totalSize, err := getDirSize(path)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // 输出结果
    fmt.Println(convertSize(totalSize))
}
