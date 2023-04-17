package file_util

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
)

// OutputFileOptions ...
type OutputFileOptions struct {
    DirPerm    os.FileMode
    FilePerm   os.FileMode
    JSONPrefix string
    JSONIndent string
}

// FileExists checks if file exists, only for file, not for dir
func FileExists(path string) bool {
    info, err := os.Stat(path)

    if err != nil {
        return false
    }

    if info.IsDir() {
        return false
    }

    return true
}

// OutputFile auto creates file if not exists, it will try to detect the data type and
// auto output binary, string or json
func OutputFile(p string, data interface{}, options *OutputFileOptions) error {
    if options == nil {
        options = &OutputFileOptions{0775, 0664, "", "    "}
    }

    dir := filepath.Dir(p)
    os.MkdirAll(dir, 0755)

    var bin []byte

    switch t := data.(type) {
    case []byte:
        bin = t
    case string:
        bin = []byte(t)
    default:
        var err error
        bin, err = json.MarshalIndent(data, options.JSONPrefix, options.JSONIndent)

        if err != nil {
            return err
        }
    }

    return ioutil.WriteFile(p, bin, options.FilePerm)
}

func CleanFilename(name string) string {
    newName := ""
    charMap := map[int32]int32{
        '/':  '_',
        '\\': '_',
        ':':  '_',
        '*':  '_',
        '?':  '_',
        '"':  '_',
        '<':  '_',
        '>':  '_',
        '|':  '_',
    }

    for _, c := range name {
        repl, found := charMap[c]
        if found {
            newName += string(repl)
        } else {
            newName += string(c)
        }
    }
    return newName
}

func AppendFile(p string, data interface{}, options *OutputFileOptions) error {
    if options == nil {
        options = &OutputFileOptions{0775, 0664, "", "    "}
    }

    dir := filepath.Dir(p)
    os.MkdirAll(dir, 0755)

    var bin []byte

    switch t := data.(type) {
    case []byte:
        bin = t
    case string:
        bin = []byte(t)
    default:
        var err error
        bin, err = json.MarshalIndent(data, options.JSONPrefix, options.JSONIndent)

        if err != nil {
            return err
        }
    }

    fp, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
    if err != nil {
        err = fmt.Errorf("OpenFile: %w", err)
        return err
    }

    _, err = fp.Write(bin)
    if err != nil {
        err = fmt.Errorf("fp.Write: %w", err)
        return err
    }

    err = fp.Close()
    return err
}
