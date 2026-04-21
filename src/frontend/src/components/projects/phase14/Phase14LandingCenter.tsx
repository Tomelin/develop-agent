"use client";

import { useCallback, useEffect, useMemo, useState, type ReactNode } from "react";
import { Project } from "@/types/project";
import { Phase5CodeFile } from "@/types/phase5";
import { LandingPageManualBrief, Phase14DeliveryReport } from "@/types/phase14";
import { Phase14Service } from "@/services/phase14";
import { ProjectService } from "@/services/project";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Progress } from "@/components/ui/progress";
import { toast } from "sonner";
import { CheckCircle2, Download, ExternalLink, Loader2, Monitor, Play, Tablet, Smartphone, WandSparkles } from "lucide-react";

interface Phase14LandingCenterProps {
  project: Project;
}

const VIEWPORTS = {
  desktop: 1440,
  tablet: 768,
  mobile: 375,
} as const;

const DEFAULT_BRIEF: LandingPageManualBrief = {
  product_name: "",
  problem_solved: "",
  target_audience: "",
  unique_value_proposed: "",
  key_features: ["", "", ""],
  color_palette: ["#111827", "#2563eb", "#f59e0b"],
  theme: "dark",
  communication_tone: "profissional",
  language: "pt-BR",
  preferred_typography: "Inter",
  output_format: "html",
  primary_keyword: "",
  primary_cta: "Começar agora",
  secondary_cta: "Agendar demo",
  social_proof_highlight: "+1.000 empresas usando",
};

