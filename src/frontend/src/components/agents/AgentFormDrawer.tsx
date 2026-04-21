import { useState, useEffect, useMemo } from "react";
import { useForm, useFieldArray, useWatch } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Agent, AgentProvider, AgentSkill } from "@/types/agent";
import { agentService } from "@/services/agentService";
import { Button } from "@/components/ui/button";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetFooter,
} from "@/components/ui/sheet";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Badge } from "@/components/ui/badge";
import { Plus, Trash, Play, Bot } from "lucide-react";
import { toast } from "sonner";
import { ScrollArea } from "@/components/ui/scroll-area";

const formSchema = z.object({
  name: z.string().min(2, "Nome deve ter pelo menos 2 caracteres"),
  description: z.string().min(10, "Descrição muito curta"),
  provider: z.enum(["OPENAI", "ANTHROPIC", "GOOGLE", "OLLAMA"]),
  model: z.string().min(2, "Modelo é obrigatório"),
  api_key_ref: z.string().optional(),
  system_prompts: z.array(
    z.object({ value: z.string().min(1, "Prompt não pode ser vazio") })
  ).min(1, "Adicione pelo menos um prompt"),
  skills: z.array(z.string()).min(1, "Selecione pelo menos uma skill"),
  enabled: z.boolean(),
});

