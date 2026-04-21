"use client";

import { useCallback, useEffect, useState } from "react";
import { ArtifactViewer } from "./ArtifactViewer";
import { TriadProgressPanel } from "./TriadProgressPanel";
import { TrackFeedbackPanel } from "./TrackFeedbackPanel";
import { Phase8Service } from "@/services/phase8";
import { PhaseArtifact, PhaseTrack, PhaseTrackStatus, TrackFeedbackItem, TriadTrackRuntime } from "@/types/phase8";
import { Card, CardContent } from "@/components/ui/card";
import { toast } from "sonner";

interface Phase8WorkspaceProps {
  projectId: string;
  phaseNumber?: number;
}

export function Phase8Workspace({ projectId, phaseNumber = 2 }: Phase8WorkspaceProps) {
  const [loading, setLoading] = useState(true);
  const [trackStatuses, setTrackStatuses] = useState<PhaseTrackStatus[]>([]);
  const [artifacts, setArtifacts] = useState<PhaseArtifact[]>([]);
  const [triad, setTriad] = useState<TriadTrackRuntime[]>([]);
  const [feedbackHistory, setFeedbackHistory] = useState<Record<PhaseTrack, TrackFeedbackItem[]>>({ FRONTEND: [], BACKEND: [] });

  const fetchData = useCallback(async () => {
    try {
      const [statusData, artifactData, triadData, frontHistory, backHistory] = await Promise.all([
        Phase8Service.getTrackStatus(projectId, phaseNumber),
        Phase8Service.getArtifacts(projectId, phaseNumber),
        Phase8Service.getTriadProgress(projectId, phaseNumber),
        Phase8Service.getFeedbackHistory(projectId, phaseNumber, "FRONTEND"),
        Phase8Service.getFeedbackHistory(projectId, phaseNumber, "BACKEND"),
      ]);

      setTrackStatuses(statusData);
      setArtifacts(artifactData);
      setTriad(triadData);
      setFeedbackHistory({ FRONTEND: frontHistory, BACKEND: backHistory });
    } catch (error) {
      console.error(error);
    } finally {
      setLoading(false);
    }
  }, [projectId, phaseNumber]);

  useEffect(() => {
    const bootstrap = setTimeout(() => {
      void fetchData();
    }, 0);

    const interval = setInterval(() => {
      void fetchData();
    }, 10000);

    return () => {
      clearTimeout(bootstrap);
      clearInterval(interval);
    };
  }, [fetchData]);

  const handleSendFeedback = async (track: PhaseTrack, content: string) => {
    try {
      await Phase8Service.sendFeedback(projectId, phaseNumber, track, content);
      toast.success(`Feedback enviado para ${track}`);
      await fetchData();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao enviar feedback");
    }
  };

  const handleApproveTrack = async (track: PhaseTrack) => {
    try {
      await Phase8Service.approveTrack(projectId, phaseNumber, track);
      toast.success(`Trilho ${track} aprovado`);
      await fetchData();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao aprovar trilho");
    }
  };

  if (loading) {
    return (
      <Card className="border-border/60 bg-card/60">
        <CardContent className="py-12 text-center text-muted-foreground">Carregando workspace da Fase 08...</CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <TriadProgressPanel tracks={triad} />

      <TrackFeedbackPanel
        trackStatuses={trackStatuses}
        artifacts={artifacts}
        feedbackHistory={feedbackHistory}
        onSendFeedback={handleSendFeedback}
        onApproveTrack={handleApproveTrack}
      />

      <div className="space-y-4">
        {artifacts.map((artifact) => (
          <ArtifactViewer key={artifact.id} artifact={artifact} />
        ))}
        {!artifacts.length && (
          <Card className="border-dashed border-border/60 bg-card/40">
            <CardContent className="py-10 text-center text-muted-foreground">Nenhum artefato disponível para esta fase ainda.</CardContent>
          </Card>
        )}
      </div>
    </div>
  );
}
