"use client";

import { useState, useEffect } from "react";
import { useAuth } from "@/contexts/AuthContext";
import { agentService } from "@/services/agentService";
import { Agent } from "@/types/agent";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Switch } from "@/components/ui/switch";
import { Bot, Search, Plus, Play, MoreVertical, Edit, Trash } from "lucide-react";
import { toast } from "sonner";
import Link from "next/link";
import { AgentFormDrawer } from "@/components/agents/AgentFormDrawer";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

const SKILLS = [
  "PROJECT_CREATION",
  "ENGINEERING",
  "ARCHITECTURE",
  "PLANNING",
  "DEVELOPMENT_FRONTEND",
  "DEVELOPMENT_BACKEND",
  "TESTING",
  "SECURITY",
  "DOCUMENTATION",
  "DEVOPS",
  "LANDING_PAGE",
  "MARKETING"
];

const PROVIDERS = ["OPENAI", "ANTHROPIC", "GOOGLE", "OLLAMA"];

export default function AgentsPage() {
  const { user } = useAuth();
  const isAdmin = user?.role === "ADMIN";

  const [agents, setAgents] = useState<Agent[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [filterProvider, setFilterProvider] = useState<string>("ALL");
  const [filterSkill, setFilterSkill] = useState<string>("ALL");

  const [drawerOpen, setDrawerOpen] = useState(false);
  const [selectedAgent, setSelectedAgent] = useState<Agent | null>(null);

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        setLoading(true);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const params: any = {};
        if (search) params.search = search;
        if (filterProvider && filterProvider !== "ALL") params.provider = filterProvider;
        if (filterSkill && filterSkill !== "ALL") params.skill = filterSkill;

        const response = await agentService.getAgents(params);
        setAgents(response.items || []);
      } catch (error) {
        console.error(error);
        toast.error("Erro ao carregar agentes");
      } finally {
        setLoading(false);
      }
    };

    fetchAgents();
  }, [search, filterProvider, filterSkill]);

  const fetchAgentsManual = async () => {
      try {
        setLoading(true);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const params: any = {};
        if (search) params.search = search;
        if (filterProvider && filterProvider !== "ALL") params.provider = filterProvider;
        if (filterSkill && filterSkill !== "ALL") params.skill = filterSkill;

        const response = await agentService.getAgents(params);
        setAgents(response.items || []);
      } catch (error) {
        console.error(error);
        toast.error("Erro ao carregar agentes");
      } finally {
        setLoading(false);
      }
  };

  const handleToggleEnabled = async (agent: Agent, checked: boolean) => {
    try {
      if (!isAdmin) {
        toast.error("Apenas administradores podem alterar o status do agente");
        return;
      }
      await agentService.updateAgent(agent.id, {
        name: agent.name,
        description: agent.description,
        provider: agent.provider,
        model: agent.model,
        system_prompts: agent.system_prompts,
        skills: agent.skills,
        api_key_ref: agent.api_key_ref,
        enabled: checked,
      });
      setAgents(agents.map(a => a.id === agent.id ? { ...a, enabled: checked } : a));
      toast.success(`Agente ${agent.name} ${checked ? 'habilitado' : 'desabilitado'}`);
    } catch (error) {
      console.error(error);
      toast.error("Erro ao atualizar status do agente");
      fetchAgentsManual(); // revert on error
    }
  };

  const handleTestConnection = async (agentId: string) => {
    try {
      toast.info("Testando conexão...");
      const response = await agentService.testConnection(agentId);
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

  const handleDelete = async (agentId: string) => {
    if (!confirm("Tem certeza que deseja remover este agente?")) return;
    try {
      await agentService.deleteAgent(agentId);
      toast.success("Agente removido com sucesso");
      fetchAgentsManual();
    } catch (error) {
      console.error(error);
      toast.error("Erro ao remover agente");
    }
  };

  const openEditDrawer = (agent: Agent) => {
    if (!isAdmin) {
      toast.error("Apenas administradores podem editar agentes");
      return;
    }
    setSelectedAgent(agent);
    setDrawerOpen(true);
  };

  const openCreateDrawer = () => {
    if (!isAdmin) {
      toast.error("Apenas administradores podem criar agentes");
      return;
    }
    setSelectedAgent(null);
    setDrawerOpen(true);
  };

  const getProviderColor = (provider: string) => {
    switch (provider) {
      case "OPENAI": return "bg-blue-500/10 text-blue-500 border-blue-500/20";
      case "ANTHROPIC": return "bg-orange-500/10 text-orange-500 border-orange-500/20";
      case "GOOGLE": return "bg-green-500/10 text-green-500 border-green-500/20";
      case "OLLAMA": return "bg-purple-500/10 text-purple-500 border-purple-500/20";
      default: return "bg-gray-500/10 text-gray-500 border-gray-500/20";
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case "IDLE": return "bg-gray-500/20 text-gray-400";
      case "RUNNING": return "bg-blue-500/20 text-blue-400 animate-pulse";
      case "PAUSED": return "bg-yellow-500/20 text-yellow-400";
      case "QUEUED": return "bg-orange-500/20 text-orange-400";
      case "ERROR": return "bg-red-500/20 text-red-400";
      case "COMPLETED": return "bg-green-500/20 text-green-400";
      default: return "bg-gray-500/20 text-gray-400";
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Catálogo de Agentes</h1>
          <p className="text-muted-foreground">Gerencie a biblioteca de especialistas em IA.</p>
        </div>
        {isAdmin && (
          <Button onClick={openCreateDrawer}>
            <Plus className="mr-2 h-4 w-4" />
            Novo Agente
          </Button>
        )}
      </div>

      <div className="flex flex-col sm:flex-row gap-4 bg-card p-4 rounded-lg border">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            placeholder="Buscar agentes..."
            className="pl-8"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
        <Select value={filterProvider} onValueChange={(val) => setFilterProvider(val || "")}>
          <SelectTrigger className="w-full sm:w-[180px]">
            <SelectValue placeholder="Provider" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="ALL">Todos Providers</SelectItem>
            {PROVIDERS.map(p => (
              <SelectItem key={p} value={p}>{p}</SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Select value={filterSkill} onValueChange={(val) => setFilterSkill(val || "")}>
          <SelectTrigger className="w-full sm:w-[180px]">
            <SelectValue placeholder="Skill" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="ALL">Todas Skills</SelectItem>
            {SKILLS.map(s => (
              <SelectItem key={s} value={s}>{s}</SelectItem>
            ))}
          </SelectContent>
        </Select>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {[1, 2, 3, 4, 5, 6].map(i => (
            <Card key={i} className="animate-pulse h-[250px] bg-muted/20" />
          ))}
        </div>
      ) : agents.length === 0 ? (
        <div className="text-center py-12">
          <Bot className="mx-auto h-12 w-12 text-muted-foreground/50 mb-4" />
          <h3 className="text-lg font-medium">Nenhum agente encontrado</h3>
          <p className="text-muted-foreground">Tente ajustar os filtros de busca.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {agents.map((agent) => (
            <Card key={agent.id} className="flex flex-col">
              <CardHeader className="pb-3">
                <div className="flex justify-between items-start">
                  <div className="flex gap-2">
                    <Badge variant="outline" className={getProviderColor(agent.provider)}>
                      {agent.provider}
                    </Badge>
                    <Badge variant="secondary" className={getStatusColor(agent.status || "IDLE")}>
                      {agent.status || "IDLE"}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-2">
                    <Switch
                      checked={agent.enabled}
                      onCheckedChange={(c) => handleToggleEnabled(agent, c)}
                      disabled={!isAdmin}
                    />
                    {isAdmin && (
                      <DropdownMenu>
                        <DropdownMenuTrigger className="inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-clip-padding text-sm font-medium whitespace-nowrap transition-all outline-none select-none hover:bg-muted hover:text-foreground h-8 w-8 size-8">
                          <MoreVertical className="h-4 w-4" />
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuItem onClick={() => openEditDrawer(agent)}>
                            <Edit className="mr-2 h-4 w-4" />
                            Editar
                          </DropdownMenuItem>
                          <DropdownMenuItem className="text-destructive" onClick={() => handleDelete(agent.id)}>
                            <Trash className="mr-2 h-4 w-4" />
                            Remover
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    )}
                  </div>
                </div>
                <Link href={`/dashboard/agents/${agent.id}`}>
                  <CardTitle className="mt-2 text-xl hover:text-primary transition-colors cursor-pointer">
                    {agent.name}
                  </CardTitle>
                </Link>
                <div className="text-xs text-muted-foreground font-mono">{agent.model}</div>
              </CardHeader>
              <CardContent className="flex-1">
                <p className="text-sm text-muted-foreground line-clamp-2 mb-4">
                  {agent.description}
                </p>
                <div className="flex flex-wrap gap-1">
                  {agent.skills?.slice(0, 3).map(skill => (
                    <Badge key={skill} variant="secondary" className="text-[10px] px-1.5 py-0">
                      {skill.replace("DEVELOPMENT_", "")}
                    </Badge>
                  ))}
                  {(agent.skills?.length || 0) > 3 && (
                    <Badge variant="secondary" className="text-[10px] px-1.5 py-0">
                      +{agent.skills.length - 3}
                    </Badge>
                  )}
                </div>
              </CardContent>
              <CardFooter className="pt-3 border-t">
                <div className="flex w-full gap-2">
                  <Button variant="outline" className="flex-1" onClick={() => handleTestConnection(agent.id)}>
                    <Play className="mr-2 h-4 w-4" />
                    Testar
                  </Button>
                  <Link href={`/dashboard/agents/${agent.id}`} className="flex-1 inline-flex shrink-0 items-center justify-center rounded-lg border border-transparent bg-secondary text-secondary-foreground hover:bg-secondary/80 h-8 px-2.5 w-full text-sm font-medium">
                    Detalhes
                  </Link>
                </div>
              </CardFooter>
            </Card>
          ))}
        </div>
      )}

      {drawerOpen && (
        <AgentFormDrawer
          open={drawerOpen}
          onOpenChange={setDrawerOpen}
          agent={selectedAgent}
          onSuccess={fetchAgentsManual}
        />
      )}
    </div>
  );
}
