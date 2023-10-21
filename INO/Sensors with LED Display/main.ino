#include <LiquidCrystal_I2C.h>

const int trigPin = 9;
const int echoPin = 10;
LiquidCrystal_I2C lcd(0x27,16,2);
long duration;
int distance;
int position;
void setup() {
  pinMode(trigPin, OUTPUT); 
  pinMode(echoPin, INPUT); 
  Serial.begin(9600); 
  lcd.init();
  lcd.backlight();
  }
void loop() {
   digitalWrite(trigPin, LOW);
  delayMicroseconds(2);
  digitalWrite(trigPin, HIGH);
  delayMicroseconds(10);
  digitalWrite(trigPin, LOW);
  duration = pulseIn(echoPin, HIGH);
  distance = duration * 0.034 / 2;
  lcd.setCursor(0,0);
  lcd.print("Distance");
  lcd.setCursor(11,0);
  lcd.print(distance);
  lcd.setCursor(0,1);
  lcd.print("Position");
  lcd.setCursor(11,1);
  lcd.print(position);
  Serial.print("DIS=");
  Serial.println(distance);
  Serial.print("POS=");
  Serial.println(distance);
}
