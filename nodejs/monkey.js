var net = require('net')

var HOST = '127.0.0.1'
var PORT = '2016'

const CREATE_CONNECTION = 0
const CLOSE_CONNECTION = 1
const DIRECT_QUERY = 2
const RESPONSE = 3

function MackPack(type, data) {
    return {
        Head: 2016,
        Len:  data.length + 12,
        Type: type,
        Data: data,
    }
}

function uint322bytes(u) {
    var b = new Buffer([0,0,0,0])
    for (var i = 0;i<4;i++) {
        b[i] = u & 0x000000ff
        u >>= 8
    }
    return b
}

function bytes2uint32(b) {
    var p = 0
    for (var i = 3;i >= 0;i--) {
        p <<= 8
        p |= b[i]
    }
    return p
}

var monkey = {
    seq: 0,
    temp: "",
    tempcallback: null,
    templen: 0,
    client: new net.Socket(),
    connect: function(callback) {
        this.client.connect({port:'2016', host:'127.0.0.1'}, function() {
            this.on('data', monkey.ondata)
            if (callback != null) {
                callback()
            }
        })
    },
    query: function(queryStr, callback) {
        var buff = Buffer.concat([uint322bytes(2016),
            uint322bytes(queryStr.length + 12),
            uint322bytes(DIRECT_QUERY)])
        monkey.client.write(buff, function(){
            monkey.client.write(new Buffer(queryStr), function() {
                monkey.seq = 0
                monkey.temp = ""
                monkey.templen = 0
                monkey.tempcallback = callback
            })
        })
    },
    ondata: function(data) {
        if (monkey.seq == 0) {
            var head = bytes2uint32(new Buffer([data[0],data[1],data[2],data[3]]))
            var len = bytes2uint32(new Buffer([data[4],data[5],data[6],data[7]]))
            var type = bytes2uint32(new Buffer([data[8],data[9],data[10],data[11]]))
            data = data.toString('ASCII',12)
            if (head != 2016) {
                monkey.seq = 0
                monkey.temp = ""
                monkey.templen = 0
                return
            }
            var leng = len - 12
            if (leng > 0) {
                monkey.templen = len - data.length - 12
                monkey.temp += data.slice(12)
                if (monkey.templen == 0) {
                    monkey.temp += data.slice(0, monkey.templen)
                    monkey.seq = 0
                    if (monkey.tempcallback != null)
                        monkey.tempcallback(monkey.temp)
                    monkey.temp = ""
                    monkey.templen = 0
                    return
                }
            } else {
                monkey.seq = 0
                if (monkey.tempcallback != null)
                    monkey.tempcallback(monkey.temp)
                monkey.temp = ""
                monkey.templen = 0
                return
            }
            monkey.seq++
        } else {
            data = data.toString('ASCII')
            if (monkey.templen > 0) {
                if (data.length > monkey.templen) {
                    monkey.temp += data.slice(0, monkey.templen)
                    monkey.seq = 0
                    if (monkey.tempcallback != null)
                        monkey.tempcallback(monkey.temp)
                    monkey.temp = ""
                    monkey.templen = 0
                    return
                    //TODO: data remain
                }
                monkey.templen -= data.length
                monkey.temp += data
                if (monkey.templen == 0) {
                    monkey.temp += data.slice(0, monkey.templen)
                    monkey.seq = 0
                    if (monkey.tempcallback != null)
                        monkey.tempcallback(monkey.temp)
                    monkey.temp = ""
                    monkey.templen = 0
                    return
                }
                monkey.seq++
            }
        }
    },
    close: function() {
        monkey.client.destroy()
    }
}

/*
monkey.connect(function() {
    monkey.query("select * from t1 where id=2", function(data) {
        console.log(data)
        monkey.query("select * from t1 where id=2", function(data) {
            console.log(data)
            monkey.query("select * from t1 where id=2", function(data) {
                console.log(data)
            })
        })
    })
   
})*/

module.exports = monkey