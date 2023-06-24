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