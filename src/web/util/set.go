// set
package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	//"github.com/google/gopacket/pcap"
)

//字符串转成数组
func StrtoArrs(str string) []int {
	str = str + ","
	arrs := make([]int, 0)
	var s []byte
	st := []byte(str)
	for i := 0; i < len(st); i++ {
		if st[i] >= 48 && st[i] <= 57 {
			s = append(s, st[i])
		} else if s != nil {
			v := string(s)
			a, _ := strconv.Atoi(v)
			arrs = append(arrs, a)
			s = nil
		}
	}
	return arrs
}

//初始化文件
func InitFolder(path string) {
	os.MkdirAll(path, 0777)
}

//更新配置
func UpdateConfig(set PacketGraspingSet, config Config) (bool, error) {
	sets := config.Sets
	success := false
	for i, s := range sets {
		if s.SetId == set.SetId {
			config.Sets[i].NativeServerPort = set.NativeServerPort
			config.Sets[i].HeartPort = set.HeartPort
			config.Sets[i].PacketHoldingTime = set.PacketHoldingTime
			config.Sets[i].FileMaxSize = set.FileMaxSize
			config.Sets[i].SavePath = set.SavePath
			config.Sets[i].FileName = set.FileName
			config.Sets[i].RemoteIp = set.RemoteIp
			config.Sets[i].NativeIp = set.NativeIp
			success = true
			break
		}
	}
	output, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		err = ioutil.WriteFile("config/set.config", output, os.ModeAppend)
	}
	return success, err
}

//检查是不是所有的抓包都结束了
func IsTasksFinished() bool {
	for i := 0; i < len(Tests); i++ {
		if Tests[i].Finished == false {
			return false
		}
	}
	return true
}

//根据id停止抓包
func StopPGById(id int) bool {
	for i := 0; i < len(Tests); i++ {
		if Tests[i].Finished == false && id == Tests[i].SetId {
			Tests[i].PGEnd()
			return true
		}
	}
	return false
}

//检查此id对象是否存在且状态为结束
func IsAliveAndStoped(ids int) bool {
	for i := 0; i < len(Tests); i++ {
		if Tests[i].SetId == ids && Tests[i].Finished {
			return true
		}
	}
	return false
}

//检查文件路径是否重复
func CheckFileUpdateRepeat(set PacketGraspingSet, sets []PacketGraspingSet) bool {
	for _, s := range sets {
		if set.SetId == s.SetId {
			continue
		}
		if s.SavePath == set.SavePath &&
			s.FileName == set.FileName {
			return true
		}
	}
	return false

}

//检查端口更新是否重复
func CheckPortUpdateRepeat(set PacketGraspingSet, sets []PacketGraspingSet) bool {
	for _, s := range sets {
		if set.SetId == s.SetId {
			continue
		}
		if s.NativeServerPort == set.NativeServerPort || s.NativeServerPort == set.HeartPort || set.NativeServerPort == s.HeartPort || set.HeartPort == s.HeartPort {
			return true
		}
	}
	return false

}

//将配置信息写进配置文件里
func WriteToConfig(set PacketGraspingSet, config Config, max_id int) error {
	set.SetId = max_id + 1
	config.Sets = append(config.Sets, set)
	config.MaxId = max_id + 1
	output, err := json.Marshal(config)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	} else {
		err = ioutil.WriteFile("config/set.config", output, os.ModeAppend)
	}
	if err != nil {
		err = errors.New("文件路径config/set.config打不开，请手动创建config文件夹以及在config文件夹里新建set.config文件")
	}
	return err
}

//根据id删除设置
func DeleteById(id int) (bool, error) {
	config, error := ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
	if error != nil {
		return false, error
	} else {
		sets := config.Sets
		for i, set := range sets {
			if set.SetId == id {
				config.Sets = Remove(config.Sets, i)
				break
			}
		}
		output, error := json.Marshal(config)
		if error != nil {
			return false, error
		} else {
			error = ioutil.WriteFile("config/set.config", output, os.ModeAppend)
			if error != nil {
				return false, error
			} else {
				return true, nil
			}
		}

	}
}

//初始化配置文件
func InitConfig(path string) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		file, err = os.Create(path)
	}
}

//获取所有id
func GetAllId() ([]int, error) {
	var ids []int
	config, err := ConfigUnmarshal("config/set.config")
	if err != nil {
		return nil, err
	} else {
		sets := config.Sets
		for i := 0; i < len(sets); i++ {
			var id int
			id = sets[i].SetId
			ids = append(ids, id)
		}
		return ids, nil
	}
}

