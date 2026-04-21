"use client";

import { OrganizationMembersPanel } from "@/components/phase20/OrganizationMembersPanel";

export default function OrganizationPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold tracking-tight">Organização</h1>
        <p className="text-muted-foreground mt-1">Gerencie membros, convites e permissões por tenant.</p>
      </div>
      <OrganizationMembersPanel />
    </div>
  );
}
