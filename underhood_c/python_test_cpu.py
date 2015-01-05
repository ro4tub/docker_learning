#echo 50000 > /sys/fs/cgroup/cpu/demo/cpu.cfs_quota_us
#echo PID > /sys/fs/cgroup/cpu/demo/tasks
import os
pidnum = os.getpid()
print pidnum
var = 1
count = 1
while var == 1:
    count = count + 1

