import controlP5.*;
import mqtt.*;

color white = #f1f1f1;
color dark = #e1e1e1;
color blue = #5175F3;
color dblue = #224AD5;
color lgrey = #464646;
color dgrey = #0D0D0D;

MQTTClient client;
ControlP5 cp5;

ArrayList<Barmband> barmbands;

PFont trispace;

class Barmband {
  String id;
  PVector pos;
  ArrayList<Barmband> targets;
  color col;
  String state;
  int size;
  MultiList context;
  boolean contextVis = false;

  Barmband(String i, PVector p, ArrayList<Barmband> t, color c, String s) {
    id = i;
    pos = p;
    targets = t;
    col = c;
    state = s;
    size = 100;

    /*
    ControlFont font = new ControlFont(trispace, 16);
     
     context = cp5.addMultiList("context", int(pos.x), int(pos.y), 110, 35);
     // create a multiListButton which we will use to
     // add new buttons to the multilist
     MultiListButton b;
     b = context.add("Connect", 0);
     b.setColorBackground(blue);
     b.setColorForeground(dblue);
     b.setFont(font);
     
     b = context.add("Remove", 1);
     b.setColorForeground(dblue);
     b.setColorBackground(blue);
     b.setFont(font);
     */
  }

  void update() {
    fill(col);
    ellipse(pos.x, pos.y, size, size);
    text(id, pos.x - 50, pos.y + 70);
    fill(white);
    ellipse(pos.x, pos.y, size - 30, size - 30);

    //context.setPosition(pos.x + size, pos.y - 100);
  }

  void reposition() {
    if (mouseX > pos.x - size/2 && mouseX < pos.x + size/2 && mouseY > pos.y - size/2 && mouseY < pos.y + size/2 && mousePressed) {
      if (mouseButton == LEFT) {
        pos = new PVector(mouseX, mouseY);
        noFill();
        stroke(lgrey, 80);
        rect(pos.x, pos.y, 200, 200);
        noStroke();
      }
      if (mouseButton == RIGHT) {
      }
    }
  }
}

void setup() {
  size(800, 500);
  windowResizable(true);
  noStroke();
  rectMode(CENTER);
  client = new MQTTClient(this);
  client.connect("mqtt://localhost");
  trispace = createFont("data/trispace.ttf", 20);
  ControlFont font = new ControlFont(trispace, 20);

  cp5 = new ControlP5(this);

  barmbands = new ArrayList<Barmband>();

  cp5.addButton("AddBarmband")
    .setColorBackground(lgrey)
    .setColorForeground(dgrey)
    .setBroadcast(false)
    .setFont(font)
    .setCaptionLabel("New Barmband")
    .setPosition(width/2 - 100, height - 70)
    .setSize(200, 50)
    .setBroadcast(true)
    ;
}

void draw() {
  background(white);
  for (Barmband b : barmbands) {
    b.update();
    b.reposition();
  }
}

void keyPressed() {
  if (key == 'b') {
    println(barmbands.size());
  }
  if (key == 'm') {
    client.publish("barmband/test", "I'm processing");
  }
}

void clientConnected() {
  println("client connected");
  client.subscribe("barmband/setup");
  client.subscribe("barmband/challenge");
}

void messageReceived(String topic, byte[] payload) {
  String[] p = split(new String(payload), ' ');

  if (topic.equals("barmband/setup")) {
    println("Barmband with ID: " + p[1]);

    boolean canBeAdded = true;

    for (Barmband bb : barmbands) {
      if (p[1].equals(bb.id)) {
        canBeAdded = false;
        break;
      }
    }

    if (canBeAdded) {
      Barmband b = new Barmband(p[1], new PVector(random(0, width), random(0, height)), barmbands, color(#000000), "inactive");
      barmbands.add(b);
    }
  }

  if (topic.equals("barmband/challenge")) {
    switch(p[0]) {
    case "New":
      for (Barmband bb : barmbands) {
        if (p[2].equals(bb.id) || p[3].equals(bb.id)) {
          String red = str(p[4].charAt(0)) + str(p[4].charAt(1));
          String green = str(p[4].charAt(2)) + str(p[4].charAt(3));
          String blue = str(p[4].charAt(4)) + str(p[4].charAt(5));
          color newCol = unhex("FF" + green + red + blue);
          bb.col = color(newCol);
          println("changed " + bb.id + " to " + newCol);
        }
      }
      break;

    case "Pair":
    case "Abort":
      for (Barmband bb : barmbands) {
        bb.col = #000000;
      }
      break;
    }
  }
}


void connectionLost() {
  println("connection lost");
}

public void AddBarmband() {
  // this should broadcast a new mqtt message with the actual ID
  Barmband b = new Barmband("Test" + str(random(0, 100)), new PVector(random(0 + 200, width - 200), random(0 + 200, height - 200)), barmbands, color(#000000), "inactive");
  barmbands.add(b);
}

public void PrintBarmbands() {
  for (Barmband b : barmbands) {
    println(b.id);
  }
}
