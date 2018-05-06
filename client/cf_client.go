package main
import (
	"net"
	"log"
	//"time"
	"time"
	"strconv"
)
const(
	BUFFERSIZE	=	1500
	FLAG		=	9981
	HOSTNAME	=	"node1"
)
type clientConfig struct {
	clientAddr string
	listenAddr string
}

type cfClient interface {
	getConn() *net.UDPConn
	contact(conn *net.UDPConn,cliAddr *net.UDPAddr)
}

func checkFatalErr(err error,msg string){
	if err!=nil{
		log.Println(msg)
		log.Fatal(err)
	}
}
func (cconfig *clientConfig) init() *clientConfig {
	 cconfig=&clientConfig{clientAddr:"127.0.0.1:9981",listenAddr:":1997"}
	 return cconfig
}
func (cconfig *clientConfig) getConn() (*net.UDPConn,*net.UDPAddr) {
	//解析地址字符串
	cliAddr,err:=net.ResolveUDPAddr("udp",cconfig.clientAddr)
	checkFatalErr(err,"Unable to resolve server UDP socket")
	lisAddr,err:=net.ResolveUDPAddr("udp",cconfig.listenAddr)
	checkFatalErr(err,"Unable to resolve server UDP socket")
	conn,err:=net.ListenUDP("udp",lisAddr)
	checkFatalErr(err,"Unable connect to server")
	return conn,cliAddr
}

func (cconfig *clientConfig) contact(conn *net.UDPConn,cliAddr *net.UDPAddr){
	buf:=make([]byte,BUFFERSIZE)
	conn.WriteToUDP([]byte(HOSTNAME+": world "+strconv.Itoa(FLAG)),cliAddr)
	n,addr,err:=conn.ReadFromUDP(buf)
	if err!=nil{
		log.Println("Unable recived data!")
	}
	log.Println("Recived server data ",addr,string(buf[:n]))
}

func start()  {
	var cconfig *clientConfig
	cconfig=cconfig.init()
	conn,serAddr:=cconfig.getConn()
	for{
		time.Sleep(3*time.Second)
		go cconfig.contact(conn,serAddr)
	}
	conn.Close()
}

func main() {
	start()
}