//根据id找设置
func FoundConfigById(set_id int) (bool, PacketGraspingSet) {
	var s PacketGraspingSet
	config, error := ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
	if error != nil {
		return false, s
	} else {
		sets := config.Sets
		for _, set := range sets {
			if set.SetId == set_id {
				s = set
				return true, s
			}
		}
		return false, s
	}
}

//解析配置文件，得到Config结构体
func ConfigUnmarshal(path string) (Config, error) {
	var config Config
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(data, &config)
	return config, err
} //检查配置信息的格式
func CheckFormat(r *http.Request) (error, bool) {
	port := r.Form.Get("port")
	during_time := r.Form.Get("during_time")
	file_max_size := r.Form.Get("file_max_size")
	file_name := r.Form.Get("file_name")
	heart_port := r.Form.Get("heart_port")
	if port == heart_port {
		error := errors.New("心跳端口和本地端口不能一样")
		return error, false
	}
	if result, error := CheckIp(r.Form.Get("nativeIp"), port); !result {
		return error, result
	}
	if result, error := CheckPortFormat(port); !result {
		return error, result
	}
	if result, error := CheckHeartPortFormat(heart_port); !result {
		return error, result
	}
	if port == heart_port {
		error := errors.New("心跳端口不能和本地端口一样")
		return error, false
	}
	if result, error := CheckDuringTimeFormat(during_time); !result {
		return error, result
	}
	if result, error := CheckFileMaxSizeFormat(file_max_size); !result {
		return error, result
	}
	if result, error := CheckFileNameFormat(file_name); !result {
		return error, result
	}
	return nil, true

}

//检查传过来的心跳端口
func CheckHeartPortFormat(heart_port string) (bool, error) {
	if heart_port == "" {
		error := errors.New("心跳端口不能为空")
		return false, error
	}
	if m, _ := regexp.MatchString("^[0-9]+$", heart_port); !m {
		error := errors.New("心跳端口必须为正整数")
		return false, error
	}
	port_value, _ := strconv.Atoi(heart_port)
	if port_value > 65535 {
		error := errors.New("心跳端口值大于65535")
		return false, error
	}
	return true, nil
}

//检查传过来的端口
func CheckPortFormat(port string) (bool, error) {
	if port == "" {
		error := errors.New("本地端口不能为空")
		return false, error
	}
	if m, _ := regexp.MatchString("^[0-9]+$", port); !m {
		error := errors.New("本地端口必须为正整数")
		return false, error
	}
	port_value, _ := strconv.Atoi(port)
	if port_value > 65535 {
		error := errors.New("本地端口值不能大于65535")
		return false, error
	}
	return true, nil
}

//检查传过来的抓包时间
func CheckDuringTimeFormat(during_time string) (bool, error) {
	if during_time == "" {
		return true, nil
	}
	isNum, _ := regexp.MatchString("^[0-9]+$", during_time)
	if !isNum {
		error := errors.New("抓包时间不是正整数")
		return false, error
	}
	return true, nil
}

//检查传过来的文件大小限制
func CheckFileMaxSizeFormat(file_max_size string) (bool, error) {
	if file_max_size == "" {
		return true, nil
	}
	isNum, _ := regexp.MatchString("^[0-9]+$", file_max_size)
	isDecimal, _ := regexp.MatchString("^[0-9]+\\.[0-9]+$", file_max_size)
	if !isNum && !isDecimal {
		error := errors.New("文件限制大小不是整数或小数")
		return false, error
	}
	return true, nil
}

//检查文件名是否为空
func CheckFileNameFormat(file_name string) (bool, error) {
	if file_name == "" {
		error := errors.New("没输入文件名")
		return false, error
	}
	s := []byte(file_name)
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			error := errors.New("文件名不能有斜杆")
			return false, error
		}
	}

	return true, nil
}

//检查文件保存是否重复
func CheckRepath(set PacketGraspingSet, config Config) bool {
	sets := config.Sets
	if sets == nil || len(sets) == 0 {
		return false
	} else {
		for _, s := range sets {
			if s.SavePath == set.SavePath &&
				s.FileName == set.FileName {
				return true
			}
		}
		return false
	}

}

//检查端口号是否已被其他配置占用
func CheckReport(set PacketGraspingSet, config Config) bool {
	sets := config.Sets
	if sets == nil || len(sets) == 0 {
		return false
	} else {
		for _, s := range sets {
			if s.NativeServerPort == set.NativeServerPort || s.NativeServerPort == set.HeartPort || set.NativeServerPort == s.HeartPort || set.HeartPort == s.HeartPort {
				return true
			}
		}
		return false
	}

}

//删除数组的一个元素
func Remove(s []PacketGraspingSet, i int) []PacketGraspingSet {
	return append(s[:i], s[i+1:]...)
}
