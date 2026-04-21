"use client";

import { useCallback, useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { CollaboratorRole, ProjectCollaborator } from "@/types/phase20";
import { Phase20Service } from "@/services/phase20";
import { toast } from "sonner";
import { UserPlus, UsersRound, X } from "lucide-react";

const collaboratorRoles: CollaboratorRole[] = ["OWNER", "EDITOR", "VIEWER"];

export function ProjectTeamPanel({ projectId }: { projectId: string }) {
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<CollaboratorRole>("EDITOR");
  const [collaborators, setCollaborators] = useState<ProjectCollaborator[]>([]);

  const loadCollaborators = useCallback(async () => {
    try {
      const data = await Phase20Service.listProjectCollaborators(projectId);
      setCollaborators(data);
    } catch (error) {
      console.error(error);
      toast.error("Falha ao carregar equipe do projeto.");
    }
  }, [projectId]);

  useEffect(() => {
    const timer = setTimeout(() => {
      void loadCollaborators();
    }, 0);
    return () => clearTimeout(timer);
  }, [loadCollaborators]);

  const addCollaborator = async () => {
    if (!email) return;
    try {
      await Phase20Service.addProjectCollaborator(projectId, email, role);
      setEmail("");
      await loadCollaborators();
      toast.success("Colaborador adicionado com sucesso.");
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível adicionar colaborador.");
    }
  };

  const updateCollaboratorRole = async (userId: string, nextRole: CollaboratorRole) => {
    try {
      await Phase20Service.updateProjectCollaboratorRole(projectId, userId, nextRole);
      setCollaborators((state) => state.map((c) => (c.user_id === userId ? { ...c, role: nextRole } : c)));
      toast.success("Permissão de colaborador atualizada.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao atualizar papel do colaborador.");
    }
  };

  const removeCollaborator = async (userId: string) => {
    try {
      await Phase20Service.removeProjectCollaborator(projectId, userId);
      setCollaborators((state) => state.filter((c) => c.user_id !== userId));
      toast.success("Colaborador removido.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao remover colaborador.");
    }
  };

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><UsersRound className="h-4 w-4 text-primary" /> Equipe do Projeto</CardTitle>
        <CardDescription>Gestão de múltiplos usuários por projeto com controle de OWNER, EDITOR e VIEWER.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="grid gap-2 rounded-xl border p-3 md:grid-cols-[1fr_160px_auto]">
          <Input value={email} placeholder="email@colaborador.com" onChange={(e) => setEmail(e.target.value)} />
          <Select value={role} onValueChange={(value) => setRole(value as CollaboratorRole)}>
            <SelectTrigger><SelectValue /></SelectTrigger>
            <SelectContent>
              {collaboratorRoles.map((item) => <SelectItem key={item} value={item}>{item}</SelectItem>)}
            </SelectContent>
          </Select>
          <Button onClick={addCollaborator}><UserPlus className="mr-2 h-4 w-4" /> Adicionar</Button>
        </div>

        <div className="space-y-2">
          {collaborators.map((collaborator) => (
            <div key={collaborator.user_id} className="flex flex-col gap-2 rounded-xl border p-3 md:flex-row md:items-center md:justify-between">
              <div>
                <p className="font-medium">{collaborator.name}</p>
                <p className="text-sm text-muted-foreground">{collaborator.email}</p>
              </div>
              <div className="flex items-center gap-2">
                <Select value={collaborator.role} onValueChange={(value) => updateCollaboratorRole(collaborator.user_id, value as CollaboratorRole)}>
                  <SelectTrigger className="w-[140px]"><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {collaboratorRoles.map((item) => <SelectItem key={item} value={item}>{item}</SelectItem>)}
                  </SelectContent>
                </Select>
                <Button size="icon" variant="outline" onClick={() => removeCollaborator(collaborator.user_id)}>
                  <X className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
