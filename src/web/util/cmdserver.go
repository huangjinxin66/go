// cmdserver
package util

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/modood/table"
)

//命令行控制
func Command() {
	comm := flag.String("c", "", "指令类型(view查看所有配置，start开始抓包，set添加配置)")
	file_max_size := flag.Float64("size", 20.0, "单个文件最大大小(单位MB)")
	file_name := flag.String("name", "", "文件名")
	port := flag.String("port", "", "本地监听端口号")
	heart_port := flag.String("heart_port", "", "心跳端口")
	during_time := flag.String("t", "", "抓包持续时间(单位min)")
	nativeIp := flag.String("nativeIp", "", "本地ip")
	remoteIp := flag.String("remoteIp", "", "远端ip正则表达式")
	path := flag.String("path", "packet", "文件保存路径")
	startid := flag.String("id", "", "输入抓包对应的配置id进行抓包")
	flag.Parse()
	switch *comm {
	case "set":
		if res, _ := CheckDuringTimeFormat(*during_time); !res {
			fmt.Println("抓包时间必须为正整数")
			return
		}
		if *file_name == "" {
			fmt.Println("配置文件名不能为空")
			return
		}
		if *file_max_size <= 0 {
			fmt.Println("文件大小不能为负数")
			return
		}
		if result, error := CheckPortFormat(*port); !result {
			fmt.Println(error)
			return
		}
		if result, error := CheckHeartPortFormat(*heart_port); !result {
			fmt.Println(error)
			return
		}
		if result, error := CheckIp(*nativeIp, *port); !result {
			fmt.Println(error)
			return
		}
		if result, error := CheckFileNameFormat(*file_name); !result {
			fmt.Println(error)
			return
		}
		if *heart_port == *port {
			fmt.Println("本地端口号和心跳端口号不能一样")
			return
		}
		var set PacketGraspingSet
		set.FileMaxSize = *file_max_size
		set.FileName = *file_name
		set.HeartPort, _ = strconv.Atoi(*heart_port)
		set.NativeServerPort, _ = strconv.Atoi(*port)
		set.NativeIp = *nativeIp
		set.RemoteIp = *remoteIp
		set.PacketHoldingTime = *during_time
		set.SavePath = *path
		InitConfig("config/set.config")                     //初始化配置文件，如果不存在就新建配置文件
		config, err := ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
		if CheckRepath(set, config) {
			fmt.Println("不同的设置不能保存报文到同一个文件上")
			os.Exit(1)
		}
		max_id := config.MaxId
		err = WriteToConfig(set, config, max_id) //往配置文件写配置信息
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("保存配置成功")
		}

	case "view":
		config, error := ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
		if error != nil {
			fmt.Println("打开配置文件错误")
		} else if len(config.Sets) == 0 {
			fmt.Println("配置文件为空")
		} else {
			var sets []PacketGraspingSet
			for _, set := range config.Sets {
				sets = append(sets, set)
			}
			table.Output(sets)
		}
	case "start":
		if *startid == "" {
			fmt.Println("请输入抓包配置id")
		} else {
			ids := StrtoArrs(*startid)
			for j := 0; j < len(ids); j++ {
				exit, _ := FoundConfigById(ids[j])
				if !exit {
					err := "ID为" + strconv.Itoa(ids[j]) + "的抓包线程启动失败，此id不存在"
					fmt.Println(err)
					continue
				} else if CheckPGhold(ids[j]) {
					err := "ID为" + strconv.Itoa(ids[j]) + "的抓包线程启动失败，此id正在抓包"
					fmt.Println(err)
					continue
				} else {
					var Test PacketGrasping
					Test.Set(ids[j])
					Tests = append(Tests, Test)
					k := len(Tests) - 1
					go Tests[k].PGStart()
					str1 := "ID为" + strconv.Itoa(Tests[k].SetId) + "的抓包线程启动成功"
					fmt.Println(str1)
				}
			}
		}
	}
}

func CheckPGhold(id int) bool {
	for i := 0; i < len(Tests); i++ {
		if Tests[i].SetId == id && Tests[i].Finished == false {
			return true
		}
	}
	return false
}
