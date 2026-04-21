"use client";

import { useEffect, useMemo, useState } from "react";
import {
  DndContext,
  PointerSensor,
  closestCenter,
  useSensor,
  useSensors,
  DragEndEvent,
} from "@dnd-kit/core";
import {
  SortableContext,
  arrayMove,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import {
  FileDown,
  FileUp,
  GripVertical,
  Lightbulb,
  Plus,
  Sparkles,
  Trash2,
  WandSparkles,
} from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import { Badge } from "@/components/ui/badge";
import { promptService } from "@/services/promptService";
import { PROMPT_GROUPS, PromptGroup, PromptTemplate, UserPrompt } from "@/types/prompt";
import Link from "next/link";

const MAX_PROMPT_LENGTH = 2000;

function SortablePromptCard({
  prompt,
  onEdit,
  onDelete,
  onToggle,
}: {
  prompt: UserPrompt;
  onEdit: (prompt: UserPrompt) => void;
  onDelete: (promptId: string) => void;
  onToggle: (prompt: UserPrompt, enabled: boolean) => void;
}) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({ id: prompt.id });

  return (
    <div
      ref={setNodeRef}
      style={{ transform: CSS.Transform.toString(transform), transition, opacity: isDragging ? 0.5 : 1 }}
      className="rounded-xl border border-border bg-card/60 p-4"
    >
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-start gap-3">
          <button className="mt-1 text-muted-foreground hover:text-foreground" {...attributes} {...listeners}>
            <GripVertical className="h-4 w-4" />
          </button>
          <div className="space-y-2">
            <div className="flex items-center gap-2">
              <h4 className="font-semibold">{prompt.title}</h4>
              <Badge variant="secondary">Prioridade {prompt.priority}</Badge>
            </div>
            <p className="text-sm text-muted-foreground line-clamp-2">{prompt.content}</p>
            <div className="flex flex-wrap gap-1">
              {prompt.tags?.map((tag) => (
                <Badge key={tag} variant="outline" className="text-xs">
                  #{tag}
                </Badge>
              ))}
            </div>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Switch checked={prompt.enabled} onCheckedChange={(checked) => onToggle(prompt, checked)} />
          <Button variant="outline" size="sm" onClick={() => onEdit(prompt)}>
            Editar
          </Button>
          <Button variant="ghost" size="icon" onClick={() => onDelete(prompt.id)}>
            <Trash2 className="h-4 w-4 text-destructive" />
          </Button>
        </div>
      </div>
    </div>
  );
}

