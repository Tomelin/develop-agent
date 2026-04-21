"use client";

import { useEffect, useState } from "react";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { CheckCircle2 } from "lucide-react";
import { PricingPlan } from "@/types/phase20";
import { Phase20Service } from "@/services/phase20";
import { toast } from "sonner";

export default function PricingPage() {
  const [plans, setPlans] = useState<PricingPlan[]>([]);

  useEffect(() => {
    const loadPlans = async () => {
      try {
        const data = await Phase20Service.getPricingPlans();
        setPlans(data);
      } catch (error) {
        console.error(error);
        toast.error("Não foi possível carregar planos de assinatura.");
      }
    };

    loadPlans();
  }, []);

  const checkout = async (planCode: string) => {
    try {
      const { checkout_url } = await Phase20Service.createStripeCheckout(planCode);
      window.location.assign(checkout_url);
    } catch (error) {
      console.error(error);
      toast.error("Falha ao iniciar checkout.");
    }
  };

  return (
    <div className="container py-10 max-w-7xl space-y-6">
      <div className="text-center space-y-2">
        <h1 className="text-4xl font-bold tracking-tight">Planos e Preços</h1>
        <p className="text-muted-foreground max-w-2xl mx-auto">Precificação transparente para escalar do MVP ao nível enterprise com multi-tenancy e integrações.</p>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        {plans.map((plan) => (
          <Card key={plan.code} className={`bg-card/50 border-border ${plan.highlighted ? "border-primary shadow-lg shadow-primary/10" : ""}`}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>{plan.name}</CardTitle>
                {plan.highlighted && <Badge>Mais escolhido</Badge>}
              </div>
              <CardDescription>{plan.description}</CardDescription>
              <p className="text-3xl font-bold mt-3">${plan.monthly_price_usd}<span className="text-sm text-muted-foreground font-normal">/mês</span></p>
            </CardHeader>
            <CardContent className="space-y-2">
              {plan.features.map((feature) => (
                <div key={feature.label} className="text-sm flex items-center gap-2">
                  <CheckCircle2 className={`h-4 w-4 ${feature.included ? "text-primary" : "text-muted-foreground"}`} />
                  <span className={feature.included ? "text-foreground" : "text-muted-foreground line-through"}>
                    {feature.label}{feature.value ? `: ${feature.value}` : ""}
                  </span>
                </div>
              ))}
            </CardContent>
            <CardFooter>
              <Button className="w-full" variant={plan.highlighted ? "default" : "outline"} onClick={() => checkout(plan.code)}>
                {plan.cta_label}
              </Button>
            </CardFooter>
          </Card>
        ))}
      </div>
    </div>
  );
}
