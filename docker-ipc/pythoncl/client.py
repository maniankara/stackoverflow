#!/usr/bin/env python
#client example
import socket
import time

print "Waiting for socket"
time.sleep(3)
client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
client_socket.connect(('java-server', 10001))
client_socket.send("hello")
client_socket.send(" stackoverflow")
client_socket.close()
while 1:
  print "pending"
  time.sleep(3)
  pass # do nothing
