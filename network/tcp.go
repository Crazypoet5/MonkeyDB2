package network

import "net"

func uint322bytes(u uint32) []byte {
	b := make([]byte, 4)
	for i := 0; i < 4; i++ {
		b[i] = (byte)(u)
		u >>= 8
	}
	return b
}

func bytes2uint32(b []byte) uint32 {
	p := uint32(0)
	for i := 3; i >= 0; i-- {
		p <<= 8
		p |= uint32(b[i])
	}
	return p
}

type TCPSession struct {
	Conn   *net.TCPConn
	Closed bool //是否已经关闭
}

func (s *TCPSession) SendPack(p *Pack) {
	var buff []byte
	head := uint322bytes(p.Head)
	leng := uint322bytes(p.Len)
	typp := uint322bytes(p.Type)
	buff = append(head, leng...)
	buff = append(buff, typp...)
	_, err := s.Conn.Write(buff)
	if err != nil {
		s.Closed = true
		return
	}
	offset := 0
	for len(p.Data) > 1024 {
		_, err = s.Conn.Write(p.Data[offset : offset+1024])
		if err != nil {
			s.Closed = true
			return
		}
		offset += 1024
	}
	_, err = s.Conn.Write(p.Data[offset:])
	if err != nil {
		s.Closed = true
		return
	}
}

func (s *TCPSession) RecvPack() *Pack {
	var buff = make([]byte, 12)
	readed, err := s.Conn.Read(buff)
	if err != nil {
		s.Closed = true
		return nil
	}
	if readed != 12 {
		return nil
	}
	leng := bytes2uint32(buff[4:8])
	ret := &Pack{
		Head: bytes2uint32(buff[0:4]),
		Len:  leng,
		Type: bytes2uint32(buff[8:12]),
	}
	leng -= 12
	buff = make([]byte, 1024)
	for leng > 1024 {
		_, err = s.Conn.Read(buff)
		if err != nil {
			s.Closed = true
			return nil
		}
		ret.Data = append(ret.Data, buff...)
	}
	readed, err = s.Conn.Read(buff)
	if err != nil {
		s.Closed = true
		return nil
	}
	ret.Data = append(ret.Data, buff[0:readed]...)
	return ret
}

func Listen() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", ":2016")
	if err != nil {
		panic(err)
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server started.")
	for {
		conn, err := l.AcceptTCP()
		//TODO
	}
}
