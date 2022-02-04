package common

import (
	"math/rand"
	"os"
	"time"
)

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func IsFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func ByteToBytes(bt byte) (bytes []byte) {
	bytes = append(bytes, bt)
	return bytes
}

//生成count个[start,end)结束的不重复的随机数
func GenerateRandomNumber(start int, end int, count int) []int {
	//范围检查
	if end < start || (end-start) < count {
		return nil
	}

	//存放结果的slice
	nums := make([]int, 0)
	//随机数生成器，加入时间戳保证每次生成的随机数不一样
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(nums) < count {
		//生成随机数
		num := r.Intn((end - start)) + start
		//查重
		exist := false
		for _, v := range nums {
			if v == num {
				exist = true
				break
			}
		}
		if !exist {
			nums = append(nums, num)
		}
	}
	return nums
}
