#!/usr/bin/python3
import sys
import requests
import os
import subprocess

if len(sys.argv) < 2:
    print("Usage: " + sys.argv[0] + " http://[spacedock-backend-url]")
    sys.exit()

serverpart = sys.argv[1][sys.argv[1].index('//') + 2:]
port = 80
if ":" in serverpart:
    splitindex = serverpart.rfind(":")
    port = serverpart[splitindex + 1:]
    serverpart = serverpart[:splitindex]

r = requests.get(sys.argv[1] + "/documentation")
data = r.json()

for key in data:
    if not "<" in key:
        subprocess.call(['httperf', '--server=' + serverpart, '--port=' + port, '--uri=' + key, '--num-conns=10', '--rate=1'])
