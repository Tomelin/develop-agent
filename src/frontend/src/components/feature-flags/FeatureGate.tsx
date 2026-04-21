"use client";

import { ReactNode, useEffect, useState } from "react";
import { Phase17Service } from "@/services/phase17";

interface FeatureGateProps {
  flag: string;
  fallback?: ReactNode;
  children: ReactNode;
}

export function FeatureGate({ flag, fallback = null, children }: FeatureGateProps) {
  const [enabled, setEnabled] = useState<boolean | null>(null);

  useEffect(() => {
    const load = async () => {
      try {
        const flags = await Phase17Service.getFeatureFlagsPublic();
        const selected = flags.find((f) => f.key === flag);
        setEnabled(Boolean(selected?.enabled));
      } catch (error) {
        console.error("Unable to read feature flag", error);
        setEnabled(false);
      }
    };
    load();
  }, [flag]);

  if (enabled === null) return null;
  if (!enabled) return <>{fallback}</>;

  return <>{children}</>;
}
