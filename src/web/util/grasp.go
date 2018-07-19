package util

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	//"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/pcapgo"
)

//初始化抓包设置结构体
func InitPacketGraspingSet(r *http.Request) PacketGraspingSet {
	var set PacketGraspingSet
	port := r.Form.Get("port")
	heart_port := r.Form.Get("heart_port")
	during_time := r.Form.Get("during_time")
	file_max_size := r.Form.Get("file_max_size")
	file_name := r.Form.Get("file_name")
	folder := r.Form.Get("folder")
	if file_max_size == "" {
		file_max_size = "20"
	}
	if folder == "" {
		folder = "packet"
	}
	nativeIp := r.Form.Get("nativeIp")
	remoteIp := r.Form.Get("remoteIp")
	set.NativeIp = nativeIp
	set.RemoteIp = remoteIp
	heartport, _ := strconv.Atoi(heart_port)
	set.NativeServerPort, _ = strconv.Atoi(port)
	set.HeartPort = heartport
	set.PacketHoldingTime = during_time
	set.FileMaxSize, _ = strconv.ParseFloat(file_max_size, 64)
	set.SavePath = folder
	set.FileName = file_name
	return set
}

//初始化抓包对象
func (Test *PacketGrasping) Set(id int) (bool, error) {
	existence, set := FoundConfigById(id) //根据id从配置文件获取配置信息
	if existence == true {
		Test.PacketGraspingSet = set
		Test.PGState.PacketCount = 0
		Test.PGState.FileSize = float64(0)
		Test.PGState.FileCount = 0
		Test.PGState.Finished = false
		Test.PGState.PGHoldedTime = 0
		Test.PGState.State = "正在启动抓包任务"
		Test.Display = true
		return true, nil
	}
	err := errors.New("未找到id对应的配置")
	return false, err
}

//根据配置创建文件
func (Test *PacketGrasping) PGCreateFile() {
	//fmt.Printf("id为%d的线程创建文件\n", Test.SetId)
	InitFolder(Test.SavePath)
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	fil := Test.SavePath + "/" + Test.FileName + "_" + timestamp + ".pcap"
	Test.F, _ = os.Create(fil)
	Test.PGState.FileCount++
	Test.W = pcapgo.NewWriter(Test.F)
	Test.W.WriteFileHeader(1024, layers.LinkTypeEthernet)
}

//开始抓包
func (Test *PacketGrasping) PGStart() {
	var addr *net.UDPAddr
	var err error
	if Test.NativeIp == "" {
		addr, err = net.ResolveUDPAddr("udp", "0.0.0.0:"+strconv.Itoa(Test.NativeServerPort))
	} else {
		addr, err = net.ResolveUDPAddr("udp", "["+Test.NativeIp+"]:"+strconv.Itoa(Test.NativeServerPort))
	}
	if err != nil {
		Test.PGState.State = "本地ip+端口不合法"
		log.Println(Test.PGState.State)
		Test.Finished = true
	} else {
		Test.conn, err = net.ListenUDP("udp", addr)
		if err != nil {
			Test.PGState.State = "无法建立抓包连接"
			log.Println(Test.PGState.State)
			Test.Finished = true
		} else {
			Test.PGState.State = "正在抓包"
			go Test.PGSaveFile()
			go Test.ListenHeart()
			go Test.Holdtime()
			for {
				time.Sleep(100)
				if Test.Finished == true {
					Test.Close()
					break
				}
			}
		}
	}
}
func (Test *PacketGrasping) ListenHeart() {
	buffer := make([]byte, 65535)
	dst, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(Test.HeartPort))
	if err != nil {
		log.Println("心跳端口不合法")
		Test.PGState.State = "心跳端口不合法"
		Test.Finished = true
		return
	}
	Test.heartConn, err = net.ListenUDP("udp", dst)
	if err != nil {
		log.Println("无法建立心跳连接")
		Test.PGState.State = "无法建立心跳连接"
		Test.Finished = true
		return
	}
	timeout := 0
	for {
		Test.heartConn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if !Test.Finished {
			num, src, err := Test.heartConn.ReadFromUDP(buffer)
			if err != nil {
				timeout++
				if timeout > 5 {
					log.Println("心跳超时，一分钟以内没收到心跳报文")
					Test.PGState.State = "心跳超时"
					Test.Finished = true
					return
				}
				continue
			}
			timeout = 0
			Test.heartConn.WriteToUDP(buffer[:num], src)
		} else {
			break
		}

	}
}
func (Test *PacketGrasping) Holdtime() {
	flag := false
	if Test.PacketHoldingTime == "" {
		for {
			if !Test.Finished {
				Test.PGHoldedTime++
			} else if Test.Finished {
				flag = true
				break
			}
			time.Sleep(1 * time.Second)
		}
	} else {
		packetHoldingTime, _ := strconv.Atoi(Test.PacketHoldingTime)
		for i := 0; i < 60*int(packetHoldingTime); i++ {
			if !Test.Finished {
				Test.PGHoldedTime++
			} else if Test.Finished {
				flag = true
				break
			}
			time.Sleep(1 * time.Second)
		}
	}
	if !flag {
		Test.PGEnd()
	}

}

