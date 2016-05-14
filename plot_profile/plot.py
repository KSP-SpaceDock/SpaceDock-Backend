#!/usr/bin/python3
import sys
import requests
import os
import subprocess

if len(sys.argv) < 2:
    print("Usage: " + sys.argv[0] + " http://[spacedock-backend-url]")
    sys.exit()

r = requests.get(sys.argv[1] + "/profiler/histogram")
data = r.json()

id = 0
idmap = {}
for api_path in data:
    idmap[api_path] = id
    id = id + 1

#Always work from the script directory
work_path = os.path.join(os.path.dirname(os.path.realpath(sys.argv[0])), "output")

if os.path.exists(work_path):
    for file_name in os.listdir(work_path):
        file_path = os.path.join(work_path, file_name)
        if os.path.isfile(file_path):
            os.unlink(file_path)
else:
    os.mkdir(work_path)

plotstring = 'set terminal pngcairo size 1920,1080; set output "' + os.path.join(work_path, 'output.png') + '"; set yrange [0:]; set xlabel "Time"; set ylabel "Frequency"; plot '
firstplot = True
for api_path in data:
    id = idmap[api_path]
    api_data = data[api_path]
    api_max = 0
    for api_key in api_data.keys():
        api_current = int(api_key)
        if api_current > api_max:
            api_max = api_current
    data_file = os.path.join((work_path), str(id))
    with open(data_file + ".txt", "w") as f:
        for current_time in range(api_max + 2):
            if str(current_time) in api_data:
                f.write(str(current_time) + " " + str(api_data[str(current_time)]) + "\n")
            else:
                f.write(str(current_time) + " 0\n")
    if firstplot:
        firstplot = False
    else:
        plotstring = plotstring + ", "
    plotstring = plotstring + '"' + data_file + '.txt" with linespoints title "' + api_path + '"'
    subprocess.call(['gnuplot', '-e', 'set terminal pngcairo size 1920,1080; unset key; set xlabel "Time"; set ylabel "Frequency"; set title "' + api_path + '"; set yrange [0:]; set output "' + data_file + '.png" ; plot "' + data_file + '.txt" with linespoints'])

subprocess.call(['gnuplot', '-e', plotstring])
    