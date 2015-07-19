# MatrixLEDGUI #################################
![RoboLive](https://github.com/RobotClubKut/MatrixLEDGUI/blob/master/img/robolive.png?raw=true)
![RoboLive](https://github.com/RobotClubKut/MatrixLEDGUI/blob/master/img/IMG_5362.jpg?raw=true)
##概要
* Dot Matrixを表示させるためのプロジェクト
* Dot Matrixの詳細は[MatrixLED](https://github.com/RobotClubKut/MatrixLED)を参照すること

##Java版の使用方法
* MatrixLEDのためのGUI

###コマンド一覧
|Command|                  |
|:-----:|------------------|
|s      |画像のスクロール/停止|
|c      |色の変更           |
|u      |frameRate up      |
|d      |frameRate down    |
|j      |左スクロール        |
|k      |右スクロール        |
|Return |描画の停止. 送信される画像は動く|
|ESC    |プログラムの終了     |

####注意事項
* デフォルトでは描画をしないのでReturnを押すこと
* ControlP5ライブラリを入れること

##go版の使用方法
* GO版にはGUIなんて存在しない

###コマンド一覧
* 起動オプションは未実装

###起動方法
1. ラズパイの電源を入れる
    * 電源を入れるときはUSB USERTは抜いておく
    * 刺したままだとうまく起動しない場合がある
2. SSHでloginする
3. $ ./main
4. fontの選択
    * 一覧が表示されるので表示に対応したfontを選択してreturn
5. serial portの選択
    * 一覧が表示されるので表示に対応したportを選択してreturn
    * デバイスが刺さってない時は表示されないので注意

###Web api
|URL    | 機能                      |例 |
|:-----:|:------------------------:|:-:|
|/      |Hello, world              |   |
|/update|strで表示文字, colで色の選択 |/update?str=にゃんぱす&col=red|
