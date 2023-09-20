
void setup() {
  Serial.begin(115200);
  pinMode(2,INPUT);
  pinMode(3,INPUT);
  pinMode(4,INPUT);
  pinMode(5,INPUT);
  
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
  delay(20);
}
