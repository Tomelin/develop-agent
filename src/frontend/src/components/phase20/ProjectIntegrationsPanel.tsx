"use client";

import { useCallback, useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { toast } from "sonner";
import { Phase20Service } from "@/services/phase20";
import { IntegrationConnectionState } from "@/types/phase20";
import { CheckCircle2, GitBranch, MessageSquare, RefreshCcw } from "lucide-react";

export function ProjectIntegrationsPanel({ projectId }: { projectId: string }) {
  const [items, setItems] = useState<IntegrationConnectionState[]>([]);
  const [jiraBaseUrl, setJiraBaseUrl] = useState("");
  const [jiraEmail, setJiraEmail] = useState("");
  const [jiraApiToken, setJiraApiToken] = useState("");
  const [jiraProjectKey, setJiraProjectKey] = useState("");
  const [slackWebhookUrl, setSlackWebhookUrl] = useState("");
  const [slackChannel, setSlackChannel] = useState("#general");

  const loadIntegrations = useCallback(async () => {
    try {
      const data = await Phase20Service.getIntegrationsStatus(projectId);
      setItems(data);
    } catch (error) {
      console.error(error);
      toast.error("Falha ao carregar status de integrações.");
    }
  }, [projectId]);

  useEffect(() => {
    const timer = setTimeout(() => {
      void loadIntegrations();
    }, 0);
    return () => clearTimeout(timer);
  }, [loadIntegrations]);

  const connectGithub = async () => {
    try {
      const { auth_url } = await Phase20Service.getGithubAuthUrl();
      window.location.assign(auth_url);
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível iniciar OAuth do GitHub.");
    }
  };

  const saveJira = async () => {
    try {
      await Phase20Service.configureJiraIntegration({
        base_url: jiraBaseUrl,
        email: jiraEmail,
        api_token: jiraApiToken,
        project_key: jiraProjectKey,
      });
      toast.success("Integração com Jira configurada.");
      await loadIntegrations();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao configurar Jira.");
    }
  };

  const syncJira = async () => {
    try {
      await Phase20Service.syncProjectToJira(projectId);
      toast.success("Sincronização com Jira iniciada.");
      await loadIntegrations();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao sincronizar com Jira.");
    }
  };

  const saveSlack = async () => {
    try {
      await Phase20Service.configureSlackWebhook({ webhook_url: slackWebhookUrl, channel: slackChannel });
      toast.success("Webhook do Slack configurado.");
      await loadIntegrations();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao configurar Slack.");
    }
  };

  const integrationState = (provider: IntegrationConnectionState["provider"]) => items.find((item) => item.provider === provider)?.connected;

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle>Integrações (GitHub, Jira e Slack)</CardTitle>
        <CardDescription>Conecte seu projeto às ferramentas do ciclo de desenvolvimento institucional.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="rounded-xl border p-4">
          <div className="mb-2 flex items-center justify-between">
            <p className="font-medium flex items-center gap-2"><GitBranch className="h-4 w-4" /> GitHub OAuth</p>
            {integrationState("github") && <span className="text-xs text-green-500 flex items-center gap-1"><CheckCircle2 className="h-3 w-3" /> Conectado</span>}
          </div>
          <Button onClick={connectGithub}>Conectar com GitHub</Button>
        </div>

        <div className="rounded-xl border p-4 space-y-2">
          <p className="font-medium">Jira Cloud</p>
          <div className="grid gap-2 md:grid-cols-2">
            <Input placeholder="https://empresa.atlassian.net" value={jiraBaseUrl} onChange={(e) => setJiraBaseUrl(e.target.value)} />
            <Input placeholder="email@empresa.com" value={jiraEmail} onChange={(e) => setJiraEmail(e.target.value)} />
            <Input placeholder="API Token" value={jiraApiToken} onChange={(e) => setJiraApiToken(e.target.value)} type="password" />
            <Input placeholder="Project Key (ex: ENG)" value={jiraProjectKey} onChange={(e) => setJiraProjectKey(e.target.value)} />
          </div>
          <div className="flex gap-2">
            <Button onClick={saveJira}>Salvar configuração</Button>
            <Button variant="outline" onClick={syncJira}><RefreshCcw className="mr-2 h-4 w-4" /> Sincronizar roadmap</Button>
          </div>
        </div>

        <div className="rounded-xl border p-4 space-y-2">
          <p className="font-medium flex items-center gap-2"><MessageSquare className="h-4 w-4" /> Slack Webhook</p>
          <div className="grid gap-2 md:grid-cols-2">
            <Input placeholder="https://hooks.slack.com/services/..." value={slackWebhookUrl} onChange={(e) => setSlackWebhookUrl(e.target.value)} />
            <Input placeholder="#canal-alertas" value={slackChannel} onChange={(e) => setSlackChannel(e.target.value)} />
          </div>
          <Button onClick={saveSlack}>Salvar webhook</Button>
        </div>
      </CardContent>
    </Card>
  );
}
