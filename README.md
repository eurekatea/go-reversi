# othello-board
海大1091 人工智慧 黑白棋前端界面

# 遊戲截圖
![alt 文字](https://raw.githubusercontent.com/lemon37564/othello-board/main/game/img/screenshot.webp "Logo 標題文字 1")

# 使用方式
在engine資料夾下放AI1.exe為先攻，執黑子，若沒有此檔案則先攻為人類玩家  
同理AI2.exe執白子後攻，無此檔案則為人類玩家  
linux平台下請命名為AI1與AI2

# 規則
規則與一般黑白棋相同，版面大小為6x6，若有一方無處可下會自動PASS，換另一方下  
雙方皆無處可下時遊戲結束，依棋子數目決定輸贏或平手  

# 說明
此程式會將版面以"++++++++++++++OX++++XO++++++++++++++ 1"形式傳送給AI之stdin  
(X表示黑方，O表示白方；1表示現在為黑方行動，2反之)  
程式在AI之stdout接收回傳值，若回傳值不合法，則在log file內印出錯誤訊息並退出  
