interface InfoCardProps {
  title: string;
  metric: string;
  description: string;
  accent: "indigo" | "cyan" | "emerald";
}

const accentClass: Record<InfoCardProps["accent"], string> = {
  indigo: "card indigo",
  cyan: "card cyan",
  emerald: "card emerald",
};

export function InfoCard({ title, metric, description, accent }: InfoCardProps) {
  return (
    <article className={accentClass[accent]}>
      <p className="muted small">{title}</p>
      <div className="metric">{metric}</div>
      <p className="subtle">{description}</p>
    </article>
  );
}
