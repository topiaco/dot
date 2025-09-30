package dot

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	DIGIT      = "0123456789"
	MixedLower = "qwertyuiopasdfghjklzxcvbnm0123456789"
	MixedAll   = "QWERTYUIOPASDFGHJKLZXCVBNMqwertyuiopasdfghjklzxcvbnm0123456789"
)

// GenNumNo 生成纯数字编号
// len: 编号长度
// prefix: 前缀字符串
func GenNumNo(length int, prefix string) string {
	return generator(length, DIGIT, prefix)
}

// GenMixedNo 生成字母与数字混合的编号
// length: 编号长度
// prefix: 前缀字符串
// isLower: 是否只使用小写字母
func GenMixedNo(length int, prefix string, isLower bool) string {
	charset := MixedAll
	if isLower {
		charset = MixedLower
	}
	return generator(length, charset, prefix)
}

// GenDateNumNo 生成带日期前缀的纯数字编号
// length: 编号主体部分长度
// dtFmt: 时间格式（如 "060102150405" 对应 ymdHis）
// prefix: 自定义前缀
func GenDateNumNo(length int, dtFmt string, prefix string) string {
	currentTime := time.Now().Format(dtFmt)
	prefix += currentTime
	return generator(length, DIGIT, prefix)
}

// GenDateMixedNo 生成带日期前缀的字母和数字混合编号
// length: 编号主体部分长度
// dtFmt: 时间格式（如 "060102150405" 对应 ymdHis）
// prefix: 自定义前缀
// isLower: 是否只使用小写字母
func GenDateMixedNo(length int, dtFmt string, prefix string, isLower bool) string {
	currentTime := time.Now().Format(dtFmt)
	prefix += currentTime

	charset := MixedAll
	if isLower {
		charset = MixedLower
	}
	return generator(length, charset, prefix)
}

// generator 内部通用生成函数
// length: 需要生成的随机字符数
// chars: 可用字符集
// prefix: 前缀
func generator(length int, chars string, prefix string) string {
	if length <= 0 || chars == "" {
		return ""
	}

	// 扩展字符集以满足所需长度
	repeatCount := (length + len(chars) - 1) / len(chars) // 向上取整
	extendedChars := strings.Repeat(chars, repeatCount)

	// 洗牌字符顺序并截取指定长度
	shuffled := shuffleString(extendedChars)[:length]

	return fmt.Sprintf("%s%s", prefix, shuffled)
}

// shuffleString 将字符串中的字符打乱顺序
func shuffleString(s string) string {
	runes := []rune(s)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}
