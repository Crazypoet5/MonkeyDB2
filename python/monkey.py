#!/usr/bin/python

import socket

CREATE_CONNECTION = 0
CLOSE_CONNECTION = 1
DIRECT_QUERY = 2
RESPONSE = 3

s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

def uint322bytes(a):
    i = 0
    s = ''
    while (i < 4):
        s = s + chr(a % 256)
        a /= 256
        i = i + 1
    return s

def bytes2uint32(a):
    p = 0
    i = 3
    while (i >= 0):
        p = p * 256
        p = p | ord(a[i])
        i = i - 1
    return p

def MakePack(tp, data):
    pack = {'Head': 2016, 'Len': len(data), 'Type': tp, 'Data': data}
    return pack
    
def Connect():
    s.connect(('127.0.0.1', 2016))

def SendPack(p):
    head = uint322bytes(p['Head'])
    leng = uint322bytes(p['Len'])
    typp = uint322bytes(p['Type'])
    buff = head + leng + typp
    s.send(buff)
    if len(p['Data']) == 0:
        return
    else:
        s.send(p['Data'])   #If you send much data, please use buffer blocks # Maybe... Python implement it internal
                            #oh... I don't know Python

def RecvPack():
    buff = s.recv(12)
    leng = bytes2uint32(buff[4:8])
    ret = {'Head': bytes2uint32(buff[0:4]), 'Len': leng, 'Type': bytes2uint32(buff[8:12])}
    leng = leng - 12
    data = s.recv(leng)
    ret['Data'] = data
    return ret
    
def Close():
    s.close()

def SendCmd(cmd):
    SendPack(MakePack(DIRECT_QUERY, cmd))
    p = RecvPack()
    return p['Data']


# Usage

# import monkey

# monkey.Connect()
# res = monkey.SendCmd('createkv kv1 string string')
# res = monkey.SendCmd('set kv1 \'a\' \'123456\'')
# res = monkey.SendCmd('get kv1 \'a\'')
# print res
# monkey.Close()