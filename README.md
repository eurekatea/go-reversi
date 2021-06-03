# othello-board
othello(reversi) arena for humans and AIs

# 遊戲截圖
![](https://user-images.githubusercontent.com/65079334/120668578-2d240680-c4c1-11eb-957f-25caebe218ec.png)
![](https://user-images.githubusercontent.com/65079334/120668607-3319e780-c4c1-11eb-8096-c23c2e3a05dc.png)

# 規則
規則與一般黑白棋相同，版面大小為6x6，若有一方無處可下會自動PASS，換另一方下  
雙方皆無處可下時遊戲結束，依棋子數目決定輸贏或平手  

# 說明
外部的AI程式須接收input，並輸出結果  
如：輸入```++++++++++++++OX++++XO++++++++++++++ 1```，輸出```Bc```  
(X表示黑方，O表示白方；1表示為黑方，2為白方)  
程式接收回傳值時，若回傳值不合法，GUI會顯示外部AI出錯，並將錯誤內容輸出至log file內  

# 下載
https://github.com/lemon37564/othello-board/releases

# 自行編譯
require go 1.16+  
先```go get fyne.io/fyne/v2```
### windows
```go build -ldflags="-H windowsgui"```  
### linux or macOS
```go build```  
