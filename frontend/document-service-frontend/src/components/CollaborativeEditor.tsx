import React, { useEffect, useRef, useState } from "react";
import * as Y from "yjs";
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Collaboration from "@tiptap/extension-collaboration";

type Props = {
  readonly docId: string;
  readonly wsUrl?: string; // optional
};

export default function CollaborativeEditor({ docId, wsUrl }: Props) {
  // Create Y.Doc synchronously so useEditor sees it on first render
  const ydocRef = useRef<Y.Doc>(new Y.Doc());
  const wsRef = useRef<WebSocket | null>(null);
  const originRef = useRef<object | null>(null);

  const [connectionStatus, setConnectionStatus] = useState<
    "connecting" | "connected" | "error"
  >("connecting");
  const [errorMessage, setErrorMessage] = useState<string>("");

  // Load snapshot and apply BEFORE starting websocket (runs when docId changes)
  useEffect(() => {
    const ydoc = ydocRef.current;

    // Reset doc when docId changes: destroy old and create fresh doc
    ydoc.destroy();
    const newDoc = new Y.Doc();
    ydocRef.current = newDoc;

    fetch(`http://localhost:7000/documents/${docId}/snapshot`)
      .then(async (res) => {
        if (!res.ok) {
          console.log("No snapshot found. Starting fresh.");
          return;
        }
        const buf = await res.arrayBuffer();
        Y.applyUpdate(newDoc, new Uint8Array(buf));
        console.log("Snapshot loaded!");
      })
      .catch((err) => console.error("Snapshot load error:", err));

    // create a stable origin object for this doc/connection
    originRef.current = {};

    return () => {
      // cleanup when docId changes/unmount
      newDoc.destroy();
      originRef.current = null;
    };
  }, [docId]);

  // Autosave every 5 seconds
  useEffect(() => {
    const ydoc = ydocRef.current;
    if (!ydoc) return;

    const interval = setInterval(() => {
      try {
        const update = Y.encodeStateAsUpdate(ydoc);
        fetch(`http://localhost:7000/documents/${docId}/snapshot`, {
          method: "POST",
          headers: { "Content-Type": "application/octet-stream" },
          body: update,
        }).catch((err) => {
          console.error("Failed to autosave snapshot:", err);
        });
        console.log("Autosaved snapshot!");
      } catch (err) {
        console.log("Autosave error:", err);
      }
    }, 5000);

    return () => clearInterval(interval);
  }, [docId]);

  // Set up WebSocket connection (connect AFTER snapshot load attempt)
  useEffect(() => {
    const ydoc = ydocRef.current;
    if (!ydoc) return;

    const url =
      wsUrl ?? `ws://localhost:7000/ws/document/${encodeURIComponent(docId)}`;

    console.log("Connecting to WebSocket:", url);
    setConnectionStatus("connecting");
    setErrorMessage("");

    const ws = new WebSocket(url);
    ws.binaryType = "arraybuffer";
    wsRef.current = ws;

    // create an origin object for this connection (stable while ws open)
    const localOrigin = {};
    originRef.current = localOrigin;

    ws.addEventListener("open", () => {
      console.log("WebSocket connected");
      setConnectionStatus("connected");
      setErrorMessage("");
    });

    ws.addEventListener("error", (error) => {
      console.error("WebSocket error:", error);
      setConnectionStatus("error");
      setErrorMessage("Failed to connect to WebSocket server.");
    });

    ws.addEventListener("close", (event) => {
      console.log("WebSocket closed:", event.code, event.reason);
      if (event.code !== 1000) {
        setConnectionStatus("error");
        setErrorMessage("Connection closed unexpectedly.");
      } else {
        setConnectionStatus("connecting"); // or whatever you prefer
      }
    });

    // local updates -> send to server (don't re-broadcast updates that came from this ws)
    const onUpdate = (update: Uint8Array, origin: any) => {
      if (origin === localOrigin) return; // don't send writes that originated from this ws
      if (ws.readyState === WebSocket.OPEN) {
        try {
          ws.send(update);
        } catch (e) {
          console.error("Error sending update:", e);
        }
      }
    };
    ydoc.on("update", onUpdate);

    // incoming updates from server
    const onMessage = (ev: MessageEvent) => {
      const data = ev.data;
      if (data instanceof ArrayBuffer) {
        const update = new Uint8Array(data);
        // Apply update using localOrigin so onUpdate knows this came from remote and won't re-send
        Y.applyUpdate(ydoc, update, localOrigin);
      } else if (data instanceof Blob) {
        data.arrayBuffer().then((buf) => {
          Y.applyUpdate(ydoc, new Uint8Array(buf), localOrigin);
        });
      } else {
        console.warn("Received non-binary WS message:", typeof data);
      }
    };
    ws.addEventListener("message", onMessage);

    return () => {
      ydoc.off("update", onUpdate);
      ws.removeEventListener("message", onMessage);
      if (
        ws.readyState === WebSocket.OPEN ||
        ws.readyState === WebSocket.CONNECTING
      ) {
        ws.close();
      }
      wsRef.current = null;
      originRef.current = null;
    };
  }, [docId, wsUrl]);

  // Initialize TipTap editor with the Y.Doc (ydocRef is synchronous)
  const editor = useEditor({
    extensions: [
      StarterKit,
      Collaboration.configure({
        document: ydocRef.current,
      }),
    ],
    content:
      connectionStatus === "connected"
        ? "<p>Start typing...</p>"
        : "<p>Connecting...</p>",
    editable: connectionStatus === "connected",
  });

  useEffect(() => {
    if (editor) editor.setEditable(connectionStatus === "connected");
  }, [editor, connectionStatus]);

  // ... UI unchanged
  const getStatusColor = () => {
    if (connectionStatus === "connected") return "#22c55e";
    if (connectionStatus === "error") return "#ef4444";
    return "#f59e0b";
  };

  const getStatusText = () => {
    if (connectionStatus === "connected") return "✓ Connected";
    if (connectionStatus === "error") return "✗ Error";
    return "⟳ Connecting...";
  };

  return (
    <div style={{ border: "1px solid #ddd", padding: 12, borderRadius: 8 }}>
      <div style={{ marginBottom: 8, fontSize: "0.875rem", color: "#666" }}>
        Status:{" "}
        <span style={{ color: getStatusColor() }}>{getStatusText()}</span>
        {errorMessage && (
          <div style={{ color: "#ef4444", marginTop: 4 }}>{errorMessage}</div>
        )}
      </div>
      <EditorContent editor={editor} />
    </div>
  );
}
