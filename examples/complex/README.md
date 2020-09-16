# Complex

This example shows you the decorator based API that is useful for constructing
long and complicated command line calls, for example like what you might do with
`docker run ...`

All of the other functionality, like buffering & prefixing is still available by
passing a built `*exec.Cmd` object to the appropriate function like
`goexec.RunBufferedCmd()`.

## Expected Output

```
total 12K
drwxrwxrwx    1 root     root        4.0K Sep 16 06:45 .
drwxr-xr-x    3 root     root        4.0K Sep 16 07:21 ..
-rwxrwxrwx    1 root     root         854 Sep 16 06:55 README.md
-rwxrwxrwx    1 root     root         650 Sep 16 07:20 main.go
-rwxrwxrwx    1 root     root         629 Sep 16 06:45 main_test.go
HOSTNAME=36cd5aa0f460
SHLVL=1
HOME=/root
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin
FOO=BAR
PWD=/c/Users/brad.jones/Projects/Personal/goexec/examples/complex
```