const SKILLS: AgentSkill[] = [
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

const MODELS: Record<AgentProvider, string[]> = {
  OPENAI: ["gpt-4o", "gpt-4-turbo", "gpt-3.5-turbo"],
  ANTHROPIC: ["claude-3-opus", "claude-3-sonnet", "claude-3-haiku", "claude-3-5-sonnet"],
  GOOGLE: ["gemini-1.5-pro", "gemini-1.5-flash"],
  OLLAMA: ["llama3", "mistral", "gemma"],
};

export function AgentFormDrawer({
  open,
  onOpenChange,
  agent,
  onSuccess
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  agent: Agent | null;
  onSuccess: () => void;
}) {
  const [loading, setLoading] = useState(false);
  const [testingConfig, setTestingConfig] = useState(false);
  const [testResult, setTestResult] = useState<string | null>(null);
  const [customModelInput, setCustomModelInput] = useState("");
  const [customModelsByProvider, setCustomModelsByProvider] = useState<Record<AgentProvider, string[]>>({
    OPENAI: [],
    ANTHROPIC: [],
    GOOGLE: [],
    OLLAMA: [],
  });

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      name: "",
      description: "",
      provider: "OPENAI",
      model: "",
      api_key_ref: "",
      system_prompts: [{ value: "" }],
      skills: [],
      enabled: true,
    },
  });

  const { fields: promptFields, append: appendPrompt, remove: removePrompt } = useFieldArray({
    control: form.control,
    name: "system_prompts",
  });

  useEffect(() => {
    if (agent) {
      form.reset({
        name: agent.name,
        description: agent.description,
        provider: agent.provider,
        model: agent.model,
        api_key_ref: agent.api_key_ref || "",
        system_prompts: agent.system_prompts.map(p => ({ value: p })),
        skills: agent.skills,
        enabled: agent.enabled,
      });
    } else {
      form.reset({
        name: "",
        description: "",
        provider: "OPENAI",
        model: "gpt-4o",
        api_key_ref: "",
        system_prompts: [{ value: "Você é um especialista em IA..." }],
        skills: [],
        enabled: true,
      });
    }
  }, [agent, form]);

  const selectedProvider = useWatch({
    control: form.control,
    name: "provider",
  });
  const selectedModel = useWatch({
    control: form.control,
    name: "model",
  });

  const modelOptions = useMemo(() => (
    Array.from(new Set([
      ...MODELS[selectedProvider as AgentProvider],
      ...customModelsByProvider[selectedProvider as AgentProvider],
      ...(selectedModel ? [selectedModel] : []),
    ]))
  ), [selectedProvider, customModelsByProvider, selectedModel]);

  useEffect(() => {
    if (!agent) { // only auto-select on create
      form.setValue("model", modelOptions[0]);
    }
  }, [selectedProvider, form, agent, modelOptions]);

  const handleAddCustomModel = () => {
    const model = customModelInput.trim();
    if (!model) return;

    const provider = selectedProvider as AgentProvider;
    const alreadyExists = [...MODELS[provider], ...customModelsByProvider[provider]].includes(model);
    if (alreadyExists) {
      form.setValue("model", model, { shouldValidate: true });
      setCustomModelInput("");
      return;
    }

    setCustomModelsByProvider((current) => ({
      ...current,
      [provider]: [...current[provider], model],
    }));
    form.setValue("model", model, { shouldValidate: true });
    setCustomModelInput("");
  };

  const toggleSkill = (skill: AgentSkill) => {
    const current = form.getValues("skills");
    if (current.includes(skill)) {
      form.setValue("skills", current.filter(s => s !== skill), { shouldValidate: true });
    } else {
      form.setValue("skills", [...current, skill], { shouldValidate: true });
    }
  };

  const onSubmit = async (values: z.infer<typeof formSchema>) => {
    try {
      setLoading(true);
      const payload = {
        ...values,
        system_prompts: values.system_prompts.map(p => p.value),
        skills: values.skills as AgentSkill[]
      };

      if (agent) {
        await agentService.updateAgent(agent.id, payload);
        toast.success("Agente atualizado com sucesso!");
      } else {
        await agentService.createAgent(payload);
        toast.success("Agente criado com sucesso!");
      }
      onSuccess();
      onOpenChange(false);
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Erro ao salvar agente");
    } finally {
      setLoading(false);
    }
  };

  const handleTestConfig = async () => {
    try {
      const values = form.getValues();
      setTestingConfig(true);
      setTestResult(null);

      const payload = {
        ...values,
        system_prompts: values.system_prompts.map(p => p.value),
        skills: values.skills as AgentSkill[]
      };

      toast.info("Testando configuração do agente...");
      const res = await agentService.testConfiguration(payload);

      if (res.success) {
        setTestResult(res.response);
        toast.success("Teste concluído com sucesso!");
      } else {
        toast.error("Falha ao testar agente");
      }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } catch (error: any) {
      toast.error(error.response?.data?.error || "Erro ao testar configuração");
    } finally {
      setTestingConfig(false);
    }
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-xl w-full flex flex-col h-full overflow-hidden p-0">
        <SheetHeader className="px-6 py-4 border-b">
          <SheetTitle>{agent ? "Editar Agente" : "Novo Agente"}</SheetTitle>
          <SheetDescription>
            Configure os parâmetros de IA, habilidades e comportamento base deste agente.
          </SheetDescription>
        </SheetHeader>

        <ScrollArea className="flex-1 px-6 py-4">
          <Form {...form}>
            <form id="agent-form" onSubmit={form.handleSubmit(onSubmit)} className="space-y-6 pb-20">
              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="name"
                  render={({ field }) => (
                    <FormItem className="col-span-2">
                      <FormLabel>Nome do Agente</FormLabel>
                      <FormControl>
                        <Input placeholder="Ex: Bob, Arquiteto Sênior" {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="description"
                  render={({ field }) => (
                    <FormItem className="col-span-2">
                      <FormLabel>Descrição Curta</FormLabel>
                      <FormControl>
                        <Input placeholder="Especialista em design de sistemas distribuídos..." {...field} />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="provider"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Provider (LLM)</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Selecione..." />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="OPENAI">OpenAI</SelectItem>
                          <SelectItem value="ANTHROPIC">Anthropic</SelectItem>
                          <SelectItem value="GOOGLE">Google Gemini</SelectItem>
                          <SelectItem value="OLLAMA">Ollama (Local)</SelectItem>
                        </SelectContent>
                      </Select>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="model"
                  render={({ field }) => (
                    <FormItem>
                      <FormLabel>Modelo</FormLabel>
                      <Select onValueChange={field.onChange} value={field.value}>
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Selecione o modelo" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {modelOptions.map(model => (
                            <SelectItem key={model} value={model}>{model}</SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                      <div className="mt-2 flex gap-2">
                        <Input
                          placeholder="Cadastrar outro modelo (ex: gpt-4.1)"
                          value={customModelInput}
                          onChange={(e) => setCustomModelInput(e.target.value)}
                        />
                        <Button type="button" variant="outline" onClick={handleAddCustomModel}>
                          Adicionar
                        </Button>
                      </div>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="api_key_ref"
                  render={({ field }) => (
                    <FormItem className="col-span-2">
                      <FormLabel>Chave de API (Secret / API Key)</FormLabel>
                      <FormControl>
                        <Input
                          placeholder="sk-api-key..."
                          type="password"
                          {...field}
                          value={field.value ?? ""}
                        />
                      </FormControl>
                      <FormDescription>
                        A chave de API para autenticação junto ao provedor de IA (OpenAI, Anthropic, etc).
                      </FormDescription>
                      <FormMessage />
                    </FormItem>
                  )}
                />

                <FormField
                  control={form.control}
                  name="enabled"
                  render={({ field }) => (
                    <FormItem className="flex flex-row items-center justify-between rounded-lg border p-3 shadow-sm col-span-2">
                      <div className="space-y-0.5">
                        <FormLabel>Agente Ativo</FormLabel>
                        <FormDescription>
                          Apenas agentes ativos são sorteados no Modo Dinâmico.
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
              </div>

              <div>
                <FormLabel className="mb-2 block">Skills (Habilidades)</FormLabel>
                <div className="flex flex-wrap gap-2 p-3 border rounded-md bg-muted/20">
                  {SKILLS.map(skill => (
                    <Badge
                      key={skill}
                      variant={form.getValues("skills").includes(skill) ? "default" : "outline"}
                      className="cursor-pointer"
                      onClick={() => toggleSkill(skill)}
                    >
                      {skill}
                    </Badge>
                  ))}
                </div>
                {form.formState.errors.skills && (
                  <p className="text-sm font-medium text-destructive mt-1">
                    {form.formState.errors.skills.message}
                  </p>
                )}
              </div>

              <div className="space-y-4">
                <div className="flex justify-between items-center">
                  <FormLabel>System Prompts (Persona base)</FormLabel>
                  <Button type="button" variant="outline" size="sm" onClick={() => appendPrompt({ value: "" })}>
                    <Plus className="h-3 w-3 mr-1" /> Add Prompt
                  </Button>
                </div>
                {promptFields.map((field, index) => (
                  <FormField
                    key={field.id}
                    control={form.control}
                    name={`system_prompts.${index}.value`}
                    render={({ field }) => (
                      <FormItem>
                        <FormControl>
                          <div className="flex gap-2">
                            <Textarea
                              {...field}
                              className="min-h-[80px]"
                              placeholder={`Prompt ${index + 1}...`}
                            />
                            {promptFields.length > 1 && (
                              <Button
                                type="button"
                                variant="destructive"
                                size="icon"
                                className="shrink-0"
                                onClick={() => removePrompt(index)}
                              >
                                <Trash className="h-4 w-4" />
                              </Button>
                            )}
                          </div>
                        </FormControl>
                        <FormMessage />
                      </FormItem>
                    )}
                  />
                ))}
              </div>

              {testResult && (
                <div className="p-4 rounded-md bg-muted border font-mono text-sm">
                  <div className="font-semibold mb-2 flex items-center">
                    <Bot className="w-4 h-4 mr-2" /> Resposta do Agente:
                  </div>
                  <div className="whitespace-pre-wrap">{testResult}</div>
                </div>
              )}
            </form>
          </Form>
        </ScrollArea>

        <SheetFooter className="px-6 py-4 border-t mt-auto flex-col sm:flex-row gap-2">
          <Button
            type="button"
            variant="secondary"
            onClick={handleTestConfig}
            disabled={testingConfig}
            className="sm:mr-auto"
          >
            <Play className={`mr-2 h-4 w-4 ${testingConfig ? 'animate-pulse' : ''}`} />
            {testingConfig ? 'Testando...' : 'Testar Configuração'}
          </Button>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={loading}>
            Cancelar
          </Button>
          <Button type="submit" form="agent-form" disabled={loading}>
            {loading ? "Salvando..." : agent ? "Salvar Alterações" : "Criar Agente"}
          </Button>
        </SheetFooter>
      </SheetContent>
    </Sheet>
  );
}
