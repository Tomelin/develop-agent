"use client";

import { useMemo, useState } from "react";
import { PhaseArtifact } from "@/types/phase8";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Download, Copy, Maximize2, GitCompareArrows, FileText } from "lucide-react";
import { toast } from "sonner";

type HeadingItem = { id: string; label: string; level: 1 | 2 | 3 };

const buildToc = (content: string): HeadingItem[] => {
  return content
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.startsWith("#"))
    .map((line, index) => {
      const level = (line.match(/^#+/)?.[0].length ?? 1) as 1 | 2 | 3;
      const label = line.replace(/^#+\s*/, "");
      return { id: `heading-${index}`, label, level: Math.min(level, 3) as 1 | 2 | 3 };
    });
};

const markdownToHtml = (raw: string) => {
  const escaped = raw
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");

  return escaped
    .replace(/^###\s(.+)$/gm, '<h3 class="mt-4 text-base font-semibold">$1</h3>')
    .replace(/^##\s(.+)$/gm, '<h2 class="mt-5 text-lg font-semibold">$1</h2>')
    .replace(/^#\s(.+)$/gm, '<h1 class="mt-6 text-2xl font-bold">$1</h1>')
    .replace(/```mermaid([\s\S]*?)```/gm, '<div class="my-4 rounded-lg border border-primary/30 bg-primary/5 p-3"><p class="mb-2 text-xs uppercase tracking-wide text-primary">Mermaid</p><pre class="overflow-x-auto text-xs">$1</pre></div>')
    .replace(/```([a-zA-Z0-9_-]*)\n([\s\S]*?)```/gm, '<pre class="my-4 overflow-x-auto rounded-lg border bg-background p-3 text-xs"><code data-lang="$1">$2</code></pre>')
    .replace(/\*\*(.*?)\*\*/g, "<strong>$1</strong>")
    .replace(/`([^`]+)`/g, '<code class="rounded bg-muted px-1 py-0.5">$1</code>')
    .replace(/\n\n/g, "<br/><br/>");
};

const getDiffLines = (base: string, next: string) => {
  const baseLines = base.split("\n");
  const nextLines = next.split("\n");
  return nextLines.map((line, index) => ({
    line,
    changed: baseLines[index] !== line,
  }));
};

export function ArtifactViewer({ artifact }: { artifact: PhaseArtifact }) {
  const [selectedVersion, setSelectedVersion] = useState<number>(0);
  const versions = artifact.versions?.length ? artifact.versions : [{ version: 1, content: artifact.current_content, created_at: artifact.updated_at }];
  const active = versions[Math.min(selectedVersion, versions.length - 1)];
  const previous = versions[Math.max(0, Math.min(selectedVersion - 1, versions.length - 1))];

  const toc = useMemo(() => buildToc(active.content), [active.content]);

  const handleCopy = async (content: string) => {
    await navigator.clipboard.writeText(content);
    toast.success("Conteúdo copiado");
  };

  const handleDownloadMd = () => {
    const blob = new Blob([active.content], { type: "text/markdown;charset=utf-8" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `${artifact.title.replace(/\s+/g, "-").toLowerCase()}.md`;
    a.click();
    URL.revokeObjectURL(url);
  };

  const handleDownloadPdf = () => {
    window.print();
  };

  return (
    <Card className="border-border/60 bg-card/60 backdrop-blur-sm">
      <CardHeader className="flex flex-row items-center justify-between gap-3">
        <div>
          <CardTitle className="text-base">{artifact.title}</CardTitle>
          <p className="mt-1 text-xs text-muted-foreground">Atualizado em {new Date(artifact.updated_at).toLocaleString()}</p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">{artifact.track}</Badge>
          <Badge variant="secondary">{artifact.type}</Badge>
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-center gap-2">
          <Button size="sm" variant="outline" onClick={handleDownloadMd}><Download className="mr-2 h-4 w-4" />Markdown</Button>
          <Button size="sm" variant="outline" onClick={handleDownloadPdf}><Download className="mr-2 h-4 w-4" />PDF</Button>
          <Button size="sm" variant="outline" onClick={() => handleCopy(active.content)}><Copy className="mr-2 h-4 w-4" />Copiar</Button>
          <Dialog>
            <DialogTrigger render={<Button size="sm" variant="outline" />}>
              <Maximize2 className="mr-2 h-4 w-4" />Fullscreen
            </DialogTrigger>
            <DialogContent className="max-w-[92vw] sm:max-w-[92vw] h-[88vh] overflow-hidden">
              <DialogHeader>
                <DialogTitle>{artifact.title}</DialogTitle>
              </DialogHeader>
              <div className="overflow-y-auto pr-2 text-sm" dangerouslySetInnerHTML={{ __html: markdownToHtml(active.content) }} />
            </DialogContent>
          </Dialog>
        </div>

        <Tabs defaultValue="read" className="w-full">
          <TabsList className="w-full justify-start">
            <TabsTrigger value="read"><FileText className="mr-2 h-4 w-4" />Leitura</TabsTrigger>
            <TabsTrigger value="diff"><GitCompareArrows className="mr-2 h-4 w-4" />Diff</TabsTrigger>
          </TabsList>
          <TabsContent value="read" className="mt-4 grid gap-4 lg:grid-cols-[240px_1fr]">
            <aside className="hidden rounded-lg border border-border/70 bg-background/60 p-3 lg:block">
              <p className="mb-2 text-xs font-semibold uppercase text-muted-foreground">Índice</p>
              <ul className="space-y-1 text-xs">
                {toc.map((item) => (
                  <li key={item.id} className={item.level === 1 ? "pl-0" : item.level === 2 ? "pl-3" : "pl-6"}>{item.label}</li>
                ))}
              </ul>
            </aside>
            <article className="rounded-lg border border-border/70 bg-background/40 p-4 text-sm leading-relaxed" dangerouslySetInnerHTML={{ __html: markdownToHtml(active.content) }} />
          </TabsContent>
          <TabsContent value="diff" className="mt-4">
            <div className="rounded-lg border border-border/70 bg-background/40 p-3">
              {getDiffLines(previous.content, active.content).map((item, index) => (
                <div key={`${item.line}-${index}`} className={`font-mono text-xs ${item.changed ? "bg-primary/15 text-primary" : "text-muted-foreground"}`}>
                  {item.changed ? "+ " : "  "}{item.line || " "}
                </div>
              ))}
            </div>
          </TabsContent>
        </Tabs>

        {versions.length > 1 && (
          <div className="flex items-center gap-2 text-xs">
            {versions.map((version, index) => (
              <Button key={version.version} size="sm" variant={index === selectedVersion ? "default" : "outline"} onClick={() => setSelectedVersion(index)}>
                v{version.version}
              </Button>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
}
