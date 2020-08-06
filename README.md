# Logger

This is a simple library to improve logging consistancy

## Config

Expect environment variables : 
```bash
<PRE>LOG_SYSLOG="true"
<PRE>LOG_LEVEL="trace"
<PRE>LOG_FILE="/dev/stdout"
<PRE>LOG_PREFIX=""
```

## Format

Standard output is used for specific system message and Warning and Errors. These logs will contains only the logged message without specific pattern<br/>

The file log will use the followin pattern : 

```bash
YYYY/MM/DD HH:mm:ss.uuuuuu: [LEVEl] "message"
```

for example :

```bash
2019/11/02 19:05:37.377333: [INFO] "Upload complete : /tmp/README.md"
```
