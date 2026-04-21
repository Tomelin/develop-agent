"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Phase20Service } from "@/services/phase20";
import { PublicTemplate } from "@/types/phase20";
import { toast } from "sonner";
import { Search, Sparkles, Star } from "lucide-react";

export default function MarketplacePage() {
  const [templates, setTemplates] = useState<PublicTemplate[]>([]);
  const [search, setSearch] = useState("");
  const [category, setCategory] = useState("ALL");
  const [group, setGroup] = useState("ALL");

  const loadTemplates = useCallback(async () => {
    try {
      const response = await Phase20Service.listPublicTemplates({
        search: search || undefined,
        category: category !== "ALL" ? category : undefined,
        group: group !== "ALL" ? group : undefined,
      });
      setTemplates(response.items);
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar o marketplace.");
    }
  }, [category, group, search]);

  useEffect(() => {
    const timer = setTimeout(() => {
      void loadTemplates();
    }, 0);
    return () => clearTimeout(timer);
  }, [loadTemplates]);

  const categories = useMemo(() => Array.from(new Set(templates.map((item) => item.category))), [templates]);
  const groups = useMemo(() => Array.from(new Set(templates.map((item) => item.group))), [templates]);

  const applyTemplate = async (templateId: string) => {
    try {
      await Phase20Service.activateTemplate(templateId);
      toast.success("Template adicionado à sua base de prompts.");
      await loadTemplates();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao usar template.");
    }
  };

  const starTemplate = async (templateId: string) => {
    try {
      await Phase20Service.starTemplate(templateId);
      toast.success("Template favoritado.");
      await loadTemplates();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao favoritar template.");
    }
  };

  return (
    <div className="container max-w-7xl py-10 space-y-6">
      <div className="space-y-2">
        <h1 className="text-3xl font-bold tracking-tight">Marketplace de Templates</h1>
        <p className="text-muted-foreground">Descubra, reutilize e publique templates de prompts para acelerar cada fase.</p>
      </div>

      <div className="rounded-xl border bg-card/40 p-3 grid gap-2 md:grid-cols-[1fr_180px_180px_auto]">
        <div className="relative">
          <Search className="h-4 w-4 text-muted-foreground absolute left-3 top-3" />
          <Input className="pl-9" placeholder="Buscar por título, tag ou descrição..." value={search} onChange={(e) => setSearch(e.target.value)} />
        </div>
        <Select value={category} onValueChange={(value) => setCategory(value ?? "ALL")}>
          <SelectTrigger><SelectValue placeholder="Categoria" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="ALL">Todas categorias</SelectItem>
            {categories.map((item) => <SelectItem key={item} value={item}>{item}</SelectItem>)}
          </SelectContent>
        </Select>
        <Select value={group} onValueChange={(value) => setGroup(value ?? "ALL")}>
          <SelectTrigger><SelectValue placeholder="Fase alvo" /></SelectTrigger>
          <SelectContent>
            <SelectItem value="ALL">Todas as fases</SelectItem>
            {groups.map((item) => <SelectItem key={item} value={item}>{item}</SelectItem>)}
          </SelectContent>
        </Select>
        <Button onClick={loadTemplates}>Filtrar</Button>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {templates.map((template) => (
          <Card key={template.id} className="bg-card/60 border-border hover:border-primary/50 transition-colors">
            <CardHeader>
              <div className="flex items-start justify-between gap-2">
                <div>
                  <CardTitle className="text-lg">{template.title}</CardTitle>
                  <CardDescription className="mt-1">{template.description}</CardDescription>
                </div>
                <Badge>{template.group}</Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-3">
              <div className="flex flex-wrap gap-1">
                {template.tags?.map((tag) => <Badge variant="outline" key={tag}>{tag}</Badge>)}
              </div>
              <div className="flex items-center justify-between text-sm text-muted-foreground">
                <span>{template.usage_count} usos</span>
                <span className="flex items-center gap-1"><Star className="h-3 w-3" /> {template.stars}</span>
              </div>
              <div className="flex flex-wrap gap-2">
                <Dialog>
                  <DialogTrigger render={<Button variant="outline" />}>Preview</DialogTrigger>
                  <DialogContent className="max-w-2xl">
                    <DialogHeader>
                      <DialogTitle>{template.title}</DialogTitle>
                      <DialogDescription>Template da fase {template.group}</DialogDescription>
                    </DialogHeader>
                    <pre className="max-h-[55vh] overflow-auto rounded-lg border p-3 text-xs whitespace-pre-wrap">{template.content}</pre>
                  </DialogContent>
                </Dialog>
                <Button onClick={() => applyTemplate(template.id)}><Sparkles className="mr-2 h-4 w-4" /> Usar template</Button>
                <Button variant="ghost" onClick={() => starTemplate(template.id)}><Star className="mr-2 h-4 w-4" /> Favoritar</Button>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