export function PromptsManager() {
  const [activeGroup, setActiveGroup] = useState<PromptGroup>("GLOBAL");
  const [promptsByGroup, setPromptsByGroup] = useState<Record<PromptGroup, UserPrompt[]>>({} as Record<PromptGroup, UserPrompt[]>);
  const [templates, setTemplates] = useState<PromptTemplate[]>([]);
  const [loading, setLoading] = useState(true);

  const [editorOpen, setEditorOpen] = useState(false);
  const [templateOpen, setTemplateOpen] = useState(false);
  const [importOpen, setImportOpen] = useState(false);
  const [editingPrompt, setEditingPrompt] = useState<UserPrompt | null>(null);

  const [form, setForm] = useState({
    title: "",
    content: "",
    group: "GLOBAL" as PromptGroup,
    tags: "",
    enabled: true,
  });

  const [importMode, setImportMode] = useState<"MERGE" | "REPLACE">("MERGE");
  const [importPayload, setImportPayload] = useState("");

  const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 5 } }));

  const activePrompts = useMemo(() => promptsByGroup[activeGroup] || [], [promptsByGroup, activeGroup]);

  const fetchAll = async () => {
    setLoading(true);
    try {
      const [promptRes, templateRes] = await Promise.all([promptService.getPrompts(), promptService.getTemplates()]);
      const grouped = {} as Record<PromptGroup, UserPrompt[]>;
      PROMPT_GROUPS.forEach((group) => (grouped[group.value] = []));

      (promptRes.items || []).forEach((prompt) => {
        grouped[prompt.group]?.push(prompt);
      });

      PROMPT_GROUPS.forEach((group) => {
        grouped[group.value] = (grouped[group.value] || []).sort((a, b) => a.priority - b.priority);
      });

      setPromptsByGroup(grouped);
      setTemplates(templateRes);
    } catch {
      toast.error("Falha ao carregar prompts e templates.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect
    fetchAll();
  }, []);

  const openCreate = () => {
    setEditingPrompt(null);
    setForm({ title: "", content: "", group: activeGroup, tags: "", enabled: true });
    setEditorOpen(true);
  };

  const openEdit = (prompt: UserPrompt) => {
    setEditingPrompt(prompt);
    setForm({
      title: prompt.title,
      content: prompt.content,
      group: prompt.group,
      tags: (prompt.tags || []).join(", "),
      enabled: prompt.enabled,
    });
    setEditorOpen(true);
  };

  const savePrompt = async () => {
    if (!form.title.trim() || !form.content.trim()) {
      toast.error("Título e conteúdo são obrigatórios.");
      return;
    }

    if (form.content.length > MAX_PROMPT_LENGTH) {
      toast.error("Conteúdo excede o limite de 2000 caracteres.");
      return;
    }

    const tags = form.tags.split(",").map((tag) => tag.trim()).filter(Boolean);

    try {
      if (editingPrompt) {
        await promptService.updatePrompt(editingPrompt.id, {
          title: form.title,
          content: form.content,
          group: form.group,
          tags,
          enabled: form.enabled,
        });
        toast.success("Prompt atualizado com sucesso.");
      } else {
        await promptService.createPrompt({
          title: form.title,
          content: form.content,
          group: form.group,
          tags,
          enabled: form.enabled,
          priority: (promptsByGroup[form.group]?.length || 0) + 1,
        });
        toast.success("Prompt criado com sucesso.");
      }
      setEditorOpen(false);
      fetchAll();
    } catch {
      toast.error("Não foi possível salvar o prompt.");
    }
  };

  const deletePrompt = async (id: string) => {
    try {
      await promptService.deletePrompt(id);
      toast.success("Prompt removido.");
      fetchAll();
    } catch {
      toast.error("Falha ao remover prompt.");
    }
  };

  const togglePrompt = async (prompt: UserPrompt, enabled: boolean) => {
    try {
      await promptService.updatePrompt(prompt.id, { enabled });
      setPromptsByGroup((prev) => ({
        ...prev,
        [prompt.group]: prev[prompt.group].map((item) => (item.id === prompt.id ? { ...item, enabled } : item)),
      }));
      toast.success(enabled ? "Prompt habilitado." : "Prompt desabilitado.");
    } catch {
      toast.error("Não foi possível alterar o status do prompt.");
    }
  };

  const handleDragEnd = async (event: DragEndEvent) => {
    const { active, over } = event;
    if (!over || active.id === over.id) return;

    const oldIndex = activePrompts.findIndex((item) => item.id === active.id);
    const newIndex = activePrompts.findIndex((item) => item.id === over.id);
    const reordered = arrayMove(activePrompts, oldIndex, newIndex).map((prompt, index) => ({ ...prompt, priority: index + 1 }));

    setPromptsByGroup((prev) => ({ ...prev, [activeGroup]: reordered }));

    try {
      await promptService.reorderPrompts(reordered.map((item) => ({ id: item.id, priority: item.priority })));
      toast.success("Prioridade atualizada.");
    } catch {
      toast.error("Falha ao reordenar prompts.");
      fetchAll();
    }
  };

  const importFromTemplate = async (templateId: string) => {
    try {
      await promptService.createFromTemplate(templateId, activeGroup);
      toast.success("Template importado como prompt.");
      setTemplateOpen(false);
      fetchAll();
    } catch {
      toast.error("Falha ao importar template.");
    }
  };

  const exportPrompts = async () => {
    try {
      const blob = await promptService.exportPrompts();
      const url = window.URL.createObjectURL(blob);
      const link = document.createElement("a");
      link.href = url;
      link.download = `prompts-backup-${new Date().toISOString().slice(0, 10)}.json`;
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
      toast.success("Export concluído.");
    } catch {
      toast.error("Não foi possível exportar.");
    }
  };

  const importPrompts = async () => {
    try {
      const parsed = JSON.parse(importPayload);
      await promptService.importPrompts({ mode: importMode, prompts: parsed.prompts || parsed });
      toast.success("Importação concluída.");
      setImportOpen(false);
      setImportPayload("");
      fetchAll();
    } catch {
      toast.error("JSON inválido ou payload incompatível.");
    }
  };

  return (
    <div className="space-y-6">
      <Card className="border-primary/20 bg-gradient-to-br from-card to-card/60">
        <CardHeader>
          <CardTitle className="flex items-center gap-2 text-2xl">
            <Sparkles className="h-6 w-6 text-primary" /> Gestão de Prompts
          </CardTitle>
          <CardDescription>
            Configure regras persistentes para toda a esteira da agência. Esses prompts serão agregados automaticamente ao contexto dos agentes.
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-3">
          <Button onClick={openCreate} className="gap-2">
            <Plus className="h-4 w-4" /> Adicionar Prompt
          </Button>
          <Button variant="outline" onClick={() => setTemplateOpen(true)} className="gap-2">
            <WandSparkles className="h-4 w-4" /> Importar do Template
          </Button>
          <Button variant="outline" onClick={exportPrompts} className="gap-2">
            <FileDown className="h-4 w-4" /> Exportar
          </Button>
          <Button variant="outline" onClick={() => setImportOpen(true)} className="gap-2">
            <FileUp className="h-4 w-4" /> Importar
          </Button>
          <Button variant="secondary" asChild>
            <Link href="/prompts/preview">Abrir Preview de Composição</Link>
          </Button>
        </CardContent>
      </Card>

      <Tabs value={activeGroup} onValueChange={(v) => setActiveGroup(v as PromptGroup)}>
        <TabsList className="w-full h-auto p-1 grid grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-1">
          {PROMPT_GROUPS.map((group) => (
            <TabsTrigger key={group.value} value={group.value} className="text-xs md:text-sm">
              {group.label}
            </TabsTrigger>
          ))}
        </TabsList>

        {PROMPT_GROUPS.map((group) => (
          <TabsContent key={group.value} value={group.value} className="space-y-4">
            <Card className="border-border/70 bg-card/50">
              <CardHeader>
                <CardTitle>{group.label}</CardTitle>
                <CardDescription>{group.description}</CardDescription>
              </CardHeader>
              <CardContent>
                {loading ? (
                  <div className="h-24 animate-pulse rounded-lg bg-muted/40" />
                ) : (promptsByGroup[group.value] || []).length === 0 ? (
                  <div className="rounded-lg border border-dashed border-border p-8 text-center text-muted-foreground">
                    Nenhum prompt cadastrado neste grupo.
                  </div>
                ) : (
                  <DndContext sensors={sensors} collisionDetection={closestCenter} onDragEnd={handleDragEnd}>
                    <SortableContext items={(promptsByGroup[group.value] || []).map((item) => item.id)} strategy={verticalListSortingStrategy}>
                      <div className="space-y-3">
                        {(promptsByGroup[group.value] || []).map((prompt) => (
                          <SortablePromptCard
                            key={prompt.id}
                            prompt={prompt}
                            onEdit={openEdit}
                            onDelete={deletePrompt}
                            onToggle={togglePrompt}
                          />
                        ))}
                      </div>
                    </SortableContext>
                  </DndContext>
                )}
              </CardContent>
            </Card>
          </TabsContent>
        ))}
      </Tabs>

      <Dialog open={editorOpen} onOpenChange={setEditorOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>{editingPrompt ? "Editar Prompt" : "Novo Prompt"}</DialogTitle>
            <DialogDescription>Crie regras claras, testáveis e alinhadas ao padrão de qualidade da sua agência.</DialogDescription>
          </DialogHeader>

          <div className="grid gap-4 py-2">
            <div className="grid gap-2">
              <Label>Título</Label>
              <Input value={form.title} onChange={(e) => setForm((prev) => ({ ...prev, title: e.target.value }))} placeholder="Ex: Backend em Golang + Gin" />
            </div>

            <div className="grid gap-2">
              <Label>Grupo</Label>
              <Select value={form.group} onValueChange={(value) => setForm((prev) => ({ ...prev, group: value as PromptGroup }))}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {PROMPT_GROUPS.map((group) => (
                    <SelectItem key={group.value} value={group.value}>
                      {group.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="grid gap-2">
              <div className="flex items-center justify-between">
                <Label>Conteúdo</Label>
                <span className={`text-xs ${form.content.length > MAX_PROMPT_LENGTH ? "text-destructive" : "text-muted-foreground"}`}>
                  {form.content.length}/{MAX_PROMPT_LENGTH}
                </span>
              </div>
              <Textarea
                rows={8}
                value={form.content}
                onChange={(e) => setForm((prev) => ({ ...prev, content: e.target.value }))}
                placeholder="Defina instruções explícitas que devem ser aplicadas automaticamente."
              />
            </div>

            <div className="grid gap-2">
              <Label>Tags (separadas por vírgula)</Label>
              <Input value={form.tags} onChange={(e) => setForm((prev) => ({ ...prev, tags: e.target.value }))} placeholder="golang, arquitetura, clean-code" />
            </div>

            <div className="flex items-center justify-between rounded-lg border border-border p-3">
              <div>
                <p className="font-medium">Prompt habilitado</p>
                <p className="text-xs text-muted-foreground">Prompts desabilitados não entram na composição.</p>
              </div>
              <Switch checked={form.enabled} onCheckedChange={(checked) => setForm((prev) => ({ ...prev, enabled: checked }))} />
            </div>

            <div className="rounded-lg border border-primary/20 bg-primary/5 p-3 text-sm">
              <button type="button" className="font-medium flex items-center gap-2">
                <Lightbulb className="h-4 w-4 text-primary" /> Dicas de escrita efetiva
              </button>
              <ul className="mt-2 list-disc pl-5 text-muted-foreground space-y-1">
                <li>Use linguagem imperativa e critérios mensuráveis.</li>
                <li>Evite ambiguidade e termos subjetivos sem contexto.</li>
                <li>Prefira padrões corporativos reaproveitáveis (stack, design, segurança).</li>
              </ul>
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setEditorOpen(false)}>
              Cancelar
            </Button>
            <Button onClick={savePrompt}>Salvar Prompt</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={templateOpen} onOpenChange={setTemplateOpen}>
        <DialogContent className="max-w-3xl">
          <DialogHeader>
            <DialogTitle>Templates de Prompt</DialogTitle>
            <DialogDescription>Biblioteca pré-definida para acelerar setup com qualidade profissional.</DialogDescription>
          </DialogHeader>
          <div className="max-h-[60vh] overflow-y-auto space-y-3 pr-1">
            {templates.map((template) => (
              <div key={template.id} className="rounded-lg border border-border p-4">
                <div className="flex items-start justify-between gap-4">
                  <div>
                    <p className="font-semibold">{template.title}</p>
                    <p className="text-sm text-muted-foreground line-clamp-2">{template.description || template.content}</p>
                    <div className="mt-2 flex flex-wrap gap-2">
                      <Badge variant="outline">{template.group}</Badge>
                      {template.category ? <Badge variant="secondary">{template.category}</Badge> : null}
                    </div>
                  </div>
                  <Button size="sm" onClick={() => importFromTemplate(template.id)}>
                    Importar
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </DialogContent>
      </Dialog>

      <Dialog open={importOpen} onOpenChange={setImportOpen}>
        <DialogContent className="max-w-2xl">
          <DialogHeader>
            <DialogTitle>Importar configuração</DialogTitle>
            <DialogDescription>Cole o JSON exportado e escolha merge ou substituição total.</DialogDescription>
          </DialogHeader>
          <div className="space-y-4">
            <div className="grid gap-2">
              <Label>Modo</Label>
              <Select value={importMode} onValueChange={(v) => setImportMode(v as "MERGE" | "REPLACE")}>
                <SelectTrigger>
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="MERGE">Merge (mesclar)</SelectItem>
                  <SelectItem value="REPLACE">Replace (substituir tudo)</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid gap-2">
              <Label>JSON</Label>
              <Textarea rows={10} value={importPayload} onChange={(e) => setImportPayload(e.target.value)} placeholder='{"prompts": [{"title": "..."}]}' />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setImportOpen(false)}>
              Cancelar
            </Button>
            <Button onClick={importPrompts}>Confirmar Importação</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
