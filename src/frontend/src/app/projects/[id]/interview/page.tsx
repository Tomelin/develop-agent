"use client";

import { use, useCallback, useEffect, useMemo, useState } from "react";
import { useRouter } from "next/navigation";
import { ArrowLeft, Bot, CheckCircle2, CircleDashed, Copy, Download, Loader2, MessageCircle, RefreshCcw, SendHorizonal, UserRound } from "lucide-react";
import Link from "next/link";
import { toast } from "sonner";
import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Progress } from "@/components/ui/progress";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { Separator } from "@/components/ui/separator";
import { Markdown } from "@/components/ui/markdown";
import { InterviewService } from "@/services/interview";
import { InterviewCoverageItem, InterviewMessage, InterviewSession } from "@/types/interview";

const DEFAULT_COVERAGE: InterviewCoverageItem[] = [
  { key: "problem", title: "Problema de negócio", status: "PENDING" },
  { key: "users", title: "Público-alvo", status: "PENDING" },
  { key: "differential", title: "Diferencial competitivo", status: "PENDING" },
  { key: "mvp", title: "Funcionalidades críticas do MVP", status: "PENDING" },
  { key: "tech", title: "Tecnologias e restrições", status: "PENDING" },
  { key: "timeline", title: "Prazo e urgência", status: "PENDING" },
  { key: "integrations", title: "Integrações externas", status: "PENDING" },
  { key: "business", title: "Modelo de negócio", status: "PENDING" },
];

