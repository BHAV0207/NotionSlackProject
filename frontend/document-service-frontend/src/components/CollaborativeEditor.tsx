import React, { useEffect, useRef, useState } from "react";
import * as Y from "yjs";
import { useEditor, EditorContent } from "@tiptap/react";
import StarterKit from "@tiptap/starter-kit";
import Collaboration from "@tiptap/extension-collaboration";

type Props = {
  readonly docId: string;
  readonly wsUrl?: string; // optional, defaults to ws://localhost:7000/ws/document/{docId}
};

export default function CollaborativeEditor({ docId, wsUrl }: Props) {
  const ydocRef = useRef<Y.Doc | null>(null);
  const wsRef = useRef<WebSocket | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<
    "connecting" | "connected" | "error"
  >("connecting");
  const [errorMessage, setErrorMessage] = useState<string>("");

  // Create Y.Doc first
  useEffect(() => {
    const ydoc = new Y.Doc();
    ydocRef.current = ydoc;

    return () => {
      ydoc.destroy();
      ydocRef.current = null;
    };
  }, []);

  // Set up WebSocket connection
  useEffect(() => {
    if (!ydocRef.current) return;

    const ydoc = ydocRef.current;
    const url =
      wsUrl ?? `ws://localhost:7000/ws/document/${encodeURIComponent(docId)}`;

    console.log("Connecting to WebSocket:", url);
    setConnectionStatus("connecting");
    setErrorMessage("");

    const ws = new WebSocket(url);
    ws.binaryType = "arraybuffer";
    wsRef.current = ws;

    // Handle WebSocket open
    ws.addEventListener("open", () => {
      console.log("WebSocket connected");
      setConnectionStatus("connected");
      setErrorMessage("");
    });

    // Handle WebSocket errors
    ws.addEventListener("error", (error) => {
      console.error("WebSocket error:", error);
      setConnectionStatus("error");
      setErrorMessage(
        "Failed to connect to WebSocket server. Make sure the server is running."
      );
    });

    // Handle WebSocket close
    ws.addEventListener("close", (event) => {
      console.log("WebSocket closed:", event.code, event.reason);
      if (event.code !== 1000) {
        // Not a normal closure
        setConnectionStatus("error");
        setErrorMessage(
          "Connection closed unexpectedly. Trying to reconnect..."
        );
      }
    });

    // When local Y.Doc emits updates, send them to the server as binary frames
    const onUpdate = (update: Uint8Array, origin: any) => {
      // Don't send updates that came from the server (avoid echo)
      if (origin === ws) {
        return;
      }

      // send raw binary update only when WebSocket is open
      if (ws.readyState === WebSocket.OPEN) {
        try {
          ws.send(update);
        } catch (e) {
          console.error("Error sending update:", e);
        }
      }
    };

    ydoc.on("update", onUpdate);

    // When receiving a binary message, apply it as a Yjs update
    ws.addEventListener("message", (ev) => {
      const data = ev.data;
      // Only handle ArrayBuffer / binary frames
      if (data instanceof ArrayBuffer) {
        const update = new Uint8Array(data);
        // Apply remote updates into the local Y.Doc
        // Pass the WebSocket as origin to prevent echo
        Y.applyUpdate(ydoc, update, ws);
      } else if (data instanceof Blob) {
        // just in case server sends blob
        data.arrayBuffer().then((buf) => {
          Y.applyUpdate(ydoc, new Uint8Array(buf), ws);
        });
      } else {
        console.warn("Received non-binary message:", typeof data);
      }
    });

    // Cleanup on unmount
    return () => {
      ydoc.off("update", onUpdate);
      if (
        ws.readyState === WebSocket.OPEN ||
        ws.readyState === WebSocket.CONNECTING
      ) {
        ws.close();
      }
      wsRef.current = null;
    };
  }, [docId, wsUrl]);

  // Create editor only when Y.Doc is ready
  const editor = useEditor(
    {
      extensions: [
        StarterKit,
        Collaboration.configure({
          document: ydocRef.current ?? new Y.Doc(),
        }),
      ],
      content:
        connectionStatus === "connected"
          ? "<p>Start typing...</p>"
          : "<p>Connecting...</p>",
      editable: connectionStatus === "connected",
    },
    [connectionStatus]
  );

  // Update editor editable state when connection status changes
  useEffect(() => {
    if (editor) {
      editor.setEditable(connectionStatus === "connected");
    }
  }, [editor, connectionStatus]);

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
