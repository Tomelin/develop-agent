"use client";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from "@/components/ui/dialog";
import { Badge } from "@/components/ui/badge";
import { CircleHelp, Dices, Scale, Wallet } from "lucide-react";

export function DynamicModeHelpDialog() {
  return (
    <Dialog>
      <DialogTrigger render={<Button size="icon" variant="ghost" className="h-8 w-8 rounded-full"><CircleHelp className="h-4 w-4" /></Button>} />
      <DialogContent className="max-w-3xl">
        <DialogHeader>
          <DialogTitle>Como funciona o Modo Dinâmico</DialogTitle>
        </DialogHeader>

        <div className="space-y-5 text-sm text-muted-foreground">
          <div className="p-4 rounded-xl border bg-card/60">
            <p className="font-medium text-foreground mb-2">Problema que resolvemos</p>
            <p>Um único modelo avaliando o próprio resultado gera viés. Regra prática: <span className="text-foreground font-semibold">o juiz não pode ser o mesmo que o réu</span>.</p>
          </div>

          <div className="grid md:grid-cols-3 gap-3">
            <div className="p-4 border rounded-xl bg-card/40">
              <Dices className="h-4 w-4 mb-2 text-primary" />
              <p className="font-medium text-foreground">Sorteio por Tríade</p>
              <p>Produtor, Revisor e Refinador são sorteados por execução.</p>
            </div>
            <div className="p-4 border rounded-xl bg-card/40">
              <Scale className="h-4 w-4 mb-2 text-emerald-400" />
              <p className="font-medium text-foreground">Diversidade máxima</p>
              <p>Quando possível, cada papel usa provider diferente.</p>
            </div>
            <div className="p-4 border rounded-xl bg-card/40">
              <Wallet className="h-4 w-4 mb-2 text-yellow-400" />
              <p className="font-medium text-foreground">Custo previsível</p>
              <p>Você acompanha histórico, score de diversidade e impacto financeiro.</p>
            </div>
          </div>

          <div className="p-4 rounded-xl border bg-background">
            <p className="font-medium text-foreground mb-3">Quando usar cada modo?</p>
            <div className="flex flex-wrap gap-2">
              <Badge variant="secondary">Dinâmico: fases críticas e revisão forte</Badge>
              <Badge variant="secondary">Fixo: padronização e previsibilidade</Badge>
              <Badge variant="secondary">Híbrido: melhor custo x qualidade</Badge>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
