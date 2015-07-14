//96x16
import processing.serial.*;

SerialSelector selector;
Circles circles = new Circles();
CircleController cc = new CircleController();

void setup(){
    selector = new SerialSelector(this);
    selector.show();
    //size(1345, 225);
    size(14*96, 14*16);
    background(255,255,255);
    //rect(60, 80, 240, 180);
}

void draw(){
    viewLCDDisplay();
    Serial port = selector.getSerial();
    if (port != null){
        //port.write("hoge");
        LCDController lcdCtl = new LCDController();
        lcdCtl.sendData(port, circles);

    }
    //String sendStr = "test string";
    //selected_port.write(100);
    /*
    if (mousePressed == true){
    //ellipse(mouseX, mouseY, 60, 60);
        try{
            CoordinateController coordCtrl = new CoordinateController();
            Coordinate coord = coordCtrl.convertR2V(mouseX, mouseY);
            color c = circles.matrix[coord.x][coord.y];
            c = c & 0x00ffffff;

            //red
            if (c == 0xff0000) {
                c = 0xff00ff00;
            }
            //green
            else if (c == 0x00ff00) {
                c = 0xffffA500;
            }
            //red + green
            else if (c == 0xffA500) {
                c = 0xff000000;
            }
            else if (c == 0x000000) {
                c = 0xffff0000;
            }
            circles.matrix[coord.x][coord.y] = c;
        }catch(Exception e){

        }
    }*/
}


void serialEvent(Serial port) {
    Serial selected_port = selector.getSerial();
    if ( port == selected_port ) {
        if ( port.available() > 0 ) {
            String buffer = selected_port.readString();
            if(buffer != null){
                //println("read: "+buffer);
            }
        }
    }

}
void mousePressed() {
    if (mousePressed == true){
    //ellipse(mouseX, mouseY, 60, 60);
        try{
            CoordinateController coordCtrl = new CoordinateController();
            Coordinate coord = coordCtrl.convertR2V(mouseX+6, mouseY+3);
            color c = circles.matrix[coord.x][coord.y];
            c = c & 0x00ffffff;

            //red
            if (c == 0xff0000) {
                c = 0xff00ff00;
            }
            //green
            else if (c == 0x00ff00) {
                c = 0xffffA500;
            }
            //red + green
            else if (c == 0xffA500) {
                c = 0xff000000;
            }
            else if (c == 0x000000) {
                c = 0xffff0000;
            }
            circles.matrix[coord.x][coord.y] = c;
        }catch(Exception e){

        }
    }
}

void viewLCDDisplay(){
    cc.draw(circles);
}
public class Coordinate{
    //仮想座標
    int x;
    int y;
    //実座標
    int realX;
    int realY;
    //単位量
    int u;
    //原点
    int ox;
    int oy;

    Coordinate(){
        this.u = 14;
        this.x = 0;
        this.y = 0;
        this.ox = this.u/2;
        this.oy = this.u/2;
        this.realX = this.u/2;
        this.realY = this.u/2;
    }

    int getX(){
        return this.x;
    }
    void setX(int x){
        this.x = x;
    }
    int getY(){
        return this.y;
    }
    void setY(int y){
        this.y = y;
    }
    int getRealX(){
        return this.realX;
    }
    void setRealX(int realX){
        this.realX = realX;
    }
    int getRealY(){
        return this.realY;
    }
    void setRealY(int realY){
        this.realY = realY;
    }
    int getU(){
        return this.u;
    }
    void setU(int u){
        this.u = u;
    }
    int getOX(){
        return this.ox;
    }
    void setOX(int ox){
        this.ox = ox;
    }
    int getOY(){
        return this.oy;
    }
    void setOY(int oy){
        this.oy = oy;
    }
}

public class CoordinateController{
    //仮想座標を実座標に変換する
    Coordinate convertV2R(int x, int y){
        Coordinate ret = new Coordinate();
        ret.setX(x);
        ret.setY(y);
        int rx = ret.getX() * ret.getU() + ret.getOX();
        int ry = ret.getY() * ret.getU() + ret.getOY();
        ret.setRealX(rx);
        ret.setRealY(ry);
        return ret;
    }
    //実座標から仮想座標へ変換
    Coordinate convertR2V(int x, int y){
        Coordinate ret = new Coordinate();
        ret.setRealX(x);
        ret.setRealY(y);
        int vx = (ret.getRealX() - ret.getOX()) / ret.getU();
        int vy = (ret.getRealY() - ret.getOY()) / ret.getU();

        ret.setX(vx);
        ret.setY(vy);
        return ret;
    }
}

public class Circles{
    color matrix[][] = new color[96][16];
    Circles(){
        init();
    }

    void init(){
        for (int y = 0; y < 16; y ++) {
            for (int x = 0; x < 96; x ++) {
                this.matrix[x][y] = 0xff000000;
            }
        }
    }
}

public class CircleController{
    void draw(Circles c){
        Coordinate coor;
        for(int y = 0;y < 16; y ++){
            for(int x = 0; x < 96; x ++){
                CoordinateController cc = new CoordinateController();
                coor = cc.convertV2R(x, y);
                //background(255,255,255);
                fill(c.matrix[x][y]);
                stroke(255,255,255);
                strokeWeight(3);
                //noStroke();
                ellipse(coor.getRealX(), coor.getRealY(), coor.getU(),coor.getU());
            }
        }
    }
}

public class LCDController{
    String header;
    String coord;
    String end;
    LCDController(){
        header = "pcmat\r";
        coord = "000\r5ff\r";
        end = "end\r";
    }



    byte[] createRdata(Circles c){
        byte ret[] = new byte[192];
        int p = 0;
        for(int y = 0;y < 16; y++){
            for(int x = 0;x < 96; x += 8){
                ret[p] = 0;
                for(int k = 0;k < 8;k ++){
                    ret[p] = (byte)((int)ret[p] << 1);
                    if ((c.matrix[x+k][y] & 0xff0000) != 0) {
                        ret[p] |= 1;
                    }
                }
                p++;
            }
        }
        return ret;
    }

    byte[] createGdata(Circles c){
        byte ret[] = new byte[192];
        int p = 0;
        for(int y = 0;y < 16; y++){
            for(int x = 0;x < 96; x += 8){
                ret[p] = 0;
                for(int k = 0;k < 8;k ++){
                    ret[p] = (byte)((int)ret[p] << 1);
                    if ((c.matrix[x+k][y] & 0x00ff00) != 0) {
                        ret[p] |= 1;
                    }
                }
                p++;
            }
        }
        return ret;
    }


    void sendData(Serial port, Circles c){
        port.write(this.header);
        port.write(this.coord);
        byte[] r = createRdata(c);
        byte[] g = createGdata(c);

        for(int i = 0;i < r.length; i ++){
            port.write(r[i]);
        }

        port.write('\r');
        for(int i = 0;i < g.length; i ++){
            port.write(g[i]);
        }
        port.write('\r');
        port.write(this.end);
    }
}
