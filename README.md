### malicious-validator

malicious-validator 실행 가이드

```
# Container 실행
docker-compose up -d

# Conatiner 동작 확인
docker ps

# Container Log 확인
docker logs $(docker container ls --all | grep malicious-validator | awk '{print $1}')

# Log 실시간 조회
docker logs $(docker container ls --all | grep malicious-validator | awk '{print $1}') -f
```


```
# 종료
docker-compose down
```


Deposit
```
# honest
docker exec -it  $(docker container ls --all | grep kroma-validator | awk '{print $1}') /bin/sh
# shell에서 실행
> kroma-validator deposit --amount 1000000000000000000

# malicious
docker exec -it  $(docker container ls --all | grep kroma-validator | awk '{print $1}') /bin/sh
# shell에서 실행
> malicious-validator deposit --amount 1000000000000000000
```