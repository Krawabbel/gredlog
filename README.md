a REDIS-based logger written in Go (Golang)

## Useful Commands

* read the current temperature of our Raspberry Pi in the terminal via
```vcgencmd measure_temp```

* REDIS commands (alternatives)
    - start REDIS server ```sudo systemctl start redis```
    - stop REDIS server ```sudo systemctl stop redis```

    - start REDIS docker ```sudo docker run -p 6379:6379 -it --rm redis/redis-stack-server```

* start tcp socket and listen for commands ```echo -n -e "+OK\r\n" | nc -l 6379``` (useful for RESP debugging)


## Resources

* Python implementation of temperature logger in GPIO Zero library: https://gpiozero.readthedocs.io/en/stable/_modules/gpiozero/internal_devices.html#CPUTemperature
* REDIS TimeSeries Quickstart with Docker: https://redis.io/docs/data-types/timeseries/quickstart/
* REDIS TimeSeries Tutorial on Medium: https://medium.com/datadenys/using-redis-timeseries-to-store-and-analyze-timeseries-data-c22c9e74ff46 
* REDIS streams: https://redis.io/docs/data-types/streams/
* Build REDIS from scratch: https://www.build-redis-from-scratch.dev/en/introduction
* Build a concurrent TCP server in Go: https://opensource.com/article/18/5/building-concurrent-tcp-server-go