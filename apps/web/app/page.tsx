const stats = [
  ["Requests", "0"],
  ["Error rate", "0%"],
  ["Avg latency", "—"],
];

export default function Home() {
  return (
    <main style={{ maxWidth: 920, margin: "0 auto", padding: "48px 20px", fontFamily: "system-ui" }}>
      <header style={{ marginBottom: 40 }}>
        <p style={{ opacity: 0.6, marginBottom: 8 }}>Developer observability</p>
        <h1 style={{ fontSize: 42, margin: 0 }}>TraceDrop</h1>
      </header>

      <section style={{ display: "grid", gridTemplateColumns: "repeat(3, 1fr)", gap: 16, marginBottom: 32 }}>
        {stats.map(([label, value]) => (
          <article key={label} style={{ border: "1px solid #ddd", borderRadius: 12, padding: 20 }}>
            <div style={{ opacity: 0.65 }}>{label}</div>
            <strong style={{ fontSize: 28 }}>{value}</strong>
          </article>
        ))}
      </section>

      <section style={{ border: "1px solid #ddd", borderRadius: 12, padding: 24 }}>
        <h2 style={{ marginTop: 0 }}>Recent traces</h2>
        <p style={{ opacity: 0.65 }}>No traces yet. The ingestion API arrives in the next phase.</p>
      </section>
    </main>
  );
}
