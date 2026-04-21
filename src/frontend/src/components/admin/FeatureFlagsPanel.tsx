"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Switch } from "@/components/ui/switch";
import { Button } from "@/components/ui/button";
import { Phase17Service } from "@/services/phase17";
import { dynamicFeatureFlagsSeed, FeatureFlag } from "@/types/phase17";
import { toast } from "sonner";

export function FeatureFlagsPanel() {
  const [flags, setFlags] = useState<FeatureFlag[]>(dynamicFeatureFlagsSeed);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    Phase17Service.getFeatureFlags().then(setFlags).catch(console.error);
  }, []);

  const save = async () => {
    setSaving(true);
    try {
      const data = await Phase17Service.updateFeatureFlags(flags);
      setFlags(data);
      toast.success("Feature flags atualizadas com sucesso.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao salvar feature flags.");
    } finally {
      setSaving(false);
    }
  };

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle>Feature Flags</CardTitle>
      </CardHeader>
      <CardContent className="space-y-3">
        {flags.map((flag) => (
          <div key={flag.key} className="flex items-center justify-between rounded-lg border p-3 bg-background/70">
            <div>
              <p className="text-sm font-semibold">{flag.key}</p>
              <p className="text-xs text-muted-foreground">{flag.description}</p>
            </div>
            <Switch checked={flag.enabled} onCheckedChange={(checked) => setFlags((prev) => prev.map((f) => f.key === flag.key ? { ...f, enabled: checked } : f))} />
          </div>
        ))}
        <Button onClick={save} disabled={saving}>{saving ? "Salvando..." : "Salvar flags"}</Button>
      </CardContent>
    </Card>
  );
}
