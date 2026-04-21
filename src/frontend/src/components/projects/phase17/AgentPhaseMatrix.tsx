"use client";

import { useEffect, useMemo, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import { AgentService } from "@/services/agent";
import { Phase17Service } from "@/services/phase17";
import { Agent } from "@/types/agent";
import { phaseMatrixDefaults, PhaseAgentMatrixRow, roleKeyMap, roleLabel, TRIAD_ROLES } from "@/types/phase17";
import { toast } from "sonner";

const UNASSIGNED_AGENT = "__UNASSIGNED_AGENT__";

export function AgentPhaseMatrix({ projectId }: { projectId: string }) {
  const [rows, setRows] = useState<PhaseAgentMatrixRow[]>(phaseMatrixDefaults);
  const [agents, setAgents] = useState<Agent[]>([]);
  const [costPreview, setCostPreview] = useState<string>("-");

  useEffect(() => {
    const load = async () => {
      try {
        const [matrix, allAgents] = await Promise.all([
          Phase17Service.getAgentMatrix(projectId),
          AgentService.getAgents(1, 200),
        ]);
        setRows(matrix.rows?.length ? matrix.rows : phaseMatrixDefaults);
        setAgents(allAgents.items.filter((a) => a.enabled));
      } catch (error) {
        console.error(error);
      }
    };
    load();
  }, [projectId]);

  const agentsMap = useMemo(() => new Map(agents.map((a) => [a.id, a])), [agents]);

  const applyToAll = (role: "PRODUCER" | "REVIEWER" | "REFINER", agentId: string) => {
    const key = roleKeyMap[role];
    setRows((prev) => prev.map((row) => ({ ...row, [key]: agentId, dynamic: false })));
  };

  const save = async () => {
    try {
      await Phase17Service.updateAgentMatrix(projectId, { rows });
      toast.success("Matriz de agentes atualizada.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao salvar matriz de agentes.");
    }
  };

  const previewCost = async () => {
    try {
      const data = await Phase17Service.previewConfigurationCost(projectId, { rows });
      setCostPreview(`US$ ${data.monthly_estimated_usd.toFixed(2)} · ${data.note}`);
    } catch (error) {
      console.error(error);
      toast.error("Falha ao calcular custo estimado.");
    }
  };

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle>Configuração de Agentes por Fase</CardTitle>
        <CardDescription>Matriz fixa/dinâmica por fase com opção de aplicar em massa por papel.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid lg:grid-cols-3 gap-2">
          {TRIAD_ROLES.map((role) => (
            <div key={role} className="flex items-center gap-2">
              <Select onValueChange={(value) => applyToAll(role, String(value))}>
                <SelectTrigger>
                  <SelectValue placeholder={`Aplicar ${roleLabel[role]} em todas`} />
                </SelectTrigger>
                <SelectContent>
                  {agents.map((agent) => <SelectItem key={`${role}-${agent.id}`} value={agent.id}>{agent.name}</SelectItem>)}
                </SelectContent>
              </Select>
            </div>
          ))}
        </div>

        <div className="overflow-x-auto rounded-xl border">
          <table className="w-full text-sm">
            <thead className="bg-muted/40">
              <tr>
                <th className="text-left p-3 min-w-[150px]">Fase</th>
                <th className="text-left p-3">Produtor</th>
                <th className="text-left p-3">Revisor</th>
                <th className="text-left p-3">Refinador</th>
                <th className="text-left p-3">Dinâmico</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((row) => (
                <tr key={row.phase_key} className="border-t">
                  <td className="p-3 font-medium">{row.phase_label}</td>
                  {TRIAD_ROLES.map((role) => {
                    const key = roleKeyMap[role];
                    const selectedAgentId = row[key] ?? null;
                    const selectValue = selectedAgentId ?? UNASSIGNED_AGENT;
                    return (
                      <td className="p-2" key={`${row.phase_key}-${role}`}>
                        <Select
                          value={selectValue}
                          onValueChange={(selected) => {
                            const selectedId = String(selected);
                            const nextAgentId = selectedId === UNASSIGNED_AGENT ? null : selectedId;
                            setRows((prev) => prev.map((item) => item.phase_key === row.phase_key ? { ...item, [key]: nextAgentId, dynamic: false } : item));
                          }}
                        >
                          <SelectTrigger>
                            <SelectValue placeholder="Selecionar" />
                          </SelectTrigger>
                          <SelectContent>
                            <SelectItem value={UNASSIGNED_AGENT}>Não atribuído</SelectItem>
                            {agents.map((agent) => (
                              <SelectItem key={`${row.phase_key}-${role}-${agent.id}`} value={agent.id}>
                                {agent.name} · {agent.provider}
                              </SelectItem>
                            ))}
                          </SelectContent>
                        </Select>
                        {selectedAgentId && <Badge variant="secondary" className="mt-1 text-[10px]">{agentsMap.get(selectedAgentId)?.provider}</Badge>}
                      </td>
                    );
                  })}
                  <td className="p-3">
                    <Switch
                      checked={Boolean(row.dynamic)}
                      onCheckedChange={(checked) => setRows((prev) => prev.map((item) => item.phase_key === row.phase_key ? { ...item, dynamic: checked } : item))}
                    />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className="flex flex-wrap items-center gap-2">
          <Button onClick={previewCost} variant="outline">Preview de custo</Button>
          <Badge variant="outline">{costPreview}</Badge>
          <Button onClick={save}>Salvar matriz</Button>
        </div>
      </CardContent>
    </Card>
  );
}
