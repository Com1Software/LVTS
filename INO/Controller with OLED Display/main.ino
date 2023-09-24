#include <Servo.h>
#include <SPI.h>
#include <Wire.h>
#include <Adafruit_GFX.h>
#include <Adafruit_SSD1306.h>

Adafruit_SSD1306 display = Adafruit_SSD1306(128, 32);

Servo myservo1; 
Servo myservo2; 
Servo myservo3; 

void setup() {
  Serial.begin(115200);
    display.begin(SSD1306_SWITCHCAPVCC, 0x3C); 
  display.clearDisplay();
  display.setTextSize(1);
  display.setTextColor(SSD1306_WHITE);
  display.setCursor(0,0);
  display.print("Throtle    Steering");
  display.setCursor(20,10);
  display.print("0");
  display.setCursor(80,10);
  display.print("0");
  display.display();
  delay(1000);
  pinMode(2,INPUT);
  pinMode(3,INPUT);
  pinMode(4,INPUT);
  pinMode(5,INPUT);
  myservo1.attach(9);
  myservo2.attach(10);
  myservo3.attach(11);
  
}
void loop() {
  Serial.print("CH1=");
  Serial.print(map(pulseIn(2,HIGH),1100,1900,0,255));
  Serial.print(",CH2=");
  Serial.print(map(pulseIn(3,HIGH),1100,1900,0,255));
  Serial.print(",CH3=");
  Serial.print(map(pulseIn(4,HIGH),1100,1900,0,255));
  Serial.print(",CH4=");
  Serial.print(map(pulseIn(5,HIGH),1100,1900,0,255));
  ;Serial.print("\n");
   myservo1.write(map(pulseIn(2,HIGH),1100,1900,0,255)); 
   myservo2.write(map(pulseIn(3,HIGH),1100,1900,0,255)-125); 
   myservo3.write(map(pulseIn(4,HIGH),1100,1900,0,255));
   display.clearDisplay();
   display.setCursor(0,0);
  display.print("Throtle    Steering");
  display.setCursor(20,10);
  display.print("     ");
  display.setCursor(20,10);
  display.print(map(pulseIn(2,HIGH),1100,1900,0,255));
  display.setCursor(80,10);
  display.print("     ");
  display.setCursor(80,10);
  display.print(map(pulseIn(3,HIGH),1100,1900,0,255));
  display.display();
   delay(10);
}
