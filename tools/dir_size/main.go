package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
)

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
    // 获取命令行参数
    if len(os.Args) < 2 {
        fmt.Println("Usage: go run dir_size.go <directory_path>")
        return
    }

    // 获取目录路径
    path := os.Args[1]

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

    // 计算目录大小
    totalSize, err := getDirSize(path)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    // 输出结果
    fmt.Println("Total directory size:", convertSize(totalSize))
}
