# StickerDiscordBot
一個簡單的複合性功能機器人

## 功能
* 貼圖回覆 - 幫你回復設定好的圖片
* 支語糾察 - 偵測到支語時，支與警察會跳出來

## 安裝
**[bot 篇]**

下載程式碼
```
git clone https://github.com/orangesobeautiful/StickerDiscordBot
```
建立python虛擬環境(可選)
```
sudo apt install python3-virtualenv
virtualenv pyvenv
source pyvenv/bin/activate
pip3 install --upgrade pip setuptools
```
安裝 python 依賴套件
```
pip3 install -r requirements.txt
```
生成設定檔(程式開始執行時如果沒有偵測到設定檔會自動產生)
```
python3 runBot.py
```
編寫設定檔(setting.ini)，修改 機器人token 和 資料庫

`bot_token = 你的機器人token`  
`database_url = mysql+pymysql://<user_name>:<password>@<ip_address>/<database_name>`

啟動機器人
```
python3 runBot.py
```
  
**[web 篇]**

網頁用來代替機器人指令管理圖片庫，如果沒有需求可以選擇略過

需要將 /sndata /sticker-image 導向 web 後端 

* 前端

安裝前端依賴套件
```
cd web/frontend
sudo npm install -g @quasar/cli
npm install
```
建構 web
```
quasar build
```
* 後端

編寫設定檔(setting.ini)，修改 機器人token 和 資料庫

`accesswebverificationguild = <Discord 群組 ID>`  


```
python3 runWebBackend.py
```
## 指令

| 命令 | 描述
|---------|-------------|
| `$add <圖片名稱> <圖片網址>` | 新增一個自動回覆的圖片 |
| `$show <圖片名稱>` | 顯示設定過的圖片 |
| `$edit <圖片ID> <圖片網址>` | 修改圖片網址 |
| `delete <圖片ID>` | 刪除指定ID圖片 |
| `deleteST <圖片名稱>` | 刪除圖片名稱設定過的所有圖片 |
