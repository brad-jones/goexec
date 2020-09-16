# Simple

So you just want to execute a child process with a minimum of fuss,
just use `goexec.Run()` as in this example.

## Expected Output

_Depending on OS & network conditions of course_

```
PING 127.0.0.1 (127.0.0.1) 56(84) bytes of data.
64 bytes from 127.0.0.1: icmp_seq=1 ttl=58 time=12.6 ms
64 bytes from 127.0.0.1: icmp_seq=2 ttl=58 time=10.2 ms
64 bytes from 127.0.0.1: icmp_seq=3 ttl=58 time=13.5 ms
64 bytes from 127.0.0.1: icmp_seq=4 ttl=58 time=19.1 ms

--- 127.0.0.1 ping statistics ---
4 packets transmitted, 4 received, 0% packet loss, time 3005ms
rtt min/avg/max/mdev = 10.156/13.850/19.109/3.274 ms
```
