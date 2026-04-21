"use client";

import { useEffect, useState } from "react";
import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Phase20Service } from "@/services/phase20";
import { PublicRoadmapFeature, PublicRoadmapResponse } from "@/types/phase20";
import { toast } from "sonner";

const statuses: PublicRoadmapFeature["status"][] = ["PLANNED", "IN_DEVELOPMENT", "COMPLETED"];

export default function AdminRoadmapPage() {
  const [roadmap, setRoadmap] = useState<PublicRoadmapResponse | null>(null);

  const loadRoadmap = async () => {
    try {
      const data = await Phase20Service.getPublicRoadmap();
      setRoadmap(data);
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar o roadmap para administração.");
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      void loadRoadmap();
    }, 0);
    return () => clearTimeout(timer);
  }, []);

  const updateStatus = async (featureId: string, status: string) => {
    try {
      await Phase20Service.updateRoadmapFeatureStatus(featureId, status);
      setRoadmap((current) => current ? {
        ...current,
        features: current.features.map((feature) => feature.id === featureId ? { ...feature, status: status as PublicRoadmapFeature["status"] } : feature),
      } : current);
      toast.success("Status do item atualizado.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao atualizar status.");
    }
  };

  return (
    <PrivateRoute>
      <div className="space-y-6">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Admin · Roadmap</h1>
          <p className="text-muted-foreground">Gerencie status e priorização dos itens públicos do roadmap.</p>
        </div>

        <Card className="bg-card/50">
          <CardHeader>
            <CardTitle>Features do roadmap</CardTitle>
            <CardDescription>Atualize os cards exibidos em /roadmap com total rastreabilidade.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-2">
            {roadmap?.features.map((feature) => (
              <div key={feature.id} className="rounded-lg border p-3 flex flex-col md:flex-row md:items-center md:justify-between gap-3">
                <div>
                  <p className="font-medium">{feature.title}</p>
                  <p className="text-sm text-muted-foreground">{feature.description}</p>
                  <div className="flex items-center gap-2 mt-2">
                    <Badge variant="secondary">{feature.milestone}</Badge>
                    <Badge variant="outline">{feature.votes} votos</Badge>
                  </div>
                </div>
                <Select value={feature.status} onValueChange={(value) => updateStatus(feature.id, value)}>
                  <SelectTrigger className="w-[220px]"><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {statuses.map((status) => <SelectItem key={status} value={status}>{status}</SelectItem>)}
                  </SelectContent>
                </Select>
              </div>
            ))}
          </CardContent>
        </Card>
      </div>
    </PrivateRoute>
  );
}
