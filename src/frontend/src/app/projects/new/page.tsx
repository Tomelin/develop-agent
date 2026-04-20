"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { PrivateRoute } from "@/components/auth/PrivateRoute";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { ArrowLeft, ArrowRight, LayoutTemplate, Megaphone, MonitorPlay, Sparkles } from "lucide-react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { ProjectService } from "@/services/project";
import { Project, FlowType } from "@/types/project";
import { toast } from "sonner";
import Link from "next/link";

const projectSchema = z.object({
  flow_type: z.enum(["A", "B", "C"], { invalid_type_error: "Selecione o tipo de fluxo", required_error: "Selecione o tipo de fluxo" } as any),
  name: z.string().min(3, "O nome deve ter no mínimo 3 caracteres").max(100),
  description: z.string().min(10, "Forneça uma descrição detalhada (mínimo 10 caracteres)"),
  dynamic_mode: z.boolean(),
  linked_project_id: z.string().optional(),
});

type ProjectFormValues = z.infer<typeof projectSchema>;

const STEPS = [
  { id: "flow", title: "Tipo de Fluxo" },
  { id: "details", title: "Configurações Base" },
  { id: "review", title: "Confirmação" },
];

export default function NewProjectPage() {
  const router = useRouter();
  const [currentStep, setCurrentStep] = useState(0);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [availableProjects, setAvailableProjects] = useState<Project[]>([]);

  const form = useForm<ProjectFormValues>({
    resolver: zodResolver(projectSchema),
    defaultValues: {
      flow_type: "A",
      name: "",
      description: "",
      dynamic_mode: false,
    },
    mode: "onChange",
  });

  const watchFlowType = form.watch("flow_type");

  useEffect(() => {
    if (watchFlowType === "B" || watchFlowType === "C") {
      // Fetch available projects to link
      ProjectService.getProjects(1, 100).then((res) => {
        // Only allow linking to projects that are well advanced (e.g. COMPLETED or IN_PROGRESS at later phases)
        // For simplicity we show all, but ideally we filter by status and phase
        setAvailableProjects(res.items.filter(p => p.flow_type === "A"));
      }).catch(console.error);
    }
  }, [watchFlowType]);

  const handleNext = async () => {
    let isValid = false;

    if (currentStep === 0) {
      isValid = await form.trigger("flow_type");
    } else if (currentStep === 1) {
      isValid = await form.trigger(["name", "description", "linked_project_id"]);
    }

    if (isValid) {
      setCurrentStep((prev) => Math.min(prev + 1, STEPS.length - 1));
    }
  };

  const handleBack = () => {
    setCurrentStep((prev) => Math.max(prev - 1, 0));
  };

  const onSubmit = async (data: ProjectFormValues) => {
    if (currentStep !== 2) return;

    setIsSubmitting(true);
    try {
      const createdProject = await ProjectService.createProject({
        name: data.name,
        description: data.description,
        flow_type: data.flow_type,
        dynamic_mode: data.dynamic_mode,
        linked_project_id: data.linked_project_id && data.linked_project_id !== "none" ? data.linked_project_id : undefined,
      });

      toast.success("Projeto criado com sucesso!");
      router.push(`/projects/${createdProject.id}`);
    } catch (error) {
      console.error(error);
      toast.error("Erro ao criar projeto. Verifique os dados e tente novamente.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const getFlowDetails = (type: FlowType) => {
    switch(type) {
      case "A": return { icon: MonitorPlay, title: "Desenvolvimento de Software", color: "text-primary", bg: "bg-primary/10 border-primary/50" };
      case "B": return { icon: LayoutTemplate, title: "Landing Page Dinâmica", color: "text-secondary", bg: "bg-secondary/10 border-secondary/50" };
      case "C": return { icon: Megaphone, title: "Estratégia de Marketing", color: "text-chart-3", bg: "bg-chart-3/10 border-chart-3/50" };
    }
  };

  return (
    <PrivateRoute>
      <div className="container max-w-4xl py-10 px-4 md:px-6 mx-auto animate-in fade-in duration-500">
        <div className="mb-8">
          <Button variant="ghost"  className="mb-4 -ml-4 text-muted-foreground hover:text-foreground">
            <Link href="/dashboard"><ArrowLeft className="mr-2 h-4 w-4" /> Voltar ao Dashboard</Link>
          </Button>
          <h1 className="text-3xl font-bold tracking-tight">Criar Novo Projeto</h1>
          <p className="text-muted-foreground">Configure os parâmetros base para o agente IA iniciar o desenvolvimento.</p>
        </div>

        {/* Stepper Progress */}
        <div className="mb-8 relative">
          <div className="absolute top-1/2 left-0 w-full h-0.5 bg-border -z-10 -translate-y-1/2"></div>
          <div className="flex justify-between items-center z-10 relative">
            {STEPS.map((step, index) => (
              <div key={step.id} className="flex flex-col items-center">
                <div
                  className={`w-10 h-10 rounded-full flex items-center justify-center font-medium border-2 transition-colors
                    ${index < currentStep ? "bg-primary border-primary text-primary-foreground" :
                      index === currentStep ? "bg-background border-primary text-primary" :
                      "bg-background border-border text-muted-foreground"}`}
                >
                  {index < currentStep ? <Sparkles className="h-4 w-4" /> : index + 1}
                </div>
                <span className={`text-xs mt-2 font-medium ${index <= currentStep ? "text-foreground" : "text-muted-foreground"}`}>
                  {step.title}
                </span>
              </div>
            ))}
          </div>
        </div>

        <Card className="bg-card/50 backdrop-blur-sm border-border">
          <Form {...form}>
            <form onSubmit={form.handleSubmit(onSubmit)}>
              <CardHeader>
                <CardTitle>{STEPS[currentStep].title}</CardTitle>
                <CardDescription>
                  {currentStep === 0 && "Escolha a jornada que a IA irá percorrer."}
                  {currentStep === 1 && "Defina os detalhes e o contexto base do projeto."}
                  {currentStep === 2 && "Revise as informações antes de iniciar o planejamento."}
                </CardDescription>
              </CardHeader>

              <CardContent className="min-h-[300px]">
                {/* Step 1: Flow Type */}
                <div className={currentStep === 0 ? "block" : "hidden"}>
                  <FormField
                    control={form.control}
                    name="flow_type"
                    render={({ field }) => (
                      <FormItem className="space-y-4">
                        <div className="grid gap-4 md:grid-cols-3">
                          {[
                            { id: "A", title: "Software (SaaS/App)", desc: "Fluxo completo de engenharia, arquitetura, desenvolvimento frontend/backend, testes e DevOps.", icon: MonitorPlay },
                            { id: "B", title: "Landing Page", desc: "Fluxo acelerado focado em conversão, design visual e copywriting persuasivo.", icon: LayoutTemplate },
                            { id: "C", title: "Marketing Digital", desc: "Estratégia de go-to-market, SEO, análise de concorrentes e funis de venda.", icon: Megaphone },
                          ].map((flow) => {
                            const Icon = flow.icon;
                            return (
                              <label
                                key={flow.id}
                                className={`
                                  cursor-pointer rounded-xl border-2 p-6 transition-all hover:bg-card relative overflow-hidden
                                  ${field.value === flow.id ? "border-primary bg-primary/5" : "border-border bg-background"}
                                `}
                              >
                                <input
                                  type="radio"
                                  className="sr-only"
                                  value={flow.id}
                                  checked={field.value === flow.id}
                                  onChange={() => field.onChange(flow.id)}
                                />
                                {field.value === flow.id && (
                                  <div className="absolute top-3 right-3 text-primary">
                                    <Sparkles className="h-5 w-5" />
                                  </div>
                                )}
                                <Icon className={`h-8 w-8 mb-4 ${field.value === flow.id ? "text-primary" : "text-muted-foreground"}`} />
                                <h3 className="font-semibold mb-2">{flow.title}</h3>
                                <p className="text-sm text-muted-foreground">{flow.desc}</p>
                              </label>
                            );
                          })}
                        </div>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                </div>

                {/* Step 2: Configs */}
                <div className={currentStep === 1 ? "block space-y-6" : "hidden"}>
                  <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Nome do Projeto</FormLabel>
                        <FormControl>
                          <Input placeholder="Ex: Plataforma de Gestão Acme" {...field} className="bg-background" />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="description"
                    render={({ field }) => (
                      <FormItem>
                        <FormLabel>Visão / Descrição Inicial</FormLabel>
                        <FormControl>
                          <textarea
                            className="flex min-h-[120px] w-full rounded-md border border-input bg-background px-3 py-2 text-base ring-offset-background placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm"
                            placeholder="Descreva o problema que deseja resolver, público-alvo ou principal objetivo do produto..."
                            {...field}
                          />
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />

                  <FormField
                    control={form.control}
                    name="dynamic_mode"
                    render={({ field }) => (
                      <FormItem className="flex flex-row items-center justify-between rounded-lg border border-border p-4 bg-background/50">
                        <div className="space-y-0.5">
                          <FormLabel className="text-base">Modo Dinâmico (Revisão Contínua)</FormLabel>
                          <FormDescription>
                            Se ativado, a IA pedirá sua aprovação ao final de cada fase. Se desativado, rodará de forma autônoma (Modo Waterfall).
                          </FormDescription>
                        </div>
                        <FormControl>
                          <Switch
                            checked={field.value}
                            onCheckedChange={field.onChange}
                          />
                        </FormControl>
                      </FormItem>
                    )}
                  />

                  {(watchFlowType === "B" || watchFlowType === "C") && (
                    <FormField
                      control={form.control}
                      name="linked_project_id"
                      render={({ field }) => (
                        <FormItem>
                          <FormLabel>Projeto Base (Herança de Contexto)</FormLabel>
                          <Select onValueChange={field.onChange} defaultValue={field.value}>
                            <FormControl>
                              <SelectTrigger className="bg-background">
                                <SelectValue placeholder="Selecione um projeto de software para herdar o branding..." />
                              </SelectTrigger>
                            </FormControl>
                            <SelectContent>
                              <SelectItem value="none">Não vincular projeto</SelectItem>
                              {availableProjects.map(p => (
                                <SelectItem key={p.id} value={p.id}>{p.name}</SelectItem>
                              ))}
                            </SelectContent>
                          </Select>
                          <FormDescription>
                            O sistema extrairá a paleta de cores, tecnologias e a identidade visual do projeto selecionado.
                          </FormDescription>
                          <FormMessage />
                        </FormItem>
                      )}
                    />
                  )}
                </div>

                {/* Step 3: Review */}
                <div className={currentStep === 2 ? "block space-y-6" : "hidden"}>
                  <div className="rounded-xl border border-border bg-background p-6">
                    <div className="flex items-center gap-4 mb-6 pb-6 border-b border-border/50">
                      {(() => {
                        const flow = getFlowDetails(watchFlowType);
                        const Icon = flow?.icon;
                        return (
                          <>
                            <div className={`p-3 rounded-lg ${flow?.bg}`}>
                              {Icon && <Icon className={`h-6 w-6 ${flow.color}`} />}
                            </div>
                            <div>
                              <h3 className="font-semibold">{form.getValues("name") || "Sem Nome"}</h3>
                              <p className="text-sm text-muted-foreground">{flow?.title}</p>
                            </div>
                          </>
                        )
                      })()}
                    </div>

                    <div className="space-y-4">
                      <div>
                        <h4 className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">Descrição</h4>
                        <p className="text-sm">{form.getValues("description") || "Sem descrição."}</p>
                      </div>

                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <h4 className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">Modo de Execução</h4>
                          <p className="text-sm flex items-center gap-2">
                            {form.getValues("dynamic_mode") ? (
                              <><Sparkles className="h-4 w-4 text-primary" /> Dinâmico (Requer Feedback)</>
                            ) : (
                              <><MonitorPlay className="h-4 w-4 text-muted-foreground" /> Autônomo (Waterfall)</>
                            )}
                          </p>
                        </div>

                        {(watchFlowType === "B" || watchFlowType === "C") && form.getValues("linked_project_id") && form.getValues("linked_project_id") !== "none" && (
                          <div>
                            <h4 className="text-xs font-medium text-muted-foreground uppercase tracking-wider mb-1">Projeto Vinculado</h4>
                            <p className="text-sm">{availableProjects.find(p => p.id === form.getValues("linked_project_id"))?.name}</p>
                          </div>
                        )}
                      </div>
                    </div>
                  </div>

                  <div className="bg-primary/10 border border-primary/20 rounded-lg p-4 text-sm text-primary flex items-start gap-3">
                    <Sparkles className="h-5 w-5 mt-0.5 shrink-0" />
                    <p>Ao criar, os Agentes começarão imediatamente a analisar o contexto fornecido e iniciarão o Planejamento (Fase 1).</p>
                  </div>
                </div>
              </CardContent>

              <CardFooter className="flex justify-between border-t border-border pt-6">
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleBack}
                  disabled={currentStep === 0 || isSubmitting}
                >
                  Voltar
                </Button>

                {currentStep < 2 ? (
                  <Button type="button" onClick={handleNext}>
                    Avançar <ArrowRight className="ml-2 h-4 w-4" />
                  </Button>
                ) : (
                  <Button type="submit" disabled={isSubmitting} className="bg-primary text-primary-foreground">
                    {isSubmitting ? (
                      <span className="flex items-center gap-2">
                        <span className="animate-spin h-4 w-4 border-2 border-current border-t-transparent rounded-full" />
                        Criando...
                      </span>
                    ) : (
                      "Criar Projeto"
                    )}
                  </Button>
                )}
              </CardFooter>
            </form>
          </Form>
        </Card>
      </div>
    </PrivateRoute>
  );
}