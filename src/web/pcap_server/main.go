// main
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"web/pcap_server/asset"
	"web/util"

	"github.com/elazarl/go-bindata-assetfs"
	"golang.org/x/net/websocket"
)

//添加设置路由函数
func Set(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("set.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		error, result := util.CheckFormat(r)
		if !result {
			fmt.Fprintf(w, error.Error())
			return
		}
		set := util.InitPacketGraspingSet(r)
		os.MkdirAll("config", 0777)
		util.InitConfig("config/set.config")                     //初始化配置文件，如果不存在就新建配置文件
		config, err := util.ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
		if util.CheckRepath(set, config) {
			fmt.Fprintf(w, "不同的设置不能保存报文到同一个文件上")
			return
		}
		max_id := config.MaxId
		err = util.WriteToConfig(set, config, max_id) //往配置文件写配置信息
		if err != nil {
			fmt.Fprintf(w, err.Error())
		} else {
			fmt.Fprintf(w, "保存配置成功")
		}

	}
}

//修改配置路由
func Update(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		error, result := util.CheckFormat(r)
		if !result {
			fmt.Fprintf(w, error.Error())
			return
		}
		set := util.InitPacketGraspingSet(r)
		util.InitConfig("config/set.config")                     //初始化配置文件，如果不存在就新建配置文件
		config, err := util.ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
		set_id := r.Form.Get("set_id")
		id, _ := strconv.Atoi(set_id)
		set.SetId = id
		exit, _ := util.FoundConfigById(id)
		if !exit {
			err = errors.New("要更新的设置在配置文件里不存在")
			fmt.Fprintf(w, "更新失败，失败原因:"+err.Error())
			return
		}
		is_repeat := util.CheckFileUpdateRepeat(set, config.Sets) //检查更新后文件路径是否有重复
		if is_repeat {
			fmt.Fprintf(w, "更新失败，失败原因:不同的设置不能保存报文到同一个文件上")
			return
		}
		success, _ := util.UpdateConfig(set, config)
		if success {
			fmt.Fprintf(w, "更新配置成功")
			return
		} else {
			fmt.Fprintf(w, "更新配置失败")
			return
		}
	}
}

//开始抓包
func Start(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		id := r.Form["id"]
		for j := 0; j < len(id); j++ {
			ids, error := strconv.Atoi(id[j])
			if error != nil {
				err := "ID为" + strconv.Itoa(ids) + "抓包配置错误，id必须是数字;"
				fmt.Println(err)
				continue
			}
			exit, _ := util.FoundConfigById(ids)

			if !exit {
				err := "ID为" + strconv.Itoa(ids) + "的抓包线程启动失败，此id不存在;"
				fmt.Println(err)
				continue
			} else if util.CheckPGhold(ids) {
				err := "ID为" + strconv.Itoa(ids) + "的抓包线程启动失败，此id正在抓包;"
				fmt.Println(err)
				continue
			} else if util.IsAliveAndStoped(ids) {
				for i := 0; i < len(util.Tests); i++ {
					if util.Tests[i].SetId == ids && !util.CheckPGhold(ids) {
						util.Tests[i].Set(ids)
						go util.Tests[i].PGStart()
						str1 := "ID为" + strconv.Itoa(util.Tests[i].SetId) + "的抓包线程启动成功"
						fmt.Println(str1)
					}
				}
			} else {
				var Test util.PacketGrasping
				Test.Set(ids)
				util.Tests = append(util.Tests, Test)
				k := len(util.Tests) - 1
				go util.Tests[k].PGStart()
				str1 := "ID为" + strconv.Itoa(util.Tests[k].SetId) + "的抓包线程启动成功"
				fmt.Println(str1)
			}
		}
		var result util.Result
		var results util.Results
		var resultArray []util.Result = make([]util.Result, 0)
		for i := 0; i < len(util.Tests); i++ {
			result.PGState = util.Tests[i].PGState
			result.SetId = util.Tests[i].SetId
			resultArray = append(resultArray, result)
		}
		results.Message = resultArray
		message, errMarshl := json.Marshal(results)
		//log.Println(string(message))
		if errMarshl != nil {
			log.Println("结构体转json错误")
		}

		fmt.Fprintf(w, "%s", string(message))
	}
}

