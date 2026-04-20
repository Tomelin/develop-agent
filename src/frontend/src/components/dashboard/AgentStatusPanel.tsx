"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Activity, ServerCrash, Bot, CheckCircle2, PauseCircle, Clock } from "lucide-react";
import { AgentService } from "@/services/agent";
import { Agent, AgentStatus } from "@/types/agent";
import { api } from "@/services/api";

export function AgentStatusPanel() {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const fetchAgents = async () => {
      try {
        const response = await AgentService.getAgents(1, 100);
        setAgents(response.items || []);
      } catch (error) {
        console.error("Failed to fetch agents:", error);
        setIsConnected(false);
      }
    };

    // Initial fetch
    fetchAgents();

    // Setup SSE
    const eventSource: EventSource | null = null;

    // Fallback polling setup
    let pollingInterval: NodeJS.Timeout | null = null;

    const startPolling = () => {
      setIsConnected(true); // Treat polling as connected
      pollingInterval = setInterval(() => {
        fetchAgents();
      }, 5000);
    };

    const setupSSE = () => {
      // In a real scenario, you'd need a way to pass auth token to EventSource,
      // which native EventSource doesn't support well with headers.
      // Often, this is done via query params or a polyfill.
      // Here we will use polling as a reliable fallback for now since we have a JWT auth system that might not play well with default EventSource.
      // We will attempt a fast poll for "real-time" feel if SSE is tricky with auth.

      startPolling();
    };

    setupSSE();

    return () => {
      if (eventSource) {
        (eventSource as EventSource).close();
      }
      if (pollingInterval) {
        clearInterval(pollingInterval);
      }
    };
  }, []);

  const getStatusConfig = (status: AgentStatus) => {
    switch (status) {
      case "IDLE":
        return { icon: Bot, color: "text-gray-500", bg: "bg-gray-500/10", border: "border-gray-500/20" };
      case "RUNNING":
        return { icon: Activity, color: "text-green-500 animate-pulse", bg: "bg-green-500/10", border: "border-green-500/20" };
      case "PAUSED":
        return { icon: PauseCircle, color: "text-yellow-500", bg: "bg-yellow-500/10", border: "border-yellow-500/20" };
      case "QUEUED":
        return { icon: Clock, color: "text-blue-500", bg: "bg-blue-500/10", border: "border-blue-500/20" };
      case "ERROR":
        return { icon: ServerCrash, color: "text-destructive", bg: "bg-destructive/10", border: "border-destructive/20" };
      case "COMPLETED":
        return { icon: CheckCircle2, color: "text-green-600", bg: "bg-green-600/10", border: "border-green-600/20" };
      default:
        return { icon: Bot, color: "text-gray-500", bg: "bg-gray-500/10", border: "border-gray-500/20" };
    }
  };

  return (
    <Card className="h-full bg-card/50 backdrop-blur-sm flex flex-col">
      <CardHeader className="pb-3 border-b border-border/50">
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg font-medium flex items-center gap-2">
            <Bot className="h-5 w-5 text-primary" />
            Status dos Agentes
          </CardTitle>
          <div className="flex items-center gap-2">
            <span className="relative flex h-3 w-3">
              {isConnected ? (
                <>
                  <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
                  <span className="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
                </>
              ) : (
                <span className="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
              )}
            </span>
            <span className="text-xs text-muted-foreground">
              {isConnected ? "Conectado" : "Desconectado"}
            </span>
          </div>
        </div>
      </CardHeader>
      <CardContent className="flex-1 overflow-auto p-0">
        <div className="divide-y divide-border/50">
          {agents.length === 0 ? (
            <div className="p-4 text-center text-sm text-muted-foreground py-8">
              Nenhum agente configurado no sistema.
            </div>
          ) : (
            agents.map((agent) => {
              const statusConfig = getStatusConfig(agent.status);
              const StatusIcon = statusConfig.icon;

              return (
                <div key={agent.id} className={`p-4 flex items-start gap-4 transition-colors hover:bg-muted/50`}>
                  <div className={`mt-1 p-2 rounded-full ${statusConfig.bg} ${statusConfig.border} border`}>
                    <StatusIcon className={`h-4 w-4 ${statusConfig.color}`} />
                  </div>
                  <div className="flex-1 space-y-1">
                    <div className="flex items-center justify-between">
                      <p className="text-sm font-medium leading-none">{agent.name}</p>
                      <Badge variant="outline" className={`text-[10px] uppercase ${statusConfig.color}`}>
                        {agent.status}
                      </Badge>
                    </div>
                    <p className="text-xs text-muted-foreground flex items-center gap-2">
                      <span className="font-mono">{agent.provider}</span> • {agent.model}
                    </p>
                    {agent.status === 'RUNNING' && (
                      <p className="text-xs text-primary mt-2">Executando task atual...</p>
                    )}
                  </div>
                </div>
              );
            })
          )}
        </div>
      </CardContent>
    </Card>
  );
}