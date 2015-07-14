//96x16
import processing.serial.*;
import java.awt.Color;
import javax.swing.*;
import java.awt.Graphics;
import java.awt.Graphics2D;
import java.awt.image.BufferedImage;
import java.io.File;
import javax.imageio.ImageIO;
import java.awt.Font;
import java.awt.RenderingHints;
import java.io.*;
import java.awt.*;
import java.awt.geom.*;
import java.awt.image.*;
import javax.imageio.*;
import java.net.URL;
import java.net.HttpURLConnection;
import java.io.InputStreamReader;
import java.io.IOException;
import java.io.BufferedReader;



SerialSelector selector;
Circles circles = new Circles();
CircleController cc = new CircleController();
String writeStr = "";

double stringCoord = 14 * 96 + 20;
//News news = new News();
JSONObject js = loadJSONObject("https://www.kimonolabs.com/api/bnl617tc?apikey=t36hExBGRYVtIATai55sBsahHkXdlt1v");
JSONObject results = js.getJSONObject("results");
JSONArray collection1 = results.getJSONArray("collection1");
int newsNum = 0;
int frameFlag = 0;
int maiden = 0;


void setup(){
    selector = new SerialSelector(this);
    frameRate(20);
    selector.show();
    //size(1345, 225);
    size(14*96, 14*16);
    background(255,255,255);
    new BitmapStrings().Create(writeStr);
    //rect(60, 80, 240, 180);
}

void draw(){
    JSONObject collection = collection1.getJSONObject(newsNum);
    JSONObject news = collection.getJSONObject("news");
    writeStr = news.getString("text");
    new BitmapStrings().Create(writeStr);
    stringCoord -= 28.0;
    int charSize = (14 * 96) / 6;
    //if (stringCoord < (charSize * writeStr.length() * -1 + writeStr.length() * 10)){
    if(circles.allColorOr == 0 && maiden != 0){
        stringCoord = 14 * 96;
        newsNum ++;
        maiden = 0;
        if(collection1.size() == newsNum){
            newsNum = 0;
        }
    }
    viewLCDDisplay();
    Serial port = selector.getSerial();
    if (port != null){
        //port.write("hoge");
        LCDController lcdCtl = new LCDController();
        lcdCtl.sendData(port, circles);

    }
    BitmapStrings b = new BitmapStrings();

    if (frameFlag % 3 == 0){
        Circles c = b.convertImage2Matrix("/Users/masato/git/aoi_shirase/matrix_led/cmd/ui/test.jpg");
        circles = c;
    }
    else {
        if (frameFlag > 65536) {
                frameFlag = 0;
        }
        Circles c = new Circles();
        circles = c;
    }
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
    int allColorOr;
    Circles(){
        allColorOr = 0;
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

public class News{
    String url;
    //http://appli.ntv.co.jp/ntv_WebAPI/news/?key=13NfqB5gW46kFmZAyp6Jjj4HSrL27pK1Nj7Bxv5jXDftzL7yzNQKuoDR7fv7&word=news
    //https://www.kimonolabs.com/api/bnl617tc?apikey=t36hExBGRYVtIATai55sBsahHkXdlt1v
    String[] json;
    News(){
        url = "https://www.kimonolabs.com/api/bnl617tc?apikey=t36hExBGRYVtIATai55sBsahHkXdlt1v";
        json = loadStrings(url);
    }

    String getUrl(){
        return this.url;
    }
    void setUrl(String url){
        this.url = url;
    }

    String[] getJson(){
        return this.json;
    }
    void setJson(String[] json){
        this.json = json;
    }
}

public class NewsController{

}

public class BitmapStrings{
    public void Create(String str) {
        int w=14*96;
        int h=14*16;
        try {
            //受け取った文字列を画像化
            BufferedImage image=new BufferedImage(w,h,BufferedImage.TYPE_INT_RGB);
            Graphics2D g2d=image.createGraphics();
            Font font = new Font(Font.DIALOG_INPUT, Font.PLAIN, 230);
            //Font font = new Font(Font.DIALOG_INPUT, Font.ITALIC, 200);
            g2d.setFont(font);
            g2d.setBackground(Color.WHITE);
            g2d.clearRect(0,0,w,h);
            g2d.setColor(Color.BLACK);
            g2d.drawString(str,0 + (int)stringCoord,h-20);

            ImageIO.write(image, "JPEG", new File("/Users/masato/git/aoi_shirase/matrix_led/cmd/ui/test.jpg"));
        } catch(Exception e) {
            e.printStackTrace();
        }
    }


    public Circles convertImage2Matrix(String fileName){
        try{
            BufferedImage img = ImageIO.read(new File(fileName));
            img = convertBinaryImg(img);
            int w = img.getWidth();
            int h = img.getHeight();
            Coordinate coord = new Coordinate();
            int buffer = 0xffffffff;
            Circles c = new Circles();
            int margin = 5;
            c.allColorOr = 0;

            for (int y = 0;y < h; y += coord.getU()){
                for (int x = 0; x < w; x += coord.getU()){
                    for(int by = y + margin; by < y + coord.getU() - margin; by ++){
                        for(int bx = x + margin; bx < x + coord.getU() - margin; bx ++){
                            buffer = buffer & img.getRGB(bx, by);
                        }
                    }
                    buffer ^= 0xffffffff;
                    buffer |= 0xff000000;
                    buffer &= 0xffff5a00;
                    c.matrix[x/coord.getU()][y/coord.getU()] = buffer;
                    c.allColorOr |= buffer;
                    c.allColorOr &= 0x00ffffff;
                    if(maiden == 0){
                        maiden = c.allColorOr;
                    }
                    buffer = 0xffffffff;
                }
            }
            return c;
        }
        catch (Exception e) {
            return null;
        }
    }
    public BufferedImage convertBinaryImg(BufferedImage img)throws Exception{
        WritableRaster wr = img.getRaster();
        int buf[] = new int[wr.getNumDataElements()];
        for(int ly=0;ly<wr.getHeight();ly++){
            for(int lx=0;lx<wr.getWidth();lx++){
                wr.getPixel(lx, ly, buf);

                int maxval = Math.max(Math.max(buf[0], buf[1]), buf[2]);
                int minval = Math.min(Math.min(buf[0], buf[1]), buf[2]);
                buf[0] = buf[1] = buf[2] = (maxval+minval)/2;

                wr.setPixel(lx, ly, buf);
            }
        }
        /* lookupデータ作成 */
        byte dat[] = new byte[256];
        for(int di=0;di<256;di++){
            dat[di] = di>256*0.55?(byte)255:(byte)0;
        }
        LookupOp lo = new LookupOp(new ByteLookupTable(0, dat), null);
        BufferedImage img2 = lo.filter(img, null);
        return img2;
    }
}

public class ImageUtility{
    public int a(int c){
        return c>>>24;
    }
    public int r(int c){
        return c>>16&0xff;
    }
    public int g(int c){
        return c>>8&0xff;
    }
    public int b(int c){
        return c&0xff;
    }
    public int rgb
    (int r,int g,int b){
        return 0xff000000 | r <<16 | g <<8 | b;
    }
    public int argb
    (int a,int r,int g,int b){
        return a<<24 | r <<16 | g <<8 | b;
    }
}
