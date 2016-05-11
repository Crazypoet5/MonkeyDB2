import monkey

monkey.Connect()
cmd = raw_input("Monkey>>")
while (cmd != "quit" and cmd != "quit;"):
    res = monkey.SendCmd(cmd)
    print res
    cmd = raw_input("Monkey>>")
monkey.Close()