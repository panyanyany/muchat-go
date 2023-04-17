package utf8_util

func Substr(str string, begin, length int) (substr string) {
    // 将字符串的转换成[]rune
    rs := []rune(str)
    lth := len(rs)

    // 简单的越界判断
    if begin < 0 {
        begin = 0
    }
    if begin >= lth {
        begin = lth
    }
    end := begin + length
    if end > lth {
        end = lth
    }

    // 返回子串
    return string(rs[begin:end])
}
func Len(str string) int {
    return len([]rune(str))
}
func SplitByLen(str string, length int) (results []string) {
    results = []string{}

    cnt := 0
    var sub string
    for _, rune := range str {
        sub += string(rune)
        cnt++
        if cnt%length == 0 {
            results = append(results, sub)
            sub = ""
        }
    }
    if sub != "" {
        results = append(results, sub)
        sub = ""
    }
    return
}
func ReplaceAt(s string, i int, c rune) string {
    r := []rune(s)
    r[i] = c
    return string(r)
}
