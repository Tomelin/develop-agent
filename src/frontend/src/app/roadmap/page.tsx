"use client";

import { useEffect, useMemo, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { PublicRoadmapResponse } from "@/types/phase20";
import { Phase20Service } from "@/services/phase20";
import { toast } from "sonner";
import { Rocket } from "lucide-react";

const featureStatusLabel: Record<string, string> = {
  PLANNED: "Planejado",
  IN_DEVELOPMENT: "Em desenvolvimento",
  COMPLETED: "Concluído",
};

export default function PublicRoadmapPage() {
  const [roadmap, setRoadmap] = useState<PublicRoadmapResponse | null>(null);
  const [suggestionTitle, setSuggestionTitle] = useState("");
  const [suggestionDescription, setSuggestionDescription] = useState("");

  const loadRoadmap = async () => {
    try {
      const data = await Phase20Service.getPublicRoadmap();
      setRoadmap(data);
    } catch (error) {
      console.error(error);
      toast.error("Falha ao carregar roadmap público.");
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      void loadRoadmap();
    }, 0);
    return () => clearTimeout(timer);
  }, []);

  const milestones = useMemo(() => Array.from(new Set(roadmap?.features.map((item) => item.milestone) ?? [])), [roadmap]);

  const voteFeature = async (id: string) => {
    try {
      await Phase20Service.voteRoadmapFeature(id);
      toast.success("Voto registrado com sucesso.");
      await loadRoadmap();
    } catch (error) {
      console.error(error);
      toast.error("Você precisa estar logado para votar.");
    }
  };

  const suggestFeature = async () => {
    if (!suggestionTitle || !suggestionDescription) return;

    try {
      await Phase20Service.suggestRoadmapFeature({ title: suggestionTitle, description: suggestionDescription });
      toast.success("Sugestão enviada para análise do time de produto.");
      setSuggestionTitle("");
      setSuggestionDescription("");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao enviar sugestão.");
    }
  };

  return (
    <div className="container py-10 max-w-7xl space-y-6">
      <div className="space-y-2">
        <h1 className="text-4xl font-bold tracking-tight">Roadmap Público de Produto</h1>
        <p className="text-muted-foreground">Visão de 12-18 meses do produto, com milestones, votação e changelog oficial.</p>
      </div>

      <Card className="bg-card/50">
        <CardHeader>
          <CardTitle>Visão estratégica</CardTitle>
          <CardDescription>{roadmap?.vision || "Carregando visão de produto..."}</CardDescription>
        </CardHeader>
      </Card>

      <Tabs defaultValue={milestones[0] || "all"} className="space-y-4">
        <TabsList>
          {milestones.map((milestone) => <TabsTrigger key={milestone} value={milestone}>{milestone}</TabsTrigger>)}
        </TabsList>

        {milestones.map((milestone) => (
          <TabsContent key={milestone} value={milestone} className="grid gap-4 md:grid-cols-2">
            {roadmap?.features.filter((feature) => feature.milestone === milestone).map((feature) => (
              <Card key={feature.id}>
                <CardHeader>
                  <div className="flex items-center justify-between gap-2">
                    <CardTitle className="text-lg">{feature.title}</CardTitle>
                    <Badge>{featureStatusLabel[feature.status] || feature.status}</Badge>
                  </div>
                  <CardDescription>{feature.description}</CardDescription>
                </CardHeader>
                <CardContent className="flex items-center justify-between">
                  <p className="text-sm text-muted-foreground">{feature.votes} votos • Meta {feature.target_quarter || "TBD"}</p>
                  <Button size="sm" onClick={() => voteFeature(feature.id)}>Votar</Button>
                </CardContent>
              </Card>
            ))}
          </TabsContent>
        ))}
      </Tabs>

      <div className="grid gap-4 lg:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle>Changelog de versões</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {roadmap?.changelog.map((item) => (
              <div key={item.version} className="rounded-lg border p-3">
                <p className="font-medium flex items-center gap-2"><Rocket className="h-4 w-4 text-primary" /> {item.version}</p>
                <p className="text-xs text-muted-foreground">{new Date(item.date).toLocaleDateString()}</p>
                <ul className="mt-2 space-y-1 list-disc pl-5 text-sm">
                  {item.highlights.map((highlight) => <li key={highlight}>{highlight}</li>)}
                </ul>
              </div>
            ))}
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Sugerir nova feature</CardTitle>
            <CardDescription>Compartilhe demandas da sua operação para influenciar o backlog.</CardDescription>
          </CardHeader>
          <CardContent className="space-y-3">
            <Input placeholder="Título da sugestão" value={suggestionTitle} onChange={(e) => setSuggestionTitle(e.target.value)} />
            <Textarea placeholder="Explique o problema, impacto e resultado esperado..." value={suggestionDescription} onChange={(e) => setSuggestionDescription(e.target.value)} />
            <Button onClick={suggestFeature}>Enviar sugestão</Button>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
