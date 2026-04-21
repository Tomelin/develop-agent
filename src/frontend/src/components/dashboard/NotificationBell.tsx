"use client";

import { useEffect, useMemo, useState } from "react";
import { Bell } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Phase8Service } from "@/services/phase8";
import { NotificationItem } from "@/types/phase8";
import { useRouter } from "next/navigation";

export function NotificationBell() {
  const [items, setItems] = useState<NotificationItem[]>([]);
  const router = useRouter();

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const data = await Phase8Service.getNotifications();
        setItems(data.filter((item) => !item.read));
      } catch (error) {
        console.error(error);
      }
    };

    fetchNotifications();
    const interval = setInterval(fetchNotifications, 15000);
    return () => clearInterval(interval);
  }, []);

  const unreadCount = useMemo(() => items.length, [items]);

  const handleNotificationClick = async (item: NotificationItem) => {
    try {
      await Phase8Service.markNotificationAsRead(item.id);
    } catch (error) {
      console.error(error);
    }
    setItems((prev) => prev.filter((entry) => entry.id !== item.id));
    router.push(`/projects/${item.project_id}?phase=${item.phase_number}`);
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Button variant="ghost" size="icon" className="relative">
          <Bell className="h-5 w-5" />
          {unreadCount > 0 && (
            <Badge className="absolute -right-1 -top-1 h-5 min-w-5 rounded-full px-1 text-[10px]">{unreadCount}</Badge>
          )}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-96">
        <DropdownMenuLabel>Notificações</DropdownMenuLabel>
        {items.map((item) => (
          <DropdownMenuItem key={item.id} className="cursor-pointer flex flex-col items-start" onClick={() => handleNotificationClick(item)}>
            <p className="text-sm">{item.message}</p>
            <p className="text-xs text-muted-foreground">Fase {item.phase_number} • {new Date(item.created_at).toLocaleString()}</p>
          </DropdownMenuItem>
        ))}
        {!items.length && (
          <div className="px-2 py-4 text-center text-xs text-muted-foreground">Sem notificações não lidas.</div>
        )}
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
