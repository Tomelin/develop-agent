"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Phase17Service } from "@/services/phase17";
import { AdminPlatformSettings } from "@/types/phase17";
import { toast } from "sonner";

const emptySettings: AdminPlatformSettings = {
  workers: { max_concurrency: 4, agent_timeout_seconds: 300, triad_timeout_seconds: 1200 },
  models: { default_model: "gpt-4o-mini", spec_generation_model: "gpt-4o-mini" },
  limits: { max_projects_per_user: 20, max_parallel_phases_per_user: 2, max_spec_tokens: 4000 },
  retry: { max_attempts: 3, backoff_seconds: 5 },
};

const num = (v: string) => Number(v || 0);

export function AdminSettingsPanel() {
  const [settings, setSettings] = useState<AdminPlatformSettings>(emptySettings);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    Phase17Service.getAdminSettings().then(setSettings).catch(console.error);
  }, []);

  const save = async () => {
    setSaving(true);
    try {
      const response = await Phase17Service.saveAdminSettings(settings);
      setSettings(response);
      toast.success("Configurações globais atualizadas.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao atualizar configurações globais.");
    } finally {
      setSaving(false);
    }
  };

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle>Configurações Globais da Plataforma</CardTitle>
      </CardHeader>
      <CardContent className="space-y-6">
        <section className="space-y-2">
          <h3 className="font-semibold">Workers</h3>
          <div className="grid md:grid-cols-3 gap-3">
            <div><Label>Máx. concorrência</Label><Input type="number" value={settings.workers.max_concurrency} onChange={(e) => setSettings((s) => ({ ...s, workers: { ...s.workers, max_concurrency: num(e.target.value) } }))} /></div>
            <div><Label>Timeout por agente (s)</Label><Input type="number" value={settings.workers.agent_timeout_seconds} onChange={(e) => setSettings((s) => ({ ...s, workers: { ...s.workers, agent_timeout_seconds: num(e.target.value) } }))} /></div>
            <div><Label>Timeout da tríade (s)</Label><Input type="number" value={settings.workers.triad_timeout_seconds} onChange={(e) => setSettings((s) => ({ ...s, workers: { ...s.workers, triad_timeout_seconds: num(e.target.value) } }))} /></div>
          </div>
        </section>

        <section className="space-y-2">
          <h3 className="font-semibold">Modelos</h3>
          <div className="grid md:grid-cols-2 gap-3">
            <div><Label>Modelo padrão</Label><Input value={settings.models.default_model} onChange={(e) => setSettings((s) => ({ ...s, models: { ...s.models, default_model: e.target.value } }))} /></div>
            <div><Label>Modelo para SPEC.md</Label><Input value={settings.models.spec_generation_model} onChange={(e) => setSettings((s) => ({ ...s, models: { ...s.models, spec_generation_model: e.target.value } }))} /></div>
          </div>
        </section>

        <section className="space-y-2">
          <h3 className="font-semibold">Limites e Retry</h3>
          <div className="grid md:grid-cols-3 gap-3">
            <div><Label>Máx. projetos por usuário</Label><Input type="number" value={settings.limits.max_projects_per_user} onChange={(e) => setSettings((s) => ({ ...s, limits: { ...s.limits, max_projects_per_user: num(e.target.value) } }))} /></div>
            <div><Label>Fases simultâneas</Label><Input type="number" value={settings.limits.max_parallel_phases_per_user} onChange={(e) => setSettings((s) => ({ ...s, limits: { ...s.limits, max_parallel_phases_per_user: num(e.target.value) } }))} /></div>
            <div><Label>Limite de tokens no SPEC</Label><Input type="number" value={settings.limits.max_spec_tokens} onChange={(e) => setSettings((s) => ({ ...s, limits: { ...s.limits, max_spec_tokens: num(e.target.value) } }))} /></div>
            <div><Label>Máx. tentativas</Label><Input type="number" value={settings.retry.max_attempts} onChange={(e) => setSettings((s) => ({ ...s, retry: { ...s.retry, max_attempts: num(e.target.value) } }))} /></div>
            <div><Label>Backoff (s)</Label><Input type="number" value={settings.retry.backoff_seconds} onChange={(e) => setSettings((s) => ({ ...s, retry: { ...s.retry, backoff_seconds: num(e.target.value) } }))} /></div>
          </div>
        </section>

        <Button onClick={save} disabled={saving}>{saving ? "Salvando..." : "Salvar configurações"}</Button>
      </CardContent>
    </Card>
  );
}
