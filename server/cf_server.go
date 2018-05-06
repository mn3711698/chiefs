package main

import (
	"log"
	"net"
	"fmt"
	"strings"
)

const(
	BUFFERSIZE	=	1500
	MTU			=	"1300"
)
type serverConfig struct {
	listenAddr string
	nodeMap map[string]string

}

type cfServer interface {
	listen() *net.UDPConn
	contact(conn *net.UDPConn,cliAddr *net.UDPAddr)
}


func (sconfig *serverConfig) init() *serverConfig {
	sconfig=&serverConfig{nodeMap:make(map[string]string),listenAddr:":9981"}
	return sconfig
}
func (sconfig *serverConfig) listen() (*net.UDPAddr) {
	listenAddr,err:=net.ResolveUDPAddr("udp",sconfig.listenAddr)
	checkFatalErr(err,"Unable to get UDP socket:")
	return listenAddr
}

func (sconfig *serverConfig) contact(listenAddr *net.UDPAddr){
	//listen port 9981
	conn,err:=net.ListenUDP("udp",listenAddr)
	checkFatalErr(err,"Unable to listen on UDP socket:")
	log.Println("server start at 0.0.0.0",sconfig.listen)
	defer conn.Close()
	buf:=make([]byte,BUFFERSIZE)
	for{
		//var clientAddr *net.UDPAddr
		n,addr,err:=conn.ReadFromUDP(buf)
		fmt.Println("Reciverd data: ",string(buf[:n]))
		if err!=nil || n==0 {
			checkFatalErr(err, "Reciverd data error!")
			continue
		}
		log.Println("client addr:",addr)
		cname:=strings.Split(string(buf[:n]),":")[0]
		sconfig.nodeMap[cname]=fmt.Sprintf("%v",addr)
		//遍历所有的node节点,发送消息
		for _,ar:= range sconfig.nodeMap{
			clientAddr,err:=net.ResolveUDPAddr("udp",ar)
			if err!=nil{
				continue
			}
			conn.WriteToUDP([]byte("hello "),clientAddr)
			//log.Println("clientAddr:",clientAddr)
		}
	}
	conn.Close()
}


func checkFatalErr(err error,msg string){
	if err!=nil{
		log.Println(msg)
		log.Fatal(err)
	}
}
func start()  {
	var sconfig=&serverConfig{}
	sconfig=sconfig.init()
	listenAddr:=sconfig.listen()
	sconfig.contact(listenAddr)
}
func main(){
	start()
}
