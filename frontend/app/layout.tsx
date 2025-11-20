import type { Metadata } from "next";
import "../styles/globals.css";

export const metadata: Metadata = {
  title: "Soros Airbyte Clone",
  description: "Lightweight data movement control plane",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="app-body">
        <header className="app-header">
          <div className="container header-content">
            <div className="brand">
              <div className="logo-mark" />
              <div>
                <p className="muted">Soros</p>
                <p className="title">Airbyte-inspired control</p>
              </div>
            </div>
            <span className="pill success">Fast & minimal</span>
          </div>
        </header>
        <main className="container main-content">{children}</main>
      </body>
    </html>
  );
}
