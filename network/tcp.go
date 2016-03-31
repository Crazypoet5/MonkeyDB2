package tcp
import "net"
import "fmt"

func bytes4uint(bytes []byte) uint32 {
	total := uint32(0);	
	for i := 0;i < 4;i++ {
		total <<= 8;
		total += uint32(bytes[i]);
	}
	return total
}

func uint32bytes(n uint32) []byte {
	header := make([]byte,4)
	i := 0
	for n > 0 {
		header[3-i] = byte(n % 256)
		n /= 256
		i++
	}
	return header
}


type TCPSession struct {
	Conn *net.TCPConn
	Closed bool //是否已经关闭
}

func (s *TCPSession) SendMessage(bytes []byte) {
	total := len(bytes)
	header := uint32bytes(uint32(total))	//计算字节数
	_,err := s.Conn.Write(header)
	if err != nil {
		s.Closed = true
		return
	}
	//fmt.Println("send:",header)
	//fmt.Println(header)
	for i := 0;i < total - 1024;i += 1024 {
		buf := bytes[0:1024]	//发送这一段
		bytes = bytes[1024:]
		_,err := s.Conn.Write(buf)
		if err != nil {
			s.Closed = true
			return
		}
		//fmt.Println("send:",buf)
		continue
	}
	buf := bytes[0:]	//发送这一段
	_,err = s.Conn.Write(buf)
	if err != nil {
		s.Closed = true
		return
	}
	//fmt.Println("send:",buf)
}

func (s *TCPSession) ReadMessage() []byte {
	buf := make([]byte,4)
	_,err := s.Conn.Read(buf)
	if err != nil {
		s.Closed = true
		return []byte{}
	}
	//fmt.Println(buf)
	total := int(bytes4uint(buf))
	var buff []byte
	buf = make([]byte,1024)
	for total > 1024 {
		_,err := s.Conn.Read(buf)
		if err != nil {
			s.Closed = true
			return []byte{}
		}
		buff = append(buff,buf...)
		total -= 1024
	}
	buf = make([]byte,total)
	_,err = s.Conn.Read(buf)
	if err != nil {
		s.Closed = true
		return []byte{}
	}
	buff = append(buff,buf...)
	return buff
}