export default function InterviewPage({ params }: { params: Promise<{ id: string }> }) {
  const resolvedParams = use(params);
  const router = useRouter();
  const [session, setSession] = useState<InterviewSession | null>(null);
  const [loading, setLoading] = useState(true);
  const [sending, setSending] = useState(false);
  const [input, setInput] = useState("");
  const [streamingAnswer, setStreamingAnswer] = useState("");
  const [processingVision, setProcessingVision] = useState(false);

  const refreshSession = useCallback(async () => {
    const data = await InterviewService.getSession(resolvedParams.id);
    setSession(data);
  }, [resolvedParams.id]);

  useEffect(() => {
    const load = async () => {
      try {
        await refreshSession();
      } catch (error) {
        console.error(error);
        toast.error("Não foi possível carregar a entrevista.");
        router.push(`/projects/${resolvedParams.id}`);
      } finally {
        setLoading(false);
      }
    };
    load();
  }, [refreshSession, resolvedParams.id, router]);

  const totalIterations = session?.max_iterations ?? 10;
  const usedIterations = session?.iteration_count ?? 0;
  const coverage = session?.coverage?.length ? session.coverage : DEFAULT_COVERAGE;
  const completion = Math.round((usedIterations / totalIterations) * 100);
  const canConfirm = usedIterations >= 3 && session?.status !== "COMPLETED";

  const historyMessages = useMemo(() => {
    const base = session?.messages ?? [];
    if (!streamingAnswer) return base;
    return [...base, { role: "ASSISTANT", content: streamingAnswer, timestamp: new Date().toISOString() } as InterviewMessage];
  }, [session?.messages, streamingAnswer]);

  const handleSend = async () => {
    if (!input.trim() || sending || !session) return;
    setSending(true);
    setStreamingAnswer("");
    const userMessage = input.trim();
    setInput("");

    setSession((prev) => {
      if (!prev) return prev;
      return {
        ...prev,
        messages: [...prev.messages, { role: "USER", content: userMessage, timestamp: new Date().toISOString() }],
      };
    });

    try {
      await InterviewService.streamMessage(resolvedParams.id, userMessage, {
        onToken: (token) => setStreamingAnswer((prev) => prev + token),
        onDone: async () => {
          setStreamingAnswer("");
          await refreshSession();
        },
        onError: (message) => {
          toast.error(message);
        },
      });
    } catch (error) {
      console.error(error);
      toast.error("Erro ao enviar mensagem.");
    } finally {
      setSending(false);
    }
  };

  const handleRegenerate = async () => {
    try {
      setProcessingVision(true);
      await InterviewService.regenerateVision(resolvedParams.id);
      await refreshSession();
      toast.success("Documento de visão regenerado com sucesso.");
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível regenerar o documento.");
    } finally {
      setProcessingVision(false);
    }
  };

  const handleConfirm = async () => {
    try {
      setProcessingVision(true);
      await InterviewService.confirmInterview(resolvedParams.id);
      await refreshSession();
      toast.success("Visão confirmada. Fase 1 finalizada com sucesso.");
    } catch (error) {
      console.error(error);
      toast.error("Erro ao confirmar visão.");
    } finally {
      setProcessingVision(false);
    }
  };

  const downloadVisionMarkdown = () => {
    const markdown = session?.vision_markdown;
    if (!markdown) return;
    const blob = new Blob([markdown], { type: "text/markdown;charset=utf-8" });
    const link = document.createElement("a");
    link.href = URL.createObjectURL(blob);
    link.download = `VISION-${resolvedParams.id}.md`;
    link.click();
    URL.revokeObjectURL(link.href);
  };

  const downloadVisionPdf = () => {
    const markdown = session?.vision_markdown;
    if (!markdown) return;
    const printWindow = window.open("", "_blank");
    if (!printWindow) return;
    printWindow.document.write(`<pre style="font-family:Inter,system-ui;padding:24px;white-space:pre-wrap;">${markdown}</pre>`);
    printWindow.document.close();
    printWindow.print();
  };

  if (loading || !session) {
    return (
      <div className="flex min-h-[40vh] items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-primary" />
      </div>
    );
  }

  return (
    <PrivateRoute>
      <div className="mx-auto flex w-full max-w-7xl flex-col gap-6 px-4 py-8 lg:px-10">
        <div className="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-border/60 bg-card/60 p-5 shadow-[0_0_60px_-35px_var(--color-primary)] backdrop-blur-xl">
          <div className="space-y-2">
            <Button variant="ghost" asChild className="-ml-2 px-2 text-muted-foreground">
              <Link href={`/projects/${resolvedParams.id}`}>
                <ArrowLeft className="mr-2 h-4 w-4" />
                Voltar para projeto
              </Link>
            </Button>
            <h1 className="text-2xl font-semibold tracking-tight">Entrevista de Descoberta do Produto</h1>
            <p className="text-sm text-muted-foreground">Converse com o Agente Entrevistador para consolidar a visão inicial do produto.</p>
          </div>
          <div className="min-w-[240px] space-y-2">
            <div className="flex items-center justify-between">
              <Badge variant="outline">{usedIterations} de {totalIterations} iterações utilizadas</Badge>
              <span className="text-xs text-muted-foreground">{totalIterations - usedIterations} restantes</span>
            </div>
            <Progress value={completion} className="h-2" />
          </div>
        </div>

        <div className="grid gap-6 lg:grid-cols-[1fr_320px]">
          <Card className="overflow-hidden border-border/70 bg-card/50">
            <CardHeader className="border-b border-border/50">
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg">Chat da Entrevista</CardTitle>
                <Sheet>
                  <SheetTrigger asChild>
                    <Button variant="outline" size="sm">Pré-visualizar VISION.md</Button>
                  </SheetTrigger>
                  <SheetContent className="w-full overflow-y-auto border-border/80 bg-card sm:max-w-3xl">
                    <SheetHeader>
                      <SheetTitle>Documento de Visão do Produto</SheetTitle>
                      <SheetDescription>
                        Última geração: {session.vision_generated_at ? new Date(session.vision_generated_at).toLocaleString("pt-BR") : "ainda não gerado"}
                      </SheetDescription>
                    </SheetHeader>
                    <div className="mt-6 space-y-4">
                      <div className="rounded-xl border border-border/60 bg-background/60 p-4">
                        <Markdown content={session.vision_markdown || "# Documento ainda não gerado\n\nUse o botão **Regenerar Documento** para gerar a versão inicial."} />
                      </div>
                      <div className="flex flex-wrap gap-2">
                        <Button onClick={handleRegenerate} disabled={processingVision}>
                          <RefreshCcw className="mr-2 h-4 w-4" />
                          Regenerar Documento
                        </Button>
                        <Button variant="secondary" disabled={!canConfirm || processingVision} onClick={handleConfirm}>
                          <CheckCircle2 className="mr-2 h-4 w-4" />
                          Confirmar Visão e Avançar
                        </Button>
                        <Button variant="outline" onClick={downloadVisionMarkdown} disabled={!session.vision_markdown}>
                          <Download className="mr-2 h-4 w-4" />
                          Download .md
                        </Button>
                        <Button variant="outline" onClick={downloadVisionPdf} disabled={!session.vision_markdown}>
                          <Download className="mr-2 h-4 w-4" />
                          Download PDF
                        </Button>
                      </div>
                    </div>
                  </SheetContent>
                </Sheet>
              </div>
            </CardHeader>
            <CardContent className="p-0">
              <ScrollArea className="h-[58vh] px-4 py-6">
                <div className="space-y-4">
                  {historyMessages.map((message, index) => (
                    <div key={`${message.timestamp}-${index}`} className={`flex ${message.role === "USER" ? "justify-end" : "justify-start"}`}>
                      <div className={`max-w-[80%] rounded-2xl border px-4 py-3 ${
                        message.role === "USER"
                          ? "border-primary/40 bg-primary/18 text-foreground"
                          : "border-border/70 bg-card text-foreground"
                      }`}>
                        <div className="mb-2 flex items-center gap-2 text-xs text-muted-foreground">
                          {message.role === "USER" ? <UserRound className="h-3.5 w-3.5" /> : <Bot className="h-3.5 w-3.5 text-primary" />}
                          <span>{message.role === "USER" ? "Você" : "Agente Entrevistador"}</span>
                        </div>
                        {message.role === "ASSISTANT" ? (
                          <Markdown content={message.content} />
                        ) : (
                          <p className="text-sm leading-relaxed">{message.content}</p>
                        )}
                      </div>
                    </div>
                  ))}
                  {sending && !streamingAnswer && (
                    <div className="flex justify-start">
                      <div className="inline-flex items-center gap-2 rounded-xl border border-border/60 bg-card px-4 py-3 text-xs text-muted-foreground">
                        <Loader2 className="h-4 w-4 animate-spin text-primary" />
                        Agente está digitando...
                      </div>
                    </div>
                  )}
                </div>
              </ScrollArea>
              <Separator />
              <div className="flex items-center gap-3 p-4">
                <Input
                  value={input}
                  onChange={(event) => setInput(event.target.value)}
                  onKeyDown={(event) => {
                    if (event.key === "Enter" && !event.shiftKey) {
                      event.preventDefault();
                      void handleSend();
                    }
                  }}
                  placeholder="Descreva sua ideia, contexto ou feedback..."
                  disabled={sending || usedIterations >= totalIterations}
                  className="h-11"
                />
                <Button onClick={handleSend} disabled={sending || !input.trim() || usedIterations >= totalIterations}>
                  <SendHorizonal className="h-4 w-4" />
                </Button>
              </div>
            </CardContent>
          </Card>

          <div className="space-y-4">
            <Card className="border-border/70 bg-card/50">
              <CardHeader>
                <CardTitle className="text-base">Cobertura da Entrevista</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {coverage.map((item) => (
                  <div key={item.key} className="flex items-center justify-between rounded-lg border border-border/60 bg-background/40 px-3 py-2">
                    <span className="text-sm">{item.title}</span>
                    {item.status === "DONE" ? (
                      <Badge variant="secondary" className="gap-1"><CheckCircle2 className="h-3 w-3" />OK</Badge>
                    ) : (
                      <Badge variant="outline" className="gap-1 text-muted-foreground"><CircleDashed className="h-3 w-3" />Pendente</Badge>
                    )}
                  </div>
                ))}
              </CardContent>
            </Card>

            <Card className="border-border/70 bg-card/50">
              <CardHeader>
                <CardTitle className="flex items-center gap-2 text-base">
                  <MessageCircle className="h-4 w-4 text-primary" />
                  Histórico da Entrevista
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <p className="text-xs text-muted-foreground">
                  Sessão somente leitura para auditoria de decisões e alinhamento futuro.
                </p>
                <ScrollArea className="h-56 rounded-xl border border-border/60 bg-background/30 p-3">
                  <div className="space-y-2">
                    {session.messages.map((message, index) => (
                      <div key={`${message.timestamp}-history-${index}`} className="rounded-lg border border-border/60 bg-card/60 p-2">
                        <div className="mb-1 flex items-center justify-between text-[11px] text-muted-foreground">
                          <span>{message.role === "USER" ? "Usuário" : "Agente"}</span>
                          <Button
                            size="sm"
                            variant="ghost"
                            className="h-6 px-2"
                            onClick={() => {
                              navigator.clipboard.writeText(message.content);
                              toast.success("Mensagem copiada.");
                            }}
                          >
                            <Copy className="h-3 w-3" />
                          </Button>
                        </div>
                        <p className="line-clamp-2 text-xs leading-relaxed text-muted-foreground">{message.content}</p>
                      </div>
                    ))}
                  </div>
                </ScrollArea>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </PrivateRoute>
  );
}