//每隔1s发抓包状态给抓包状态界面
func Echo(ws *websocket.Conn) {
	var err error
	var result util.Result
	var results util.Results
	for {
		var resultArray []util.Result = make([]util.Result, 0)
		time.Sleep(1 * time.Second)
		for i := 0; i < len(util.Tests); i++ {
			if !util.Tests[i].Display {
				continue
			}
			result.PGState = util.Tests[i].PGState
			result.SetId = util.Tests[i].SetId
			resultArray = append(resultArray, result)
		}
		results.Message = resultArray
		message, errMarshl := json.Marshal(results)
		//log.Println(string(message))
		if errMarshl != nil {
			log.Println("结构体转json错误")
			break
		}
		if err = websocket.Message.Send(ws, string(message)); err != nil {
			break
		}
	}
}
func state(w http.ResponseWriter, r *http.Request) {
	var result util.Result
	var results util.Results
	var resultArray []util.Result = make([]util.Result, 0)
	for i := 0; i < len(util.Tests); i++ {
		if !util.Tests[i].Display {
			continue
		}
		result.PGState = util.Tests[i].PGState
		result.SetId = util.Tests[i].SetId
		resultArray = append(resultArray, result)
	}
	results.Message = resultArray
	message, errMarshl := json.Marshal(results)
	//log.Println(string(message))
	if errMarshl != nil {
		log.Println("结构体转json错误")
		return
	}
	fmt.Fprintf(w, "%s", string(message))
}
func refresh(w http.ResponseWriter, r *http.Request) {
	util.Refresh()
	var result util.Result
	var results util.Results
	var resultArray []util.Result = make([]util.Result, 0)
	for i := 0; i < len(util.Tests); i++ {
		if util.Tests[i].Finished {
			util.Tests[i].Display = false
			continue
		}
	}
	for i := 0; i < len(util.Tests); i++ {
		if !util.Tests[i].Display {
			continue
		}
		result.PGState = util.Tests[i].PGState
		result.SetId = util.Tests[i].SetId
		resultArray = append(resultArray, result)
	}
	results.Message = resultArray
	message, errMarshl := json.Marshal(results)
	//log.Println(string(message))
	if errMarshl != nil {
		log.Println("结构体转json错误")
		return
	}
	fmt.Fprintf(w, "%s", string(message))
}

//查看设置路由函数
func Views(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sets := make([]util.PacketGraspingSet, 0)
		config, error := util.ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
		if error != nil {
			log.Println(error)
		} else {
			sets = config.Sets
		}
		message, _ := json.Marshal(sets)
		fmt.Fprintf(w, "%s", string(message))
	}
}

//首页
func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data, _ := asset.Asset("static/index.html")
		w.Header().Add("Content-Type", "text/html; charset=UTF-8")
		w.Write(data)
	}
}

//判断状态
func Rules(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sets := make([]util.PacketGraspingSet, 0)
		config, error := util.ConfigUnmarshal("config/set.config") //读取配置文件，获得结构体
		if error != nil {
			log.Println(error)
		} else {
			sets = config.Sets
		}
		var rules []util.Rule
		for _, set := range sets {
			var rule util.Rule
			rule.Set = set
			rule.Flag = util.CheckPGhold(set.SetId)
			rules = append(rules, rule)
		}
		message, _ := json.Marshal(rules)
		fmt.Fprintf(w, "%s", string(message))
	}
}

//获得某条设置
func View(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		id := r.Form.Get("id")
		set_id, error := strconv.Atoi(id)
		if error != nil {
			log.Println(error)
			return
		}
		exit, set := util.FoundConfigById(set_id)
		if exit {
			output, _ := json.Marshal(set)
			fmt.Fprintf(w, "%s", string(output))
		}
	}
}

//路由删除函数
func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		id := r.Form.Get("id")
		if id == "" {
			error := errors.New("删除id为空")
			fmt.Fprintf(w, "删除失败，失败原因:"+error.Error())
		} else {
			set_id, error := strconv.Atoi(id)
			if error != nil {
				fmt.Fprintf(w, "删除失败，失败原因:"+error.Error())
			} else {
				result, error := util.DeleteById(set_id)
				if result == true {
					fmt.Fprintf(w, "删除成功")
				} else {
					fmt.Fprintf(w, "删除失败，失败原因:"+error.Error())
				}
			}
		}
	}
}

//停止抓包路由
func Stop(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		r.ParseForm()
		id := r.Form.Get("id")
		if id == "" {
			fmt.Fprintf(w, "停止抓包失败")
		} else {
			set_id, error := strconv.Atoi(id)
			if error != nil {
				fmt.Fprintf(w, "停止抓包失败")
			} else {
				result := util.StopPGById(set_id)
				if result == true {
					fmt.Fprintf(w, "停止抓包成功")
				} else {
					fmt.Fprintf(w, "停止抓包失败")
				}
			}
		}
	}
}

//获取ip路由
func Ip(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		fmt.Fprintf(w, r.Host)
	}
}

func main() {
	go util.Command()
	os.MkdirAll("config", 0777)
	fs := assetfs.AssetFS{
		Asset:     asset.Asset,
		AssetDir:  asset.AssetDir,
		AssetInfo: asset.AssetInfo,
	}
	http.Handle("/static/", http.StripPrefix("/", http.FileServer(&fs)))
	http.HandleFunc("/set", Set)       //设置访问的路由，添加设置
	http.HandleFunc("/views", Views)   //查看所有设置
	http.HandleFunc("/view", View)     //查看某条设置
	http.HandleFunc("/delete", Delete) //删除设置
	http.HandleFunc("/update", Update) //修改设置
	http.HandleFunc("/", Index)        //首页
	http.HandleFunc("/start", Start)   //开始抓包
	http.HandleFunc("/stop", Stop)     //停止抓包
	http.HandleFunc("/rule", Rules)    //判断状态
	http.HandleFunc("/state", state)
	http.HandleFunc("/refresh", refresh)
	http.Handle("/result", websocket.Handler(Echo)) //抓包状态
	http.HandleFunc("/ip", Ip)                      //获得ip
	err := http.ListenAndServe(":9090", nil)        //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
