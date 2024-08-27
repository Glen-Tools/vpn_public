#!/bin/bash

# 產出 exe 檔案
GOOS=windows GOARCH=amd64 go build -o v2ray_client.exe

if [ $? -eq 0 ]; then
    echo "成功生成 v2ray_client.exe"
else
    echo "生成 v2ray_client.exe 失敗"
fi

# 定義新資料夾的名稱
new_folder="v2ray_client"

# 新建資料夾
mkdir -p "$new_folder"


# 將資料夾移動到新建的資料夾中
cp -r "./config" "./$new_folder/"
cp -r "./readme" "./$new_folder/"
cp -r "./v2ray_client.exe" "./$new_folder/v2ray_client.exe"
cp -r "./v2ray_file/windows" "./$new_folder/v2ray_file/windows"

echo "資料夾已成功移動到 $new_folder/"