export function Phase14LandingCenter({ project }: Phase14LandingCenterProps) {
  const [useLinkedProject, setUseLinkedProject] = useState(Boolean(project.linked_project_id));
  const [generateVariants, setGenerateVariants] = useState(true);
  const [variantCount, setVariantCount] = useState(3);
  const [manualBrief, setManualBrief] = useState<LandingPageManualBrief>(() => {
    if (typeof window === "undefined") return DEFAULT_BRIEF;
    const fromProjectCreation = sessionStorage.getItem(`phase14-brief-draft-new-${project.id}`);
    const persisted = sessionStorage.getItem(`phase14-brief-draft-${project.id}`);
    const raw = fromProjectCreation ?? persisted;
    if (!raw) return DEFAULT_BRIEF;
    try {
      const parsed = JSON.parse(raw) as LandingPageManualBrief;
      if (fromProjectCreation) {
        sessionStorage.removeItem(`phase14-brief-draft-new-${project.id}`);
      }
      return { ...DEFAULT_BRIEF, ...parsed, key_features: parsed.key_features?.length ? parsed.key_features : DEFAULT_BRIEF.key_features };
    } catch {
      return DEFAULT_BRIEF;
    }
  });
  const [running, setRunning] = useState(false);
  const [deliveryReport, setDeliveryReport] = useState<Phase14DeliveryReport | null>(null);
  const [files, setFiles] = useState<Phase5CodeFile[]>([]);
  const [viewport, setViewport] = useState<keyof typeof VIEWPORTS>("desktop");
  const [streamProgress, setStreamProgress] = useState(0);

  const storageKey = `phase14-brief-draft-${project.id}`;

  useEffect(() => {
    if (typeof window === "undefined") return;
    sessionStorage.setItem(storageKey, JSON.stringify(manualBrief));
  }, [manualBrief, storageKey]);

  const fetchArtifacts = useCallback(async () => {
    try {
      const allFiles = await ProjectService.getProjectFiles(project.id);
      setFiles(allFiles.filter((file) => file.path.startsWith("artifacts/landing/")));
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar os artefatos da landing page.");
    }
  }, [project.id]);

  useEffect(() => {
    const bootstrap = setTimeout(() => {
      void fetchArtifacts();
    }, 0);
    const timer = setInterval(() => void fetchArtifacts(), 8000);
    return () => {
      clearTimeout(bootstrap);
      clearInterval(timer);
    };
  }, [fetchArtifacts]);

  const landingHtml = useMemo(() => files.find((file) => file.path === "artifacts/landing/landing_page.html")?.content ?? "", [files]);

  const conversionReport = useMemo(
    () => files.find((file) => file.path === "artifacts/landing/CONVERSION_REPORT.md")?.content ?? "",
    [files],
  );

  const seoChecklist = useMemo(
    () => files.find((file) => file.path === "artifacts/landing/SEO_CHECKLIST.md")?.content ?? "",
    [files],
  );

  useEffect(() => {
    if (!running) return;

    const interval = setInterval(() => {
      setStreamProgress((value) => Math.min(95, value + 5));
    }, 450);

    return () => clearInterval(interval);
  }, [running]);

  const displayHtml = useMemo(() => {
    if (!landingHtml) return "";
    const pct = running ? streamProgress : 100;
    const chars = Math.max(200, Math.floor((landingHtml.length * pct) / 100));
    return landingHtml.slice(0, chars);
  }, [landingHtml, running, streamProgress]);

  const validationErrors = useMemo(() => {
    if (useLinkedProject) return [] as string[];

    const errors: string[] = [];
    if (!manualBrief.product_name.trim()) errors.push("Nome do produto é obrigatório.");
    if (!manualBrief.problem_solved.trim()) errors.push("Problema que resolve é obrigatório.");
    if (!manualBrief.target_audience.trim()) errors.push("Público-alvo é obrigatório.");
    if (!manualBrief.unique_value_proposed.trim()) errors.push("Proposta de valor única é obrigatória.");

    const validFeatures = manualBrief.key_features.map((f) => f.trim()).filter(Boolean);
    if (validFeatures.length < 3) errors.push("Informe no mínimo 3 benefícios principais.");
    if (validFeatures.length > 5) errors.push("Informe no máximo 5 benefícios principais.");

    return errors;
  }, [manualBrief, useLinkedProject]);

  const runPhase14 = async () => {
    if (validationErrors.length > 0) {
      toast.error("Revise os campos obrigatórios antes de iniciar a Tríade.");
      return;
    }

    try {
      setRunning(true);
      setStreamProgress(0);
      const report = await Phase14Service.run(project.id, {
        use_linked_project: useLinkedProject,
        generate_variants: generateVariants,
        variant_count: variantCount,
        manual_brief: {
          ...manualBrief,
          key_features: manualBrief.key_features.map((f) => f.trim()).filter(Boolean),
        },
      });
      setDeliveryReport(report);
      setStreamProgress(100);
      toast.success("Fluxo B executado com sucesso.");
      await fetchArtifacts();
    } catch (error) {
      console.error(error);
      toast.error("Falha na execução da Tríade de Landing Page.");
    } finally {
      setRunning(false);
    }
  };

  const downloadZip = async () => {
    try {
      const blob = await ProjectService.downloadProjectFilesZip(project.id);
      const anchor = document.createElement("a");
      anchor.href = URL.createObjectURL(blob);
      anchor.download = `landing-page-${project.id}.zip`;
      document.body.appendChild(anchor);
      anchor.click();
      URL.revokeObjectURL(anchor.href);
      anchor.remove();
      toast.success("Download do bundle iniciado.");
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível baixar o ZIP.");
    }
  };

  const openFullscreen = () => {
    if (!landingHtml) return;
    const blob = new Blob([landingHtml], { type: "text/html" });
    window.open(URL.createObjectURL(blob), "_blank", "noopener,noreferrer");
  };

  return (
    <div className="space-y-6">
      <Card className="border-border/60 bg-card/70">
        <CardHeader className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <CardTitle className="flex items-center gap-2 text-lg">
              <WandSparkles className="h-5 w-5 text-primary" />
              Fluxo B — Landing Page Studio
            </CardTitle>
            <CardDescription>
              Configure o brief, execute a Tríade e acompanhe preview + score de conversão com padrão premium para cliente final.
            </CardDescription>
          </div>
          <div className="flex items-center gap-2">
            <Button variant="outline" onClick={downloadZip} className="gap-2">
              <Download className="h-4 w-4" /> Download da Landing Page
            </Button>
            <Button onClick={runPhase14} disabled={running} className="gap-2">
              {running ? <Loader2 className="h-4 w-4 animate-spin" /> : <Play className="h-4 w-4" />} Iniciar Tríade
            </Button>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {running && (
            <div className="space-y-2 rounded-xl border border-primary/30 bg-primary/5 p-3">
              <p className="text-sm font-medium">Gerando landing page em tempo real…</p>
              <Progress value={streamProgress} className="h-2" />
            </div>
          )}

          <div className="grid gap-4 lg:grid-cols-2">
            <div className="space-y-3 rounded-xl border border-border/60 bg-background/60 p-4">
              <div className="flex items-center justify-between">
                <Label className="text-sm font-medium">Herdar contexto de projeto existente</Label>
                <Switch checked={useLinkedProject} onCheckedChange={setUseLinkedProject} />
              </div>
              <p className="text-xs text-muted-foreground">
                Quando ativo, usa branding e posicionamento do projeto Fluxo A vinculado. Quando desativado, utiliza brief manual completo.
              </p>

              <div className="flex items-center justify-between pt-2">
                <Label className="text-sm font-medium">Gerar variantes A/B</Label>
                <Switch checked={generateVariants} onCheckedChange={setGenerateVariants} />
              </div>

              <div className="space-y-1">
                <Label>Quantidade de variantes</Label>
                <Select value={String(variantCount)} onValueChange={(value) => setVariantCount(Number(value))}>
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="1">1 variante</SelectItem>
                    <SelectItem value="2">2 variantes</SelectItem>
                    <SelectItem value="3">3 variantes</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            <div className="rounded-xl border border-border/60 bg-background/60 p-4">
              <p className="text-sm font-medium mb-2">Resultado da última execução</p>
              {deliveryReport ? (
                <div className="space-y-2 text-sm">
                  <p className="flex items-center gap-2"><CheckCircle2 className="h-4 w-4 text-green-500" /> Formato: {deliveryReport.output_format.toUpperCase()}</p>
                  <p>Origem do brief: <Badge variant="outline">{deliveryReport.brief_source}</Badge></p>
                  <p>Score de conversão: <span className="font-semibold text-primary">{deliveryReport.conversion_score.toFixed(1)}/100</span></p>
                  <p>Artefatos gerados: {deliveryReport.artifact_paths.length}</p>
                </div>
              ) : (
                <p className="text-sm text-muted-foreground">Ainda não executado para este projeto.</p>
              )}
            </div>
          </div>

          {!useLinkedProject && (
            <div className="grid gap-4 rounded-xl border border-border/60 bg-background/60 p-4 md:grid-cols-2">
              <Field label="Nome do produto *">
                <Input value={manualBrief.product_name} onChange={(e) => setManualBrief((prev) => ({ ...prev, product_name: e.target.value }))} />
              </Field>
              <Field label="Público-alvo *">
                <Input value={manualBrief.target_audience} onChange={(e) => setManualBrief((prev) => ({ ...prev, target_audience: e.target.value }))} />
              </Field>
              <Field label="Problema que resolve *" className="md:col-span-2">
                <Textarea value={manualBrief.problem_solved} onChange={(e) => setManualBrief((prev) => ({ ...prev, problem_solved: e.target.value }))} />
              </Field>
              <Field label="Proposta de valor única *" className="md:col-span-2">
                <Textarea value={manualBrief.unique_value_proposed} onChange={(e) => setManualBrief((prev) => ({ ...prev, unique_value_proposed: e.target.value }))} />
              </Field>

              <div className="space-y-2 md:col-span-2">
                <Label>Benefícios (3 a 5)</Label>
                {manualBrief.key_features.map((feature, index) => (
                  <Input
                    key={index}
                    value={feature}
                    placeholder={`Benefício ${index + 1}`}
                    onChange={(e) => {
                      const next = [...manualBrief.key_features];
                      next[index] = e.target.value;
                      setManualBrief((prev) => ({ ...prev, key_features: next }));
                    }}
                  />
                ))}
              </div>

              <div className="space-y-2 md:col-span-2">
                <Label>Paleta de cores</Label>
                <div className="flex flex-wrap gap-3">
                  {manualBrief.color_palette.map((color, index) => (
                    <div key={index} className="flex items-center gap-2 rounded-lg border border-border/70 px-2 py-1">
                      <input
                        type="color"
                        value={color}
                        onChange={(e) => {
                          const next = [...manualBrief.color_palette];
                          next[index] = e.target.value;
                          setManualBrief((prev) => ({ ...prev, color_palette: next }));
                        }}
                        className="h-8 w-8"
                      />
                      <Input
                        value={color}
                        className="h-8 w-28"
                        onChange={(e) => {
                          const next = [...manualBrief.color_palette];
                          next[index] = e.target.value;
                          setManualBrief((prev) => ({ ...prev, color_palette: next }));
                        }}
                      />
                    </div>
                  ))}
                </div>
              </div>

              <Field label="Tema">
                <Select
                  value={manualBrief.theme}
                  onValueChange={(value) => {
                    if (!value) return;
                    setManualBrief((prev) => ({ ...prev, theme: value as "light" | "dark" }));
                  }}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="dark">Dark</SelectItem>
                    <SelectItem value="light">Light</SelectItem>
                  </SelectContent>
                </Select>
              </Field>

              <Field label="Tom de comunicação">
                <Select
                  value={manualBrief.communication_tone}
                  onValueChange={(value) => {
                    if (!value) return;
                    setManualBrief((prev) => ({ ...prev, communication_tone: value as LandingPageManualBrief["communication_tone"] }));
                  }}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="profissional">Profissional</SelectItem>
                    <SelectItem value="moderno">Moderno</SelectItem>
                    <SelectItem value="descontraído">Descontraído</SelectItem>
                    <SelectItem value="inspirador">Inspirador</SelectItem>
                  </SelectContent>
                </Select>
              </Field>

              <Field label="Idioma">
                <Select
                  value={manualBrief.language}
                  onValueChange={(value) => {
                    if (!value) return;
                    setManualBrief((prev) => ({ ...prev, language: value as LandingPageManualBrief["language"] }));
                  }}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="pt-BR">PT-BR</SelectItem>
                    <SelectItem value="en-US">EN-US</SelectItem>
                    <SelectItem value="es">ES</SelectItem>
                  </SelectContent>
                </Select>
              </Field>

              <Field label="Output preferido">
                <Select
                  value={manualBrief.output_format}
                  onValueChange={(value) => {
                    if (!value) return;
                    setManualBrief((prev) => ({ ...prev, output_format: value as "html" | "nextjs" }));
                  }}
                >
                  <SelectTrigger><SelectValue /></SelectTrigger>
                  <SelectContent>
                    <SelectItem value="html">HTML/CSS/JS</SelectItem>
                    <SelectItem value="nextjs">Next.js</SelectItem>
                  </SelectContent>
                </Select>
              </Field>
            </div>
          )}

          {validationErrors.length > 0 && !useLinkedProject && (
            <div className="rounded-lg border border-destructive/30 bg-destructive/10 p-3 text-sm">
              {validationErrors.map((error) => (
                <p key={error}>• {error}</p>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      <Tabs defaultValue="preview">
        <TabsList className="bg-card/50 border border-border p-1">
          <TabsTrigger value="preview">Preview em tempo real</TabsTrigger>
          <TabsTrigger value="analysis">Análise CRO + SEO</TabsTrigger>
          <TabsTrigger value="variants">Variantes A/B</TabsTrigger>
        </TabsList>

        <TabsContent value="preview" className="mt-5">
          <div className="grid gap-4 lg:grid-cols-2">
            <Card className="border-border/60 bg-card/70">
              <CardHeader>
                <CardTitle className="text-base">Streaming de Código</CardTitle>
                <CardDescription>Renderização progressiva do HTML gerado pela Tríade.</CardDescription>
              </CardHeader>
              <CardContent>
                <pre className="h-[500px] overflow-auto rounded-lg bg-black/90 p-4 text-xs text-green-400">
                  <code>{displayHtml || "Aguardando execução..."}</code>
                </pre>
              </CardContent>
            </Card>

            <Card className="border-border/60 bg-card/70">
              <CardHeader>
                <div className="flex items-center justify-between gap-2">
                  <div>
                    <CardTitle className="text-base">Preview Sandbox</CardTitle>
                    <CardDescription>Split-screen com simulação Desktop/Tablet/Mobile.</CardDescription>
                  </div>
                  <div className="flex items-center gap-1">
                    <Button size="icon" variant={viewport === "desktop" ? "default" : "outline"} onClick={() => setViewport("desktop")}><Monitor className="h-4 w-4" /></Button>
                    <Button size="icon" variant={viewport === "tablet" ? "default" : "outline"} onClick={() => setViewport("tablet")}><Tablet className="h-4 w-4" /></Button>
                    <Button size="icon" variant={viewport === "mobile" ? "default" : "outline"} onClick={() => setViewport("mobile")}><Smartphone className="h-4 w-4" /></Button>
                  </div>
                </div>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex justify-end">
                  <Button size="sm" variant="outline" onClick={openFullscreen} disabled={!landingHtml} className="gap-2">
                    <ExternalLink className="h-4 w-4" /> Abrir em aba completa
                  </Button>
                </div>
                <div className="mx-auto overflow-hidden rounded-xl border border-border/70 bg-white" style={{ width: "100%", maxWidth: VIEWPORTS[viewport], minHeight: 500 }}>
                  <iframe title="Landing Preview" srcDoc={displayHtml} sandbox="allow-same-origin" className="h-[500px] w-full" />
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="analysis" className="mt-5">
          <div className="grid gap-4 lg:grid-cols-2">
            <Card className="border-border/60 bg-card/70">
              <CardHeader>
                <CardTitle className="text-base">Relatório de Conversão</CardTitle>
              </CardHeader>
              <CardContent>
                <pre className="h-[420px] overflow-auto rounded-lg border border-border/60 bg-background/60 p-4 text-xs">{conversionReport || "Sem relatório ainda."}</pre>
              </CardContent>
            </Card>
            <Card className="border-border/60 bg-card/70">
              <CardHeader>
                <CardTitle className="text-base">Checklist SEO</CardTitle>
              </CardHeader>
              <CardContent>
                <pre className="h-[420px] overflow-auto rounded-lg border border-border/60 bg-background/60 p-4 text-xs">{seoChecklist || "Sem checklist ainda."}</pre>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="variants" className="mt-5">
          <Card className="border-border/60 bg-card/70">
            <CardHeader>
              <CardTitle className="text-base">Variantes geradas</CardTitle>
              <CardDescription>Comparação side-by-side dos scores para decisão de A/B test.</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              {deliveryReport?.variants?.length ? deliveryReport.variants.map((variant) => (
                <div key={variant.name} className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-border/60 p-3">
                  <div>
                    <p className="font-medium">{variant.name}</p>
                    <p className="text-xs text-muted-foreground">{variant.path}</p>
                  </div>
                  <Badge variant="secondary">Score {variant.conversion_score.toFixed(1)}</Badge>
                </div>
              )) : <p className="text-sm text-muted-foreground">Nenhuma variante disponível. Ative “Gerar variantes A/B” e execute a Tríade.</p>}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}

function Field({ label, className, children }: { label: string; className?: string; children: ReactNode }) {
  return (
    <div className={`space-y-2 ${className ?? ""}`}>
      <Label>{label}</Label>
      {children}
    </div>
  );
}