//保存文件
func (Test *PacketGrasping) PGSaveFile() {
	buffer := make([]byte, 65535)
	Test.PGCreateFile()
	for {
		if !Test.Finished {
			num, addr, err := Test.conn.ReadFromUDP(buffer)
			if Test.RemoteIp != "" {
				if addr != nil {
					if m, _ := regexp.MatchString(Test.RemoteIp, addr.IP.String()); !m {
						continue
					}
				}
			}
			if nil != err {
				continue
			} else {
				captureInfo := gopacket.CaptureInfo{
					Timestamp:      time.Now().UTC(),
					CaptureLength:  num,
					Length:         num,
					InterfaceIndex: 0,
				}
				Test.W.WritePacket(captureInfo, buffer[:num])
				Test.Flag, _ = Test.F.Seek(0, os.SEEK_END)
				Test.PGState.FileSize = float64(Test.FileCount-1)*Test.FileMaxSize + float64(Test.Flag)/(float64(1024)*float64(1024))
				// Only file size more 20MB and then stop
				size := Test.FileMaxSize * 1024 * 1024
				Test.PGState.PacketCount++
				if Test.Flag > int64(size) {
					Test.F.Close()
					Test.PGCreateFile()
				}
			}
		} else {
			break
		}

	}
}

//结束抓包
func (Test *PacketGrasping) PGEnd() {
	Test.PGState.Finished = true
	Test.PGState.State = "抓包结束"
	fmt.Printf("id为%d的抓包任务已结束，共抓包%d个，抓包文件大小为%gMB,共存文件%d个,抓了%ds包\n", Test.SetId, Test.PacketCount, Test.FileSize, Test.FileCount, Test.PGHoldedTime)
}
func (Test *PacketGrasping) Close() {
	if Test.F != nil {
		Test.F.Close()
	}
	if Test.conn != nil {
		Test.conn.Close()
	}
	if Test.heartConn != nil {
		Test.heartConn.Close()
	}

}
func CheckIp(ip string, port string) (bool, error) {
	var err error
	if ip == "" {
		_, err = net.ResolveUDPAddr("udp", "0.0.0.0:"+port)
	} else {
		_, err = net.ResolveUDPAddr("udp", "["+ip+"]:"+port)
	}
	if err != nil {
		err = errors.New("本地ip加端口不是有效的地址")
		return false, err
	} else {
		return true, nil
	}
}
func Refresh() bool {
	if IsTasksFinished() {
		for i := 0; i < len(Tests); i++ {
			Tests[i].Close()
		}
		Tests = make([]PacketGrasping, 0, 64)
	}
	return true
}
