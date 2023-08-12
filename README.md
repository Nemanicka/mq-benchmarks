# Results

| MQ               |  RW   | Total Requests #  | Total Time Write   | RPS          | MBPS         |
|------------------|-------|-------------------|--------------------|--------------|--------------|
| Redis RDB        | W     | 100000            | 59.9495298s        | ~1694        | 30.492       |
| Redis RDB        | R     | 100000            | 7.4953866s         | ~13513       | 243.234      |
| Redis AOF        | W     | 100000            | 1m6.7837954s       | ~1499        | 26.982       |
| Redis AOF        | R     | 100000            | 8.1965398s         | ~12345       | 222.210      |
| Redis No Persist | W     | 100000            | 54.3772661s        | ~1841        | 33.138       |
| Redis No Persist | R     | 100000            | 3.8718633s         | ~26315       | 473.670      |
| BeanStalkd       | W     | 100000            | 18.8162104s        | ~5319        | 95.742       |
| BeanStalkd       | R     | 100000            | 8.5316565s         | ~11764       | 211.752      |

# Methodolodgy

Every MQ was give 100 concurrent connection, each consisted of 1000 requests.
The size of one request (defined in *constants.go*) is equal to ~18Kb.
Consequently, the MBPS is calculated as (rps)*(payload size), and since I have a constant number 
for each measurement, rps is calculated as (total requests) / (total time)

# Conclusion

**What did go well?**

Apperently, Beanstalkd consistently showed better write performace, and read perforance is roughly the same as in Redis.

**What didn't go well?**

Although I made sure multiple times that rdb, aof and no-persist redis have proper configuration, and in logs I can
see approprite messages - like writing data to disk, or absence of these messages (for no persist redis), anyway
write and read performance is equal on average, with some deviation due to some random backgroung processes (I guess).