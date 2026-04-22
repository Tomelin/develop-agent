"use client";

import { useMemo, useState } from "react";
import { PhaseArtifact, PhaseTrack, PhaseTrackStatus, TrackFeedbackItem } from "@/types/phase8";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Send, CheckCheck } from "lucide-react";

interface TrackFeedbackPanelProps {
  trackStatuses: PhaseTrackStatus[];
  artifacts: PhaseArtifact[];
  feedbackHistory: Record<PhaseTrack, TrackFeedbackItem[]>;
  onSendFeedback: (track: PhaseTrack, content: string) => Promise<void>;
  onApproveTrack: (track: PhaseTrack) => Promise<void>;
}

const tracks: PhaseTrack[] = ["FRONTEND", "BACKEND"];

export function TrackFeedbackPanel({ trackStatuses, artifacts, feedbackHistory, onSendFeedback, onApproveTrack }: TrackFeedbackPanelProps) {
  const [feedbackByTrack, setFeedbackByTrack] = useState<Record<PhaseTrack, string>>({ FRONTEND: "", BACKEND: "" });
  const [loadingTrack, setLoadingTrack] = useState<PhaseTrack | null>(null);

  const statusByTrack = useMemo(() => {
    const statuses = Array.isArray(trackStatuses) ? trackStatuses : [];
    return statuses.reduce((acc, item) => ({ ...acc, [item.track]: item }), {} as Record<PhaseTrack, PhaseTrackStatus>);
  }, [trackStatuses]);

  const handleSend = async (track: PhaseTrack) => {
    const content = feedbackByTrack[track].trim();
    if (!content) return;
    setLoadingTrack(track);
    await onSendFeedback(track, content);
    setFeedbackByTrack((prev) => ({ ...prev, [track]: "" }));
    setLoadingTrack(null);
  };

  const handleApprove = async (track: PhaseTrack) => {
    setLoadingTrack(track);
    await onApproveTrack(track);
    setLoadingTrack(null);
  };

  return (
    <Card className="border-border/60 bg-card/60">
      <CardHeader>
        <CardTitle className="text-base">Feedback por Trilho</CardTitle>
      </CardHeader>
      <CardContent>
        <Tabs defaultValue="FRONTEND" className="w-full">
          <TabsList>
            <TabsTrigger value="FRONTEND">Feedback Frontend</TabsTrigger>
            <TabsTrigger value="BACKEND">Feedback Backend</TabsTrigger>
          </TabsList>

          {tracks.map((track) => {
            const status = statusByTrack[track];
            const latestArtifact = artifacts.find((item) => item.track === track);
            return (
              <TabsContent value={track} key={track} className="mt-4 space-y-4">
                <div className="flex items-center justify-between rounded-lg border border-border/60 bg-background/40 p-3">
                  <span className="text-sm text-muted-foreground">Feedbacks utilizados</span>
                  <Badge variant="secondary">{status?.feedbacks_used ?? 0} de {status?.feedbacks_limit ?? 5} feedbacks utilizados</Badge>
                </div>

                <div className="space-y-2">
                  <Textarea
                    rows={7}
                    placeholder={`Descreva melhorias para o trilho ${track}...`}
                    value={feedbackByTrack[track]}
                    onChange={(event) => setFeedbackByTrack((prev) => ({ ...prev, [track]: event.target.value }))}
                  />
                  <p className="text-xs text-muted-foreground">Suporte a markdown básico: listas, ênfase e blocos de código.</p>
                </div>

                <div className="rounded-lg border border-border/60 bg-background/40 p-3">
                  <p className="mb-2 text-xs font-semibold uppercase text-muted-foreground">Preview para o Refinador</p>
                  <p className="text-xs text-muted-foreground">Artefato atual: <span className="font-medium text-foreground">{latestArtifact?.title ?? "Sem artefato"}</span></p>
                  <pre className="mt-2 max-h-28 overflow-auto rounded bg-background p-2 text-xs text-muted-foreground whitespace-pre-wrap">{feedbackByTrack[track] || "Seu feedback aparecerá aqui em tempo real."}</pre>
                </div>

                <div className="flex flex-wrap gap-2">
                  <Button disabled={loadingTrack === track} onClick={() => handleSend(track)}><Send className="mr-2 h-4 w-4" />Enviar Feedback e Refinar</Button>
                  <Button disabled={loadingTrack === track} variant="outline" onClick={() => handleApprove(track)}><CheckCheck className="mr-2 h-4 w-4" />Aprovar e Avançar</Button>
                </div>
                <details className="rounded-lg border border-border/60 bg-background/30 p-3">
                  <summary className="cursor-pointer text-sm font-medium">Histórico de feedbacks ({feedbackHistory[track]?.length ?? 0})</summary>
                  <div className="mt-3 space-y-2">
                    {(feedbackHistory[track] ?? []).map((item) => (
                      <div key={item.id} className="rounded border border-border/60 bg-background/40 p-2 text-xs">
                        <p className="mb-1 text-muted-foreground">{new Date(item.created_at).toLocaleString()}</p>
                        <p className="whitespace-pre-wrap">{item.content}</p>
                      </div>
                    ))}
                    {!(feedbackHistory[track] ?? []).length && (
                      <p className="text-xs text-muted-foreground">Nenhum feedback enviado ainda.</p>
                    )}
                  </div>
                </details>
              </TabsContent>
            );
          })}
        </Tabs>
      </CardContent>
    </Card>
  );
}
