"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { AgentService } from "@/services/agent";
import { Phase17Service } from "@/services/phase17";
import { Agent } from "@/types/agent";
import { TriadRole, roleLabel } from "@/types/phase17";
import { toast } from "sonner";
import { Dices } from "lucide-react";
import { DynamicModeHelpDialog } from "./DynamicModeHelpDialog";

const roles: TriadRole[] = ["PRODUCER", "REVIEWER", "REFINER"];

export function DynamicModeControl({ projectId, initialEnabled }: { projectId: string; initialEnabled: boolean }) {
  const [enabled, setEnabled] = useState(initialEnabled);
  const [saving, setSaving] = useState(false);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [fixed, setFixed] = useState<Partial<Record<TriadRole, string>>>({});
  const [preview, setPreview] = useState<string>("");

  useEffect(() => {
    AgentService.getAgents(1, 100)
      .then((data) => setAgents(data.items.filter((agent) => agent.enabled)))
      .catch(console.error);
  }, []);

  const runPreview = async () => {
    try {
      const data = await Phase17Service.previewDynamicSelection(projectId);
      setPreview(`${data.triad.producer.name} • ${data.triad.reviewer.name} • ${data.triad.refiner.name}`);
    } catch (error) {
      console.error(error);
      toast.error("Falha ao gerar preview do sorteio.");
    }
  };

  const save = async () => {
    setSaving(true);
    try {
      await Phase17Service.updateDynamicMode(projectId, {
        enabled,
        fixed_agents: enabled ? undefined : fixed,
      });
      toast.success("Configuração de modo salva com sucesso.");
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível salvar o Modo Dinâmico.");
    } finally {
      setSaving(false);
    }
  };

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-lg">Modo Dinâmico Multi-Modelo</CardTitle>
            <CardDescription>Ative o sorteio automático da tríade por execução com diversidade de providers.</CardDescription>
          </div>
          <DynamicModeHelpDialog />
        </div>
      </CardHeader>
      <CardContent className="space-y-5">
        <div className="flex items-center justify-between p-4 rounded-xl border bg-background/80">
          <div>
            <p className="font-medium">Status atual</p>
            <p className="text-sm text-muted-foreground">{enabled ? "Dinâmico (sorteio ativo)" : "Fixo (seleção manual por papel)"}</p>
          </div>
          <Switch checked={enabled} onCheckedChange={setEnabled} />
        </div>

        {!enabled && (
          <div className="grid md:grid-cols-3 gap-3">
            {roles.map((role) => (
              <div key={role} className="space-y-1.5">
                <p className="text-sm font-medium">{roleLabel[role]}</p>
                <Select value={fixed[role]} onValueChange={(value) => setFixed((prev) => ({ ...prev, [role]: value }))}>
                  <SelectTrigger>
                    <SelectValue placeholder="Selecione um agente" />
                  </SelectTrigger>
                  <SelectContent>
                    {agents.map((agent) => (
                      <SelectItem key={`${role}-${agent.id}`} value={agent.id}>{agent.name} · {agent.provider}</SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
            ))}
          </div>
        )}

        {enabled && (
          <div className="rounded-xl border border-primary/30 bg-primary/5 p-4">
            <div className="flex flex-wrap items-center gap-2">
              <Button variant="outline" size="sm" onClick={runPreview}>
                <Dices className="h-4 w-4 mr-1" /> Preview do sorteio
              </Button>
              {preview && <p className="text-sm text-primary">{preview}</p>}
            </div>
          </div>
        )}

        <Button onClick={save} disabled={saving}>{saving ? "Salvando..." : "Salvar configuração"}</Button>
      </CardContent>
    </Card>
  );
}
