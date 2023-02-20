goctl docker -go main.go

docker build -t tiktok:v1

docker run -it -p :8090:8090 tiktok:v1