"use client";

import * as Y from "yjs";

export class YjsWebsocketProvider {
  doc: Y.Doc;
  ws: WebSocket;
  room: string;

  constructor(room: string, url: string, doc: Y.Doc) {
    this.room = room;
    this.doc = doc;
    this.ws = new WebSocket(url);

    //  connect to sync the current state
    this.ws.onopen = () => {
      console.log("Connected to the Go Websocket for", room);

      const update = Y.encodeStateAsUpdate(this.doc);
      this.ws.send(update);
    };

    // on receiving any message form websocket it updates its own page  for ithers changes
    this.ws.onmessage = (event) => {
      const data = new Uint8Array(event.data);
      Y.applyUpdate(this.doc, data);
    };

    // pushing the update to sockets for its own made changes
    this.doc.on("update", (update: Uint8Array) => {
      this.ws.send(update);
    });
  }
}
