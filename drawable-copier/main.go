package main

import (
    "flag"
    "log"
    "os"
    "io"
    "io/ioutil"
    "path/filepath"
    "strings"
)

var (
    isOverwrite = flag.Bool("w", false, "If True, overwrite the file if the file already exists in the target directory.")
)

const (
    DrawableDirPrefix = "drawable-"
)

func getFiles(path string) []os.FileInfo {
    stat, err := os.Stat(path)
    if err != nil {
        log.Fatalln(err)
    }
    switch md := stat.Mode(); {
    case md.IsDir():
        files, err := ioutil.ReadDir(path)
        if err != nil {
            log.Fatalln(err)
        }
        return files
    }
    return []os.FileInfo{}
}

func copyFile(src, trg string, isOverwrite bool) {
    if !isOverwrite {
        if _, err := os.Stat(trg); err == nil {
            log.Println(trg + " already exists. Not overwrite.")
            return
        }
    }

    srcFile, err := os.Open(src)
    if err != nil {
        log.Println(trg + ": " + err.Error())
    }
    defer srcFile.Close()

    trgFile, err := os.Create(trg)
    if err != nil {
        log.Println(trg + ": " + err.Error())
    }
    defer trgFile.Close()
    if _, err := io.Copy(trgFile, srcFile); err != nil {
        log.Println(trg + ": " + err.Error())
    }
    log.Println("Copied " + trg)
}

func isDrawableDir(dirname string) bool {
    return strings.Contains(dirname, DrawableDirPrefix)
}

func copyFiles(src, trg string, isOverwrite bool) {
    srcs := getFiles(src)
    trgs := getFiles(trg)

    if len(trgs) == 0 {
        for _, sf := range srcs {
            if isDrawableDir(sf.Name()) {
                os.Mkdir(filepath.Join(trg, sf.Name()), 0700)
            }
        }
        trgs = getFiles(trg)
    }

    for _, sf := range srcs {
        if !isDrawableDir(sf.Name()) {
            continue
        }
        for _, tf := range trgs {
            if sf.Name() == tf.Name() {
                sImages := getFiles(filepath.Join(src, sf.Name()))
                for _, si := range sImages {
                    sfPath := filepath.Join(src, sf.Name(), si.Name())
                    tfPath := filepath.Join(trg, tf.Name(), si.Name())
                    copyFile(sfPath, tfPath, isOverwrite)
                }
            }
        }
    }
}

func main() {
    flag.Parse()
    if flag.NArg() < 2 {
        log.Fatalln("Required argument is missing.")
    }

    src := flag.Arg(0)
    trg := flag.Arg(1)
    copyFiles(src, trg, *isOverwrite)
}
