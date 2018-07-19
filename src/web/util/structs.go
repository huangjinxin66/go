// structs
package util

import (
	"net"
	"os"

	"github.com/google/gopacket/pcapgo"
)

//抓包配置结构体
type Rule struct {
	Set  PacketGraspingSet
	Flag bool
}

//抓包配置结构体
type PacketGraspingSet struct {
	SetId             int     `json:"setId"`             //每条设置的标识符
	NativeServerPort  int     `json:"nativeServerPort"`  //本地端口
	HeartPort         int     `json:"heart_port"`        //心跳端口
	NativeIp          string  `json:"nativeIp"`          //本地ip
	RemoteIp          string  `json:"remoteIp"`          //远端ip，正则表达式形式
	PacketHoldingTime string  `json:"packetHoldingTime"` //抓包时间，单位min。默认空，永远抓包
	FileMaxSize       float64 `json:"fileMaxSize"`       //每个文件限制大小,默认20MB，单位MB
	SavePath          string  `json:"savePath"`          //文件存储路径,默认存在packet文件夹上
	FileName          string  `json:"fileName"`          //文件存储名字
}

//抓包状态返回结构体组
type Results struct {
	Message []Result
}

//抓包状态返回结构体
type Result struct {
	PGState
	SetId int
}

//配置文件结构体
type Config struct {
	Sets  []PacketGraspingSet `json:"sets"`
	MaxId int                 `json:"max_id"`
}

//抓包状态结构体
type PGState struct {
	PacketCount  int     //报文数量
	FileSize     float64 //文件大小
	FileCount    int     //文件数量
	PGHoldedTime int     //已抓包时间
	Finished     bool    //抓包是否结束
	State        string  //抓包状态
}

//声明抓包对象组
var Tests []PacketGrasping = make([]PacketGrasping, 0, 64)

//抓包对象结构体
type PacketGrasping struct {
	PacketGraspingSet                //抓包设置
	PGState                          //当前抓包的状态
	Flag              int64          //当前存的文件大小
	F                 *os.File       //用于文件操作（如创建打开关闭文件）
	W                 *pcapgo.Writer //用于写文件
	conn              *net.UDPConn
	heartConn         *net.UDPConn
	Display           bool //是否在网页显示
}
