# othello-board
othello(reversi) arena for humans and AIs

# 遊戲截圖
![](https://user-images.githubusercontent.com/65079334/120734562-a9e3ce80-c51b-11eb-8590-6033cb762637.png)
![](https://user-images.githubusercontent.com/65079334/120924683-68177b80-c707-11eb-92ce-d87c60a4db26.png)

# 規則
規則與一般黑白棋相同，版面大小為6x6(或8x8)，若有一方無處可下會自動PASS，換另一方下  
雙方皆無處可下時遊戲結束，依棋子數目決定輸贏或平手  

# 使用外部AI
程式可以導入外部AI，外部的AI程式須使用while input，並輸出結果  
範例：  
```
int main() {  
    string board;
    int color;
    while(cin >> board >> color) {
        ...
        cout << answer << endl;
    }
}
```
輸入```++++++++++++++OX++++XO++++++++++++++ 1```，輸出```Bc```  
(X表示黑方，O表示白方；1表示為黑方，2為白方)  
若顯示外部AI出錯，請到error.log查看詳細訊息  

# 自行編譯
require go 1.16+  
### windows
```go build -ldflags="-H windowsgui"```  
### linux or macOS
```go build```

# Cross Compilation
https://github.com/fyne-io/fyne-cross
