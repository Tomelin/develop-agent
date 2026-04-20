"use client";

import { useState, useEffect, use } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/contexts/AuthContext";
import { agentService } from "@/services/agentService";
import { Agent } from "@/types/agent";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { ArrowLeft, Edit, Trash, Play, Bot, BrainCircuit, Activity, Clock } from "lucide-react";
import { toast } from "sonner";
import Link from "next/link";
import { AgentFormDrawer } from "@/components/agents/AgentFormDrawer";
import { Separator } from "@/components/ui/separator";
import { Skeleton } from "@/components/ui/skeleton";

export default function AgentDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const router = useRouter();
  const { user } = useAuth();
  const isAdmin = user?.role === "ADMIN";
  const resolvedParams = use(params);

  const [agent, setAgent] = useState<Agent | null>(null);
  const [loading, setLoading] = useState(true);
  const [drawerOpen, setDrawerOpen] = useState(false);

  const fetchAgentManual = async () => {
    try {
      setLoading(true);
      const data = await agentService.getAgentById(resolvedParams.id);
      setAgent(data);
    } catch (error) {
      console.error(error);
      toast.error("Erro ao carregar detalhes do agente");
      router.push("/dashboard/agents");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const fetchAgent = async () => {
      try {
        setLoading(true);
        const data = await agentService.getAgentById(resolvedParams.id);
        setAgent(data);
      } catch (error) {
        console.error(error);
        toast.error("Erro ao carregar detalhes do agente");
        router.push("/dashboard/agents");
      } finally {
        setLoading(false);
      }
    };

    fetchAgent();
  }, [resolvedParams.id, router]);

  const handleDelete = async () => {
    if (!agent) return;
    if (!confirm("Tem certeza que deseja remover este agente? Esta ação não pode ser desfeita.")) return;
    try {
      await agentService.deleteAgent(agent.id);
      toast.success("Agente removido com sucesso");
      router.push("/dashboard/agents");
    } catch (error) {
      console.error(error);
      toast.error("Erro ao remover agente");
    }
  };

  const handleTestConnection = async () => {
    if (!agent) return;
    try {
      toast.info("Testando conexão com o LLM...");
      const response = await agentService.testConnection(agent.id);
      if (response.success) {
        toast.success(`Conexão bem-sucedida!`);
      } else {
        toast.error(`Falha na conexão: ${response.message}`);
      }
    } catch (error) {
      console.error(error);
      toast.error("Erro ao testar conexão");
    }
  };

  if (loading) {
    return (
      <div className="space-y-6">
        <div className="flex items-center gap-4">
          <Skeleton className="h-10 w-10 rounded-full" />
          <Skeleton className="h-8 w-[300px]" />
        </div>
        <Skeleton className="h-[200px] w-full" />
        <Skeleton className="h-[400px] w-full" />
      </div>
    );
  }

  if (!agent) return null;

  return (
    <div className="space-y-6 animate-in fade-in">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div className="flex items-center gap-4">
          <Link href="/dashboard/agents">
            <Button variant="ghost" size="icon">
              <ArrowLeft className="h-5 w-5" />
            </Button>
          </Link>
          <div>
            <div className="flex items-center gap-3">
              <h1 className="text-3xl font-bold tracking-tight">{agent.name}</h1>
              {!agent.enabled && <Badge variant="destructive">Inativo</Badge>}
            </div>
            <p className="text-muted-foreground">{agent.description}</p>
          </div>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" onClick={handleTestConnection}>
            <Play className="mr-2 h-4 w-4" />
            Testar Conexão
          </Button>
          {isAdmin && (
            <>
              <Button variant="secondary" onClick={() => setDrawerOpen(true)}>
                <Edit className="mr-2 h-4 w-4" />
                Editar
              </Button>
              <Button variant="destructive" onClick={handleDelete}>
                <Trash className="mr-2 h-4 w-4" />
                Remover
              </Button>
            </>
          )}
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {/* Left Column - Info */}
        <div className="space-y-6 md:col-span-1">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center text-lg">
                <BrainCircuit className="mr-2 h-5 w-5 text-primary" />
                Configuração do Modelo
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <span className="text-sm text-muted-foreground block mb-1">Provider</span>
                <Badge variant="outline" className="text-sm">{agent.provider}</Badge>
              </div>
              <div>
                <span className="text-sm text-muted-foreground block mb-1">Modelo Específico</span>
                <code className="bg-muted px-2 py-1 rounded text-sm font-mono">{agent.model}</code>
              </div>
              <Separator />
              <div>
                <span className="text-sm text-muted-foreground block mb-2">Skills (Fases)</span>
                <div className="flex flex-wrap gap-2">
                  {agent.skills.map(skill => (
                    <Badge key={skill} variant="secondary">
                      {skill.replace("DEVELOPMENT_", "")}
                    </Badge>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="flex items-center text-lg">
                <Activity className="mr-2 h-5 w-5 text-primary" />
                Status Atual
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Estado Operacional</span>
                <Badge variant="outline">{agent.status || "IDLE"}</Badge>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Criado em</span>
                <span className="text-sm font-medium">
                  {agent.created_at ? new Date(agent.created_at).toLocaleDateString("pt-BR") : "N/A"}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="text-sm text-muted-foreground">Última atualização</span>
                <span className="text-sm font-medium">
                  {agent.updated_at ? new Date(agent.updated_at).toLocaleDateString("pt-BR") : "N/A"}
                </span>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Right Column - Prompts & History */}
        <div className="space-y-6 md:col-span-2">
          <Card className="h-full">
            <CardHeader>
              <CardTitle className="flex items-center">
                <Bot className="mr-2 h-5 w-5 text-primary" />
                System Prompts (Persona)
              </CardTitle>
              <CardDescription>
                As diretrizes base que moldam o comportamento e conhecimento técnico deste agente.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              {agent.system_prompts.map((prompt, idx) => (
                <div key={idx} className="bg-muted/50 p-4 rounded-lg border">
                  <div className="text-xs text-muted-foreground mb-2 font-mono uppercase tracking-wider">
                    Diretriz {idx + 1}
                  </div>
                  <p className="text-sm whitespace-pre-wrap leading-relaxed">
                    {prompt}
                  </p>
                </div>
              ))}
            </CardContent>
          </Card>
        </div>
      </div>

      <div className="mt-6">
         <Card>
            <CardHeader>
              <CardTitle className="flex items-center text-lg">
                <Clock className="mr-2 h-5 w-5 text-primary" />
                Histórico de Uso
              </CardTitle>
              <CardDescription>
                Últimas execuções deste agente nos projetos da agência.
              </CardDescription>
            </CardHeader>
            <CardContent>
                <div className="text-center py-8 text-muted-foreground border-2 border-dashed rounded-lg bg-muted/20">
                  Histórico de execuções será populado conforme projetos avancem na esteira.
                </div>
            </CardContent>
          </Card>
      </div>

      {drawerOpen && (
        <AgentFormDrawer
          open={drawerOpen}
          onOpenChange={setDrawerOpen}
          agent={agent}
          onSuccess={fetchAgentManual}
        />
      )}
    </div>
  );
}
