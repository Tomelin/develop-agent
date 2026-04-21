export type PlatformStatusLevel = "OPERATIONAL" | "DEGRADED" | "OUTAGE";
export type ComponentHealth = "healthy" | "degraded" | "unhealthy" | "unknown";

export interface PlatformComponentStatus {
  status: ComponentHealth;
  latency_ms?: number;
  error?: string;
  uptime_30d_percent?: number;
}

export interface ProviderStatus {
  [provider: string]: string;
}

export interface StatusIncident {
  title: string;
  occurred_at: string;
  duration: string;
  resolution: string;
}

export interface PlatformStatusResponse {
  status: PlatformStatusLevel;
  checked_at: string;
  components: Record<string, PlatformComponentStatus>;
  providers: ProviderStatus;
  incidents: StatusIncident[];
}
