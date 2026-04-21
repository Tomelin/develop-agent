import React from "react";

const formatInline = (text: string) => {
  const nodes: React.ReactNode[] = [];
  const regex = /(\*\*[^*]+\*\*|`[^`]+`|\[[^\]]+\]\([^)]+\))/g;
  let lastIndex = 0;
  let match: RegExpExecArray | null = regex.exec(text);

  while (match) {
    if (match.index > lastIndex) {
      nodes.push(text.slice(lastIndex, match.index));
    }

    const token = match[0];
    if (token.startsWith("**")) {
      nodes.push(
        <strong key={`${match.index}-b`} className="font-semibold text-foreground">
          {token.slice(2, -2)}
        </strong>,
      );
    } else if (token.startsWith("`")) {
      nodes.push(
        <code key={`${match.index}-c`} className="rounded bg-background/70 px-1.5 py-0.5 font-mono text-xs text-primary">
          {token.slice(1, -1)}
        </code>,
      );
    } else {
      const linkMatch = token.match(/\[([^\]]+)\]\(([^)]+)\)/);
      if (linkMatch) {
        nodes.push(
          <a key={`${match.index}-l`} className="text-primary underline underline-offset-4" href={linkMatch[2]} target="_blank" rel="noreferrer">
            {linkMatch[1]}
          </a>,
        );
      }
    }

    lastIndex = regex.lastIndex;
    match = regex.exec(text);
  }

  if (lastIndex < text.length) {
    nodes.push(text.slice(lastIndex));
  }

  return nodes.length > 0 ? nodes : text;
};

export function Markdown({ content }: { content: string }) {
  const lines = content.split("\n");

  return (
    <div className="space-y-2 text-sm leading-relaxed text-muted-foreground">
      {lines.map((line, index) => {
        const trimmed = line.trim();
        if (!trimmed) return <div key={`${index}-sp`} className="h-2" />;
        if (trimmed.startsWith("### ")) return <h3 key={index} className="text-base font-semibold text-foreground">{formatInline(trimmed.slice(4))}</h3>;
        if (trimmed.startsWith("## ")) return <h2 key={index} className="text-lg font-semibold text-foreground">{formatInline(trimmed.slice(3))}</h2>;
        if (trimmed.startsWith("# ")) return <h1 key={index} className="text-xl font-bold text-foreground">{formatInline(trimmed.slice(2))}</h1>;
        if (trimmed.match(/^[-*] /)) return <li key={index} className="ml-5 list-disc">{formatInline(trimmed.slice(2))}</li>;
        if (trimmed.match(/^\d+\. /)) return <li key={index} className="ml-5 list-decimal">{formatInline(trimmed.replace(/^\d+\. /, ""))}</li>;
        if (trimmed.startsWith("> ")) return <blockquote key={index} className="border-l-2 border-primary/40 pl-3 italic">{formatInline(trimmed.slice(2))}</blockquote>;
        return <p key={index}>{formatInline(trimmed)}</p>;
      })}
    </div>
  );
}
