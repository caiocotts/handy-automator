#include <ESP8266WiFi.h>
#include <ESPAsyncTCP.h>
#include <ESP8266HTTPClient.h>
#include <ESPAsyncWebServer.h>
#include <ESP8266mDNS.h>

const char *ssid = "Dreamland";
const char *password = "1231231231";

const int ledPin = D5;
AsyncWebServer server(80);
bool state = false;

String mdnsHostname() {
  return "handy-" + String(ESP.getChipId(), HEX);
}

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

  String hostname = mdnsHostname();
  if (MDNS.begin(hostname)) {
    MDNS.addService("handy-automator", "tcp", 80);
    Serial.println("mDNS started: " + hostname + ".local");
  } else {
    Serial.println("mDNS failed to start");
  }

  server.on("/device/toggle", HTTP_POST, [](AsyncWebServerRequest *request) {
    if (state) {
      digitalWrite(ledPin, LOW);
      request->send(200, "application/json", "{\"state\":\"off\"}");
    } else {
      digitalWrite(ledPin, HIGH);
      request->send(200, "application/json", "{\"state\":\"on\"}");
    }
    state = !state;
  });

  server.on("/device/state", HTTP_GET, [](AsyncWebServerRequest *request) {
    request->send(200, "application/json", state ? "{\"state\":\"on\"}" : "{\"state\":\"off\"}");
  });

  server.onNotFound([](AsyncWebServerRequest *request) {
    request->send(404, "application/json", "{\"message\":\"not found\"}");
  });

  server.begin();
}

void loop() {
  MDNS.update();
}
