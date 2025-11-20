import { InfoCard } from "@/components/info-card";

const cards = [
  {
    title: "Sources",
    metric: "12 connected",
    description: "Effortless ingestion across Postgres, Mongo, Stripe, and more.",
    accent: "indigo" as const,
  },
  {
    title: "Destinations",
    metric: "7 active",
    description: "Warehouse-first delivery with Snowflake, BigQuery, and S3.",
    accent: "cyan" as const,
  },
  {
    title: "Sync health",
    metric: "99.95% uptime",
    description: "Lean Go services keep pipelines responsive and debuggable.",
    accent: "emerald" as const,
  },
];

const featureList = [
  "Composable pipelines with minimal latency",
  "Realtime status panels built with Shadcn-inspired cards",
  "Predictable Go backend with typed responses",
  "Ready for containerized deployment on minimal hardware",
];

export default function Page() {
  return (
    <section className="page">
      <div className="panel">
        <p className="muted small">Pipeline overview</p>
        <h1>Lightweight Airbyte-style control plane</h1>
        <p className="subtle">
          Soros pairs a Next.js experience layer with a compact Go API. The interface favors clarity and speed, while the backend
          exposes predictable endpoints for sources, destinations, and connections.
        </p>
        <div className="card-grid">
          {cards.map((card) => (
            <InfoCard key={card.title} {...card} />
          ))}
        </div>
      </div>

      <div className="panel">
        <div className="section-heading">
          <div>
            <p className="muted small">Performance profile</p>
            <h2>Built to stay lean</h2>
          </div>
          <span className="pill accent">Shadcn-inspired</span>
        </div>
        <ul className="feature-list">
          {featureList.map((feature) => (
            <li key={feature} className="feature-item">
              <span className="dot" />
              <span>{feature}</span>
            </li>
          ))}
        </ul>
      </div>
    </section>
  );
}
