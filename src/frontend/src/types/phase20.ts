export type OrganizationRole = "OWNER" | "ADMIN" | "MEMBER" | "VIEWER";
export type CollaboratorRole = "OWNER" | "EDITOR" | "VIEWER";

export interface OrganizationMember {
  user_id: string;
  name: string;
  email: string;
  role: OrganizationRole;
  joined_at: string;
}

export interface ProjectCollaborator {
  user_id: string;
  name: string;
  email: string;
  role: CollaboratorRole;
  added_at: string;
}

export interface PublicTemplate {
  id: string;
  title: string;
  description: string;
  category: string;
  content: string;
  group: string;
  stars: number;
  usage_count: number;
  creator_id: string;
  creator_name?: string;
  visibility: "PUBLIC" | "PRIVATE";
  tags: string[];
}

export interface PublicTemplateListResponse {
  items: PublicTemplate[];
  total: number;
  page: number;
  size: number;
}

export interface IntegrationConnectionState {
  provider: "github" | "jira" | "slack";
  connected: boolean;
  last_sync_at?: string;
  configured_by?: string;
}

export interface PricingPlanFeature {
  label: string;
  included: boolean;
  value?: string;
}

export interface PricingPlan {
  code: "FREE" | "STARTER" | "PRO" | "ENTERPRISE";
  name: string;
  monthly_price_usd: number;
  annual_price_usd?: number;
  highlighted?: boolean;
  description: string;
  cta_label: string;
  stripe_price_id?: string;
  features: PricingPlanFeature[];
}

export interface PublicRoadmapFeature {
  id: string;
  title: string;
  description: string;
  milestone: "v1.0" | "v1.5" | "v2.0" | string;
  status: "PLANNED" | "IN_DEVELOPMENT" | "COMPLETED";
  votes: number;
  target_quarter?: string;
}

export interface PublicRoadmapChangelog {
  version: string;
  date: string;
  highlights: string[];
}

export interface PublicRoadmapResponse {
  vision: string;
  features: PublicRoadmapFeature[];
  changelog: PublicRoadmapChangelog[];
}
