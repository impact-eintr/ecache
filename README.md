# ecached
分布式缓存

## 安装
ecached的落盘是基于rocksdb

> arch系
```bash
sudo pacman -S rocksdb

git clone https://github.com/impact-eintr/ecache.git

cd ecache/cmd/ecachd

go build

./ecachd -h
Usage of ecached:
  -T int
        缓存生存时间 默认为0 即不失效
  -d string
        磁盘缓存目录 请务必指定 (default "/tmp/ecached")
  -hp string
        ecached http 服务端口 (default "7895")
  -t string
        缓存类型 可选项 mem disk (default "mem")
  -tp string
        ecached tcp 服务端口 (default "6430")
```

