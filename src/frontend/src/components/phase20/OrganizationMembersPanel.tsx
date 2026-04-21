"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import { Phase20Service } from "@/services/phase20";
import { OrganizationMember, OrganizationRole } from "@/types/phase20";
import { MailPlus, Shield, Trash2, Users } from "lucide-react";

const orgRoles: OrganizationRole[] = ["OWNER", "ADMIN", "MEMBER", "VIEWER"];

export function OrganizationMembersPanel() {
  const [members, setMembers] = useState<OrganizationMember[]>([]);
  const [inviteEmail, setInviteEmail] = useState("");
  const [inviteRole, setInviteRole] = useState<OrganizationRole>("MEMBER");
  const [loading, setLoading] = useState(true);

  const loadMembers = async () => {
    try {
      const data = await Phase20Service.listOrganizationMembers();
      setMembers(data);
    } catch (error) {
      console.error(error);
      toast.error("Não foi possível carregar os membros da organização.");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const timer = setTimeout(() => {
      void loadMembers();
    }, 0);
    return () => clearTimeout(timer);
  }, []);

  const inviteMember = async () => {
    if (!inviteEmail) return;
    try {
      await Phase20Service.inviteMember(inviteEmail, inviteRole);
      toast.success("Convite enviado com sucesso.");
      setInviteEmail("");
      await loadMembers();
    } catch (error) {
      console.error(error);
      toast.error("Falha ao enviar convite.");
    }
  };

  const updateRole = async (userId: string, role: OrganizationRole) => {
    try {
      await Phase20Service.updateOrganizationMemberRole(userId, role);
      setMembers((current) => current.map((member) => (member.user_id === userId ? { ...member, role } : member)));
      toast.success("Permissão atualizada.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao atualizar role.");
    }
  };

  const removeMember = async (userId: string) => {
    try {
      await Phase20Service.removeOrganizationMember(userId);
      setMembers((current) => current.filter((member) => member.user_id !== userId));
      toast.success("Membro removido da organização.");
    } catch (error) {
      console.error(error);
      toast.error("Falha ao remover membro.");
    }
  };

  return (
    <Card className="bg-card/50 border-border">
      <CardHeader>
        <CardTitle className="flex items-center gap-2"><Users className="h-4 w-4 text-primary" /> Gestão de membros da organização</CardTitle>
        <CardDescription>Controle convites e papéis por tenant com RBAC por organização.</CardDescription>
      </CardHeader>
      <CardContent className="space-y-5">
        <div className="rounded-xl border bg-background/40 p-3">
          <p className="mb-3 text-sm font-medium">Convidar novo membro</p>
          <div className="grid gap-2 md:grid-cols-[1fr_180px_auto]">
            <Input placeholder="email@empresa.com" value={inviteEmail} onChange={(e) => setInviteEmail(e.target.value)} />
            <Select value={inviteRole} onValueChange={(value) => setInviteRole(value as OrganizationRole)}>
              <SelectTrigger><SelectValue placeholder="Selecione o role" /></SelectTrigger>
              <SelectContent>
                {orgRoles.map((role) => <SelectItem key={role} value={role}>{role}</SelectItem>)}
              </SelectContent>
            </Select>
            <Button onClick={inviteMember}><MailPlus className="mr-2 h-4 w-4" /> Convidar</Button>
          </div>
        </div>

        <div className="space-y-2">
          {loading ? <p className="text-sm text-muted-foreground">Carregando membros...</p> : members.map((member) => (
            <div key={member.user_id} className="flex flex-col gap-3 rounded-xl border p-3 md:flex-row md:items-center md:justify-between">
              <div>
                <p className="font-medium">{member.name}</p>
                <p className="text-sm text-muted-foreground">{member.email}</p>
              </div>
              <div className="flex items-center gap-2">
                <Badge variant="outline" className="gap-1"><Shield className="h-3 w-3" /> {member.role}</Badge>
                <Select value={member.role} onValueChange={(value) => updateRole(member.user_id, value as OrganizationRole)}>
                  <SelectTrigger className="w-[140px]"><SelectValue /></SelectTrigger>
                  <SelectContent>
                    {orgRoles.map((role) => <SelectItem key={role} value={role}>{role}</SelectItem>)}
                  </SelectContent>
                </Select>
                <Button variant="outline" size="icon" onClick={() => removeMember(member.user_id)}>
                  <Trash2 className="h-4 w-4" />
                </Button>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
