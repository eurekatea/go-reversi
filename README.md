# othello-board
othello(reversi) arena for humans and AIs

# 遊戲截圖
![](https://user-images.githubusercontent.com/65079334/120412652-7ec97580-c389-11eb-99f6-c338a38208c0.png)
![](https://user-images.githubusercontent.com/65079334/120412656-80933900-c389-11eb-9735-70881212c586.png)

# 規則
規則與一般黑白棋相同，版面大小為6x6，若有一方無處可下會自動PASS，換另一方下  
雙方皆無處可下時遊戲結束，依棋子數目決定輸贏或平手  

# 說明
此程式會將版面以"++++++++++++++OX++++XO++++++++++++++ 1"形式傳送給AI之stdin  
(X表示黑方，O表示白方；1表示現在為黑方行動，2反之)  
程式在AI之stdout接收回傳值，若回傳值不合法，則在log file內印出錯誤訊息並退出  

# 下載
https://github.com/lemon37564/othello-board/releases

# 自行編譯
require go 1.16+  
先```go get fyne.io/fyne/v2```
### windows
```go build -ldflags="-H windowsgui"```  
### linux or macOS
```go build```  
