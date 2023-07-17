#!/bin/bash

version=$(git describe --long --dirty --abbrev=6 --tags)
flags="-X main.buildTime=$(date -u '+%Y-%m-%d_%I:%M:%S%p') -X main.version=$version"

export CGO_ENABLED=0 GOOS=linux

echo "start to build version:$version"

echo "flags:$flags"

current_path=$(pwd)

echo "当前路径：$current_path"

cd front

# npm install yarn

# yarn install

# yarn build

cd ..

rm -rf ./output
mkdir -p ./output/linux/


for i in "arm64" "amd64" ; do
  echo "building for $i..."
  GOARCH="$i" go build -ldflags "$flags" -o ./output/linux/$i/app main.go
  cp -r ./front ./output/linux/$i
  mkdir ./output/linux/$i/config
  cp ./demo.ini ./output/linux/$i/config/conf.ini
  chmod +x ./output/linux/$i/app
done

build_darwin() {
  echo "building for darwin"
  GOOS=darwin GOARCH=amd64 go build -ldflags "$flags" -o ./output/linux/darwin/app main.go
  cp -r ./front ./output/linux/darwin
  mkdir ./output/linux/darwin/config
  cp ./demo.ini ./output/linux/darwin/config/conf.ini
  chmod +x ./output/linux/darwin/app
}

build_darwin
# image=ccr.ccs.tencentyun.com/imoe-tech/go-playground:ikuai-exporter-"$version"
# official_img=jakes/ikuai-exporter:latest
# official_img_versioned=jakes/ikuai-exporter:"$version"
# echo "packaging docker multiplatform image: $image"
# echo "packaging docker multiplatform image: $official_img"
# echo "packaging docker multiplatform image: $official_img_versioned"

# docker buildx build --push \
#   --platform linux/amd64,linux/arm64 \
#   --build-arg VERSION="$version" \
#   -t "$image" -t "$official_img" -t "$official_img_versioned" .

# echo "finished: $image"


