"use client";

import { useEffect, useMemo, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Phase17Service } from "@/services/phase17";
import { DiversityMetrics, providerPalette } from "@/types/phase17";

function pieGradient(metrics?: DiversityMetrics) {
  if (!metrics?.providers?.length) return "conic-gradient(#334155 0% 100%)";
  let cursor = 0;
  const chunks = metrics.providers.map((p) => {
    const start = cursor;
    cursor += p.usage_percentage;
    const color = p.provider === "OPENAI" ? "#60a5fa" : p.provider === "ANTHROPIC" ? "#fb923c" : p.provider === "GOOGLE" ? "#34d399" : "#a78bfa";
    return `${color} ${start}% ${cursor}%`;
  });
  return `conic-gradient(${chunks.join(",")})`;
}

export function DiversityInsightsCard({ projectId }: { projectId: string }) {
  const [metrics, setMetrics] = useState<DiversityMetrics | null>(null);

  useEffect(() => {
    Phase17Service.getDiversityMetrics(projectId).then(setMetrics).catch(console.error);
  }, [projectId]);

  const chart = useMemo(() => pieGradient(metrics || undefined), [metrics]);

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle className="text-lg">Diversidade de IA</CardTitle>
      </CardHeader>
      <CardContent>
        {!metrics ? (
          <p className="text-sm text-muted-foreground">Carregando métricas de diversidade...</p>
        ) : (
          <div className="grid md:grid-cols-[180px_1fr] gap-6 items-center">
            <div className="flex flex-col items-center gap-3">
              <div className="h-36 w-36 rounded-full border-8 border-background" style={{ background: chart }} />
              <p className="text-sm font-semibold">Score: {metrics.diversity_score}%</p>
            </div>

            <div className="space-y-3">
              {metrics.providers.map((provider) => {
                const palette = providerPalette[provider.provider];
                return (
                  <div key={provider.provider} className="flex items-center justify-between rounded-lg border p-2.5">
                    <span className={`px-2 py-1 rounded-md text-xs font-medium ${palette.bg} ${palette.text}`}>{provider.provider}</span>
                    <span className="text-sm text-muted-foreground">{provider.usage_percentage.toFixed(1)}% ({provider.count})</span>
                  </div>
                );
              })}
              <div className="grid grid-cols-2 gap-2 pt-1 text-sm">
                <div className="rounded-lg border p-2">
                  <p className="text-muted-foreground">Tríades 100% diversas</p>
                  <p className="font-semibold">{metrics.full_diversity_triads}</p>
                </div>
                <div className="rounded-lg border p-2">
                  <p className="text-muted-foreground">Com repetição</p>
                  <p className="font-semibold">{metrics.repeated_provider_triads}</p>
                </div>
              </div>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}
