# Buffered

This example shows how to execute a child process and capture it's output
instead of having the child process inherit the current process's STDIO streams.

You should notice a couple things:

- The output will be output all at once instead of in realtime
- The ip address `127.0.0.1` was changed to `1.2.3.4`

## Expected Output

_Depending on OS & network conditions of course_

```
PING 1.2.3.4 (1.2.3.4) 56(84) bytes of data.
64 bytes from 1.2.3.4: icmp_seq=1 ttl=58 time=12.6 ms
64 bytes from 1.2.3.4: icmp_seq=2 ttl=58 time=10.2 ms
64 bytes from 1.2.3.4: icmp_seq=3 ttl=58 time=13.5 ms
64 bytes from 1.2.3.4: icmp_seq=4 ttl=58 time=19.1 ms

--- 1.2.3.4 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3005ms
rtt min/avg/max/mdev = 10.156/13.850/19.109/3.274 ms
```
