#include <WiFi.h>
#include <WiFiClient.h>
#include <WebServer.h>
#include <ESPmDNS.h>
#include <Update.h>

WebServer server(80);

const char *serverIndex =
    "<script "
    "src='https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js'></"
    "script>"
    "<h1>Upload a new .bin file</h1>"
    "<form method='POST' action='#' enctype='multipart/form-data'"
    "id='upload_form'>"
    "<input type='file' name='update'"
    ">"
    "<input type='submit' value='Update'>"
    "</form>"
    "<div id='prg'>Progress: 0%</div>"
    "<style>"
    "body {"
    "width: 100%;"
    "height: 100%;"
    "display: flex;"
    "flex-direction: column;"
    "font-family: 'Lucida Console', 'Courier New', monospace;"
    "align-items: center;"
    "justify-content: center;"
    "gap: 70px;"
    "}"
    "input[type=button], input[type=submit] {"
    "border: none;"
    "padding: 16px 32px;"
    "text-decoration: none;"
    "margin: 4px 2px;"
    "cursor: pointer;"
    "border-radius: 4px;"
    "}"
    "form {"
    "padding: 50px;"
    "box-shadow: 4px 3px 4px lightgrey;"
    "color: white;"
    "background-color: #313131;"
    "border-radius: 10px;"
    "}"
    "#prg {"
    "font-size: 50px;"
    "}"
    "</style>"
    "<script>"
    "$('form').submit(function(e){"
    "e.preventDefault();"
    "var form = $('#upload_form')[0];"
    "var data = new FormData(form);"
    " $.ajax({"
    "url: '/update',"
    "type: 'POST',"
    "data: data,"
    "contentType: false,"
    "processData:false,"
    "xhr: function() {"
    "var xhr = new window.XMLHttpRequest();"
    "xhr.upload.addEventListener('progress', function(evt) {"
    "if (evt.lengthComputable) {"
    "var per = evt.loaded / evt.total;"
    "$('#prg').html('progress: ' + Math.round(per*100) + '%');"
    "}"
    "}, false);"
    "return xhr;"
    "},"
    "success:function(d, s) {"
    "console.log('success!')"
    "},"
    "error: function (a, b, c) {"
    "}"
    "});"
    "});"
    "</script>";

void initOtaUpdate(String ownID) {
  if (!MDNS.begin(ownID)) {  // http://63D592A9.local
    Serial.println("Error setting up MDNS responder!");
    while (1) {
      delay(1000);
    }
  }
  Serial.println("mDNS responder started");
  server.on("/", HTTP_GET, []() {
    server.sendHeader("Connection", "close");
    server.send(200, "text/html", serverIndex);
  });

  server.on(
      "/update", HTTP_POST,
      []() {
        server.sendHeader("Connection", "close");
        server.send(200, "text/plain", (Update.hasError()) ? "FAIL" : "OK");
        ESP.restart();
      },
      []() {
        HTTPUpload &upload = server.upload();
        if (upload.status == UPLOAD_FILE_START) {
          Serial.printf("Update: %s\n", upload.filename.c_str());
          if (!Update.begin(
                  UPDATE_SIZE_UNKNOWN)) {  // start with max available size
            Update.printError(Serial);
          }
        } else if (upload.status == UPLOAD_FILE_WRITE) {
          /* flashing firmware to ESP*/
          if (Update.write(upload.buf, upload.currentSize) !=
              upload.currentSize) {
            Update.printError(Serial);
          }
        } else if (upload.status == UPLOAD_FILE_END) {
          if (Update.end(
                  true)) {  // true to set the size to the current progress
            Serial.printf("Update Success: %u\nRebooting...\n",
                          upload.totalSize);
          } else {
            Update.printError(Serial);
          }
        }
      });
  server.begin();
}