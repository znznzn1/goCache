package Cache

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

const (
	B = 1 << (iota * 10)
	KB
	MB
	GB
	TB
	PB
)

func ParseSize(size string) (int64, string) {
	re, _ := regexp.Compile("[0-9]*")
	unit := string(re.ReplaceAll([]byte(size), []byte(""))) //将数字部分替换位数组
	num, _ := strconv.ParseInt(strings.Replace(size, unit, "", 1), 10, 64)
	unit = strings.ToUpper(unit)
	var byteNum int64 = 0
	switch unit {
	case "B":
		byteNum = num
	case "KB":
		byteNum = num * KB
	case "MB":
		byteNum = num * MB
	case "GB":
		byteNum = num * GB
	case "TB":
		byteNum = num * TB
	case "PB":
		byteNum = num * PB
	default:
		num = 0
		byteNum = 0
	}
	if num == 0 {
		log.Println("ParseSize 仅支持 B,KB,MB,GB,TB,PB为单位")
		num = 100
		byteNum = num * MB
		unit = "MB"
	}
	sizeStr := strconv.FormatInt(num, 10) + unit
	return byteNum, sizeStr
}

func GetValueSize(val interface{}) int64 {
	return int64(unsafe.Sizeof(val))
}
