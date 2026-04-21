"use client";

import { TriadTrackRuntime } from "@/types/phase8";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { LoaderCircle, Cpu, Timer, Binary } from "lucide-react";

const providerColor: Record<string, string> = {
  openai: "bg-emerald-500/20 text-emerald-300",
  anthropic: "bg-orange-500/20 text-orange-300",
  gemini: "bg-cyan-500/20 text-cyan-300",
  ollama: "bg-violet-500/20 text-violet-300",
};

const statusLabel: Record<string, string> = {
  PENDING: "Pendente",
  RUNNING: "Processando",
  REVIEW: "Revisão",
  COMPLETED: "Concluído",
  ERROR: "Erro",
};

export function TriadProgressPanel({ tracks }: { tracks: TriadTrackRuntime[] }) {
  return (
    <Card className="border-border/60 bg-card/60">
      <CardHeader>
        <CardTitle className="text-base">Acompanhamento da Tríade em Tempo Real</CardTitle>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue={tracks[0]?.track ?? "FRONTEND"} className="w-full">
          <TabsList>
            {tracks.map((track) => (
              <TabsTrigger key={track.track} value={track.track}>{track.track}</TabsTrigger>
            ))}
          </TabsList>
          {tracks.map((track) => (
            <TabsContent key={track.track} value={track.track} className="mt-4 grid gap-3 md:grid-cols-3">
              {track.steps.map((step) => {
                const isActive = step.status === "RUNNING";
                const initials = step.agent_name.split(" ").map((chunk) => chunk[0]).join("").slice(0, 2).toUpperCase();
                return (
                  <Card key={`${track.track}-${step.step}`} className={`border ${isActive ? "border-primary/60" : "border-border/70"}`}>
                    <CardContent className="space-y-3 p-4">
                      <div className="flex items-start justify-between gap-3">
                        <div className="flex items-center gap-3">
                          <Avatar className="h-10 w-10">
                            <AvatarFallback className={providerColor[step.provider?.toLowerCase?.()] ?? "bg-muted text-foreground"}>{initials}</AvatarFallback>
                          </Avatar>
                          <div>
                            <p className="font-medium leading-tight">{step.agent_name}</p>
                            <p className="text-xs text-muted-foreground">{step.step}</p>
                          </div>
                        </div>
                        <Badge variant={isActive ? "default" : "outline"} className={isActive ? "animate-pulse" : ""}>
                          {statusLabel[step.status] ?? step.status}
                        </Badge>
                      </div>

                      {isActive && (
                        <div className="rounded-md bg-primary/10 p-2 text-xs text-primary">
                          <div className="mb-1 flex items-center gap-2"><LoaderCircle className="h-3.5 w-3.5 animate-spin" /> Processando output em streaming</div>
                          <p className="line-clamp-6 whitespace-pre-wrap text-muted-foreground">{(step.partial_output ?? "Sem preview disponível").slice(0, 500)}</p>
                        </div>
                      )}

                      <div className="grid grid-cols-3 gap-2 text-[11px]">
                        <div className="rounded bg-background/70 p-2">
                          <p className="flex items-center gap-1 text-muted-foreground"><Binary className="h-3 w-3" />Tokens</p>
                          <p className="font-medium">{step.tokens_used?.toLocaleString?.() ?? "-"}</p>
                        </div>
                        <div className="rounded bg-background/70 p-2">
                          <p className="flex items-center gap-1 text-muted-foreground"><Timer className="h-3 w-3" />Duração</p>
                          <p className="font-medium">{step.duration_ms ? `${(step.duration_ms / 1000).toFixed(1)}s` : "-"}</p>
                        </div>
                        <div className="rounded bg-background/70 p-2">
                          <p className="flex items-center gap-1 text-muted-foreground"><Cpu className="h-3 w-3" />Modelo</p>
                          <p className="line-clamp-1 font-medium">{step.model ?? "-"}</p>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                );
              })}
            </TabsContent>
          ))}
        </Tabs>
      </CardContent>
    </Card>
  );
}
