#include <ESP8266WiFi.h>
#include <ESPAsyncTCP.h>
#include <ESP8266HTTPClient.h>
#include <ESPAsyncWebServer.h>

const char *ssid = "Dreamland";
const char *password = "1231231231";

const int ledPin = D5;
AsyncWebServer server(80);
bool state = false;

void setup() {
  Serial.begin(9600);
  pinMode(ledPin, OUTPUT);
  digitalWrite(ledPin, LOW);

  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(300);
    Serial.print(".");
  }
  Serial.println();
  Serial.println(WiFi.localIP());

  server.on("/device/toggle", HTTP_POST, [](AsyncWebServerRequest *request) {
    if (state) {
      digitalWrite(ledPin, LOW);
      request->send(200, "application/json", "{\"message\":\"turned off device\"}");
    } else {
      digitalWrite(ledPin, HIGH);
      request->send(200, "application/json", "{\"message\":\"turned on device\"}");
    }
    state = !state;
  });

  server.on("/device/state", HTTP_GET, [](AsyncWebServerRequest *request) {
    request->send(200, "application/json", state ? "{\"message\":\"device is on\"}" : "{\"message\":\"device is off\"}");
  });

  server.onNotFound([](AsyncWebServerRequest *request) {
    request->send(404, "application/json", "{\"message\":\"not found\"}");
  });

  server.begin();
}

void loop() {}
