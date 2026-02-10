import React from "react";
import ReactDOM from "react-dom/client";
import "./styles.css";

type HealthResponse = {
  status: string;
  timestamp: string;
};

function App() {
  const [health, setHealth] = React.useState<HealthResponse | null>(null);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    fetch("/api/health")
      .then(async (response) => {
        if (!response.ok) {
          throw new Error(`request failed with status ${response.status}`);
        }
        return (await response.json()) as HealthResponse;
      })
      .then((data) => {
        setHealth(data);
      })
      .catch((err: unknown) => {
        const message = err instanceof Error ? err.message : "unknown error";
        setError(message);
      });
  }, []);

  return (
    <main className="app">
      <h1>PulseLounge</h1>
      <p>Go API + React TypeScript UI</p>

      {error && <p className="error">Health check failed: {error}</p>}

      {health && (
        <section className="card">
          <h2>Health Check</h2>
          <p>
            <strong>Status:</strong> {health.status}
          </p>
          <p>
            <strong>Timestamp (UTC):</strong> {health.timestamp}
          </p>
        </section>
      )}
    </main>
  );
}

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
);
