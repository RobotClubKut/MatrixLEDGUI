import processing.serial.*;
import controlP5.*;
import javax.swing.*; 
import java.awt.Insets;


public class SerialSelector extends JFrame {
  final String TITLE = "Serial Selector";
  final int WIDTH = 540;
  final int HEIGHT = 145;
  final int DEFAULT_POSX = 200;
  final int DEFAULT_POSY = 200;  
  final int DEFAULT_PORT_INDEX = 0;
  final int DEFAULT_RATE_INDEX = 4;  // this means rates[4] = 9600;

  SecondApplet s;
  PApplet parent;  
  Serial port;  
  String[] ports; 
  final int[] rates = {
    300, 1200, 2400, 4800, 9600, 14400, 19200, 28800, 38400, 57600, 115200
  };
  
  String default_port;
  int    default_rate;
  int port_index = DEFAULT_PORT_INDEX;
  int rate_index = DEFAULT_RATE_INDEX;

  public SerialSelector(PApplet _applet) {
    initThis( _applet, "dummy", -1 );    
  }

  public SerialSelector(PApplet _applet, String _default_port, int _default_rate ) {
    initThis( _applet, _default_port, _default_rate );    
  }

  void initThis(PApplet _parent, String _default_port, int _default_rate ) {
    this.parent = _parent;
    this.default_port = _default_port;
    this.default_rate = _default_rate;

    // setting of applet
    s = new SecondApplet();
    add(s);
    s.init();

    // setting of frame
    Insets insets = frame.getInsets();    
    setSize( WIDTH + insets.left + insets.right, HEIGHT + insets.top + insets.bottom );
    setResizable(false);
    setLocation(DEFAULT_POSX, DEFAULT_POSY);
    setTitle(TITLE);

    // Get serial list
    ports = Serial.list();

    // Find index by default parameters
    for ( int i=0; i<ports.length; i++) {
      if (ports[i].equals(default_port)) {
        port_index = i;
        break;
      }
    }
    for (int i=0; i<rates.length; i++) {
      if (rates[i]==default_rate) {
        rate_index = i;
        break;
      }
    }
  }



  public String getPort() {
    return ports[port_index];
  }

  public int getRate() {
    return rates[rate_index];
  }
  
  public Serial getSerial() {
    return port;
  }

  private class SecondApplet extends PApplet {
    final String FONT = "Verdata";
    final int FONT_SIZE = 20;
    final int dl_FONT_SIZE = 15;

    PFont myFont;
    ControlP5 cp5;
    DropdownList dl_ports;
    DropdownList dl_rates;   

    public void setup() {

      // Font setting
      myFont = parent.createFont(FONT, FONT_SIZE);
      textFont(myFont);
      textSize(FONT_SIZE);

      // setting of controlP5
      cp5 = new ControlP5(this);
      cp5.setControlFont(new ControlFont(myFont, dl_FONT_SIZE)); 

      // setting of dropdownlist (Serial port)
      dl_ports = cp5.addDropdownList("dl-serial-ports");
      if ( ports.length <= 0 ) { 
        dl_ports.captionLabel().set("not found");
      } else {
        dl_ports.captionLabel().set("serial port");
      }
      dl_ports.setScrollbarWidth(20);
      dl_ports.setPosition(14, 61);
      dl_ports.setWidth(330);  
      dl_ports.setHeight(80);
      dl_ports.setBarHeight(21);
      dl_ports.setItemHeight(22);
      dl_ports.valueLabel().style().marginTop = 3;
      dl_ports.toUpperCase(false);
      for (int i=0; i<ports.length; i++) {
        dl_ports.addItem(ports[i], i);
      }
      dl_ports.setIndex(port_index);


      // setting of dropdownlist (Baud rate)
      dl_rates = cp5.addDropdownList("dl-baud-rates");
      dl_rates.captionLabel().set("baud rate");
      dl_rates.setScrollbarWidth(20);
      dl_rates.setPosition(357, 61);
      dl_rates.setWidth(170);
      dl_rates.setHeight(80);  
      dl_rates.setBarHeight(21);
      dl_rates.setItemHeight(22);  
      dl_rates.toUpperCase(false);
      dl_rates.valueLabel().style().marginTop = 3;
      for (int i=0; i<rates.length; i++) {
        dl_rates.addItem( rates[i] + " bps", i);
      }
      dl_rates.setIndex(rate_index);
    }


    public void draw() {
      background(200);
      fill(0);
      text("Serial port:", 14, 28 );
      text("Baud rate:", 357, 28 );
    }    

    void controlEvent(ControlEvent event) {
      if (event.isGroup()) {
        if ( dl_ports != null) {
          port_index = (int)(dl_ports.getValue());
        }
        if (dl_rates != null) {
          rate_index = (int)(dl_rates.getValue());
        }
        if (port != null) {
          port.stop();
        }
        try {
          // println("Selected: " + ports[port_index] + ", " + rates[rate_index] );
          port = new Serial(parent, Serial.list()[port_index], rates[rate_index]);
        } 
        catch (Exception e) {
          port = null;
        }
      }
    }
  }
}

