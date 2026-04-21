import re

with open("src/frontend/src/app/projects/[id]/page.tsx", "r") as f:
    content = f.read()

# Replace the KanbanBoard import to remove it if unused, but we'll just ignore it or replace it
content = content.replace('import { KanbanBoard } from "@/components/projects/KanbanBoard";', '')

# Replace the KanbanBoard rendering with a CTA to the new Roadmap page
new_content = """
            <div className="flex flex-col items-center justify-center p-12 text-center bg-card/30 rounded-xl border border-border mt-4">
              <LayoutTemplate className="h-12 w-12 text-muted-foreground mb-4" />
              <h3 className="text-xl font-semibold mb-2">Roadmap & KanBan</h3>
              <p className="text-muted-foreground mb-6 max-w-md">
                Gerencie o progresso do desenvolvimento, visualize épicos, dependências e a timeline do projeto.
              </p>
              <Link href={`/projects/${project.id}/roadmap`}>
                <Button className="gap-2">
                  <LayoutTemplate className="h-4 w-4" />
                  Abrir Dashboard do Roadmap
                </Button>
              </Link>
            </div>
"""

content = re.sub(r'<KanbanBoard projectId=\{project\.id\} />', new_content, content)

with open("src/frontend/src/app/projects/[id]/page.tsx", "w") as f:
    f.write(content)
