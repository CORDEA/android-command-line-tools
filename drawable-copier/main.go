/*
 * Copyright 2016 Yoshihiro Tanaka
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * Author: Yoshihiro Tanaka <contact@cordea.jp>
 * date  : 2016-05-29
 */

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
    isRemoveSource = flag.Bool("r", false, "If True, remove the source directory when all of process was successful.")
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

func copyFile(src, trg string, isOverwrite bool) bool {
    if !isOverwrite {
        if _, err := os.Stat(trg); err == nil {
            log.Println(trg + " already exists. Not overwrite.")
            return true
        }
    }

    srcFile, err := os.Open(src)
    if err != nil {
        log.Println(err)
        return false
    }
    defer srcFile.Close()

    trgFile, err := os.Create(trg)
    if err != nil {
        log.Println(err)
        return false
    }
    defer trgFile.Close()
    if _, err := io.Copy(trgFile, srcFile); err != nil {
        os.Remove(trg)
        log.Println(err)
        return false
    }
    log.Println("Copied " + trg)
    return true
}

func isDrawableDir(dirname string) bool {
    return strings.Contains(dirname, DrawableDirPrefix)
}

func copyFiles(src, trg string, isOverwrite bool) bool {
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

    allSucceed := true

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
                    allSucceed = copyFile(sfPath, tfPath, isOverwrite)
                }
            }
        }
    }
    return allSucceed
}

func main() {
    flag.Parse()
    if flag.NArg() < 2 {
        log.Fatalln("Required argument is missing.")
    }

    src := flag.Arg(0)
    trg := flag.Arg(1)
    allSucceed := copyFiles(src, trg, *isOverwrite)
    if *isRemoveSource && allSucceed {
        if err := os.RemoveAll(src); err != nil {
            log.Println(err)
        }
    }
}
