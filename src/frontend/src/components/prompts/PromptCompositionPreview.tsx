"use client";

import { useEffect, useMemo, useState } from "react";
import { Copy, Eye } from "lucide-react";
import { toast } from "sonner";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { PROMPT_GROUPS, PromptGroup, PromptPreviewBlock } from "@/types/prompt";
import { promptService } from "@/services/promptService";
import { agentService } from "@/services/agentService";
import { Agent } from "@/types/agent";

const sourceStyle: Record<PromptPreviewBlock["source"], string> = {
  SYSTEM: "border-blue-500/40 bg-blue-500/10",
  GLOBAL: "border-emerald-500/40 bg-emerald-500/10",
  GROUP: "border-amber-500/40 bg-amber-500/10",
  RAG: "border-purple-500/40 bg-purple-500/10",
  PHASE_INSTRUCTION: "border-slate-500/40 bg-slate-500/20",
};

export function PromptCompositionPreview() {
  const [group, setGroup] = useState<PromptGroup>("GLOBAL");
  const [agents, setAgents] = useState<Agent[]>([]);
  const [agentId, setAgentId] = useState<string>("none");
  const [loading, setLoading] = useState(false);

  const [blocks, setBlocks] = useState<PromptPreviewBlock[]>([]);
  const [composedPrompt, setComposedPrompt] = useState("");
  const [tokenEstimate, setTokenEstimate] = useState<number | undefined>(undefined);

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const data = await agentService.getAgents({ page: 1, size: 100, enabled: true });
        setAgents(data.items || []);
      } catch {
        toast.error("Não foi possível carregar os agentes.");
      }
    };
    fetchAgents();
  }, []);

  const selectedAgent = useMemo(() => agents.find((agent) => agent.id === agentId), [agents, agentId]);

  const loadPreview = async (selectedGroup: PromptGroup, selectedAgentId?: string) => {
    setLoading(true);
    try {
      const response = await promptService.getPreview(selectedGroup, selectedAgentId);
      setBlocks(response.blocks || []);
      setComposedPrompt(response.composed_prompt || "");
      setTokenEstimate(response.token_estimate);
    } catch {
      toast.error("Falha ao gerar preview da composição.");
      setBlocks([]);
      setComposedPrompt("");
      setTokenEstimate(undefined);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    loadPreview(group, agentId === "none" ? undefined : agentId);
  }, [group, agentId]);

  const copyComposed = async () => {
    if (!composedPrompt) return;
    await navigator.clipboard.writeText(composedPrompt);
    toast.success("Prompt composto copiado.");
  };

  return (
    <div className="space-y-6">
      <Card className="bg-card/70 border-border/70">
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Eye className="h-5 w-5 text-primary" /> Preview de Composição
          </CardTitle>
          <CardDescription>
            Visualize exatamente como o motor de agregação injetará os blocos de instrução antes de executar a phase.
          </CardDescription>
        </CardHeader>
        <CardContent className="grid md:grid-cols-3 gap-4">
          <div className="space-y-2">
            <p className="text-sm font-medium">Grupo/Fase</p>
            <Select value={group} onValueChange={(value) => setGroup(value as PromptGroup)}>
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {PROMPT_GROUPS.map((item) => (
                  <SelectItem value={item.value} key={item.value}>
                    {item.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <p className="text-sm font-medium">Agente</p>
            <Select value={agentId} onValueChange={setAgentId}>
              <SelectTrigger>
                <SelectValue placeholder="Selecione um agente" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="none">Sem agente específico</SelectItem>
                {agents.map((agent) => (
                  <SelectItem key={agent.id} value={agent.id}>
                    {agent.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="flex items-end justify-start md:justify-end">
            <Button className="gap-2" onClick={copyComposed} disabled={!composedPrompt}>
              <Copy className="h-4 w-4" /> Copiar Prompt Completo
            </Button>
          </div>
        </CardContent>
      </Card>

      {selectedAgent ? (
        <Card className="bg-card/40">
          <CardHeader>
            <CardTitle className="text-base">Agente selecionado</CardTitle>
            <CardDescription>
              {selectedAgent.name} · {selectedAgent.provider} · {selectedAgent.model}
            </CardDescription>
          </CardHeader>
        </Card>
      ) : null}

      {tokenEstimate !== undefined ? (
        <Badge variant="outline" className="text-sm">
          Estimativa de tokens: {tokenEstimate.toLocaleString()}
        </Badge>
      ) : null}

      <div className="space-y-3">
        {loading ? (
          <div className="h-36 animate-pulse rounded-xl bg-muted/40" />
        ) : (
          blocks.map((block, index) => (
            <div key={`${block.source}-${index}`} className={`rounded-xl border p-4 ${sourceStyle[block.source]}`}>
              <div className="mb-2 flex items-center justify-between gap-2">
                <h3 className="font-semibold">{block.title}</h3>
                <Badge>{block.source}</Badge>
              </div>
              <p className="whitespace-pre-wrap text-sm text-muted-foreground">{block.content}</p>
            </div>
          ))
        )}
      </div>

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Prompt completo composto</CardTitle>
        </CardHeader>
        <CardContent>
          <pre className="max-h-[320px] overflow-auto rounded-lg border bg-muted/20 p-4 text-xs leading-relaxed whitespace-pre-wrap">
            {composedPrompt || "Sem composição disponível para o filtro atual."}
          </pre>
        </CardContent>
      </Card>
    </div>
  );
}
