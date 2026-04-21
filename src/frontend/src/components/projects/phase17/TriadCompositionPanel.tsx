"use client";

import { useEffect, useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Phase17Service } from "@/services/phase17";
import { providerPalette, TriadSelection } from "@/types/phase17";
import { Dices, Lock, User } from "lucide-react";

export function TriadCompositionPanel({ projectId }: { projectId: string }) {
  const [selections, setSelections] = useState<TriadSelection[]>([]);

  useEffect(() => {
    Phase17Service.getTriadSelections(projectId).then(setSelections).catch(console.error);
  }, [projectId]);

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle className="text-lg">Composição da Tríade por fase</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {selections.length === 0 && <p className="text-sm text-muted-foreground">Ainda não há sorteios registrados para este projeto.</p>}

        {selections.map((item) => {
          const agents = [
            { key: "Produtor", data: item.producer },
            { key: "Revisor", data: item.reviewer },
            { key: "Refinador", data: item.refiner },
          ];

          return (
            <div key={`${item.phase_number}-${item.selection_timestamp}`} className="rounded-xl border bg-background/70 p-4 space-y-3">
              <div className="flex items-center justify-between">
                <div>
                  <p className="font-medium">{item.phase_name}</p>
                  <p className="text-xs text-muted-foreground">{new Date(item.selection_timestamp).toLocaleString()}</p>
                </div>
                <Badge variant="outline" className="gap-1">
                  {item.mode === "DYNAMIC" ? <Dices className="h-3.5 w-3.5" /> : <Lock className="h-3.5 w-3.5" />}
                  {item.mode === "DYNAMIC" ? "Sorteado dinamicamente" : "Configuração fixa"}
                </Badge>
              </div>

              <div className="grid md:grid-cols-3 gap-3">
                {agents.map(({ key, data }) => {
                  const palette = providerPalette[data.provider];
                  return (
                    <div key={key} className={`p-3 rounded-lg border ${palette.bg}`}>
                      <div className="flex items-center gap-2 mb-2">
                        <Avatar className={`h-8 w-8 ring-1 ${palette.ring}`}>
                          <AvatarFallback>
                            <User className="h-4 w-4" />
                          </AvatarFallback>
                        </Avatar>
                        <div>
                          <p className="text-xs text-muted-foreground">{key}</p>
                          <p className="text-sm font-semibold leading-none">{data.name}</p>
                        </div>
                      </div>
                      <p className={`text-xs font-medium ${palette.text}`}>{data.provider}</p>
                      <p className="text-xs text-muted-foreground truncate">{data.model}</p>
                    </div>
                  );
                })}
              </div>
            </div>
          );
        })}
      </CardContent>
    </Card>
  );
}
