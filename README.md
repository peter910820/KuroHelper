# KuroHelper

[![made-with-python](https://img.shields.io/badge/Made%20with-Golnag-00A7D0.svg)](https://www.python.org/)
[![Static Badge](https://img.shields.io/badge/Golang-1.24%2B-00A7D0)](https://www.python.org/downloads/release/python-3100/)
[![Discord](https://badgen.net/badge/icon/discord?icon=discord&label)](https://https://discord.com/)
[![GitHub license](https://img.shields.io/github/license/Naereen/StrapDown.js.svg)](https://github.com/peter910820/KuroHelper/blob/main/LICENSE)

A galgame helper use vndb and ErogameScape

# ⚙️Build:

1. clone this repo
```bash
$ git clone https://github.com/peter910820/KuroHelper.git
```
2. install package
```bash
$ go mod download
```
3. into repo and make .env, fill in required parameters(you must has a discord bot)
```bash
$ cd ./KuroHelper
$ cp .env.example .env
$ vim .env
```
4. run bot
```bash
$ go run ./cmd/core/main.go
```

# 💻Commands:

- global instructions:
    - vndb統計資料: 取得VNDB統計資料(VNDB)
    - 查詢指定遊戲: 根據VNDB ID查詢指定遊戲資料(VNDB)
    - 查詢遊戲: 根據關鍵字查詢遊戲資料(ErogameScape)
    - 查詢公司品牌: 根據關鍵字查詢公司品牌資料(VNDB&ErogameScape)
    - 查詢創作者: 根據關鍵字查詢創作者資料(ErogameScape)
    - 查詢音樂: 根據關鍵字查詢音樂資料(ErogameScape)
    - 查詢角色: 根據關鍵字查詢角色資料(ErogameScape)
    - 隨機遊戲: 隨機一部Galgame(ymgal)
    - 加已玩: 將遊戲加入已經遊玩的紀錄
    - 個人資料: 查詢使用者個人資料
- guild instructions:
    - 清除快取: 清除搜尋資料快取(管理員專用)
