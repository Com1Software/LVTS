#include <LiquidCrystal_I2C.h>
#include <Servo.h>

const int trigPin = 2;
const int echoPin = 3;
LiquidCrystal_I2C lcd(0x27,16,2);
long duration;
int distance;
Servo myservo1; 
int position;
void setup() {
  pinMode(trigPin, OUTPUT); 
  pinMode(echoPin, INPUT); 
  myservo1.attach(4);
  Serial.begin(9600); 
  lcd.init();
  lcd.backlight();
  }
void loop() {
  position++;
  if (position > 120) position = 0;
   myservo1.write(position);
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
  lcd.print("     ");
  lcd.setCursor(11,0);
  lcd.print(distance);
  lcd.setCursor(0,1);
  lcd.print("Position");
  lcd.setCursor(11,1);
  lcd.print("     ");
  lcd.setCursor(11,1);
  lcd.print(position);
  Serial.print("DIS=");
  Serial.println(distance);
  Serial.print("POS=");
  Serial.println(position);
}
