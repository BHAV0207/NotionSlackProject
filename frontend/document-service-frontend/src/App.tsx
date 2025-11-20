import CollaborativeEditor from "./components/CollaborativeEditor";

export default function App() {
  // you can change docId to test multiple documents
  return (
    <div style={{ padding: 24 }}>
      <h2>Collaborative Editor (React + TypeScript + TipTap + Yjs)</h2>
      <p>Open two tabs to test real-time sync.</p>
      <CollaborativeEditor docId="my-doc-1" />
    </div>
  );
}
