// В app/(landing)/_sections/CTASection.tsx
"use client"

import { useEffect, useState } from "react"

export function CTASection() {
  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
  }, []);

  // Для SSR рендерим заглушку с теми же размерами
  if (!isClient) {
    return (
      <section className="py-20 px-4 bg-[var(--cta)] text-white">
        <div className="max-w-5xl mx-auto text-center">
          <div className="aspect-video w-full rounded-xl bg-gray-200 animate-pulse" />
        </div>
      </section>
    );
  }

  return (
    <section className="py-20 px-4 bg-[var(--cta)] text-white">
      <div className="max-w-5xl mx-auto text-center">
        <div className="relative aspect-video w-full rounded-xl overflow-hidden shadow-lg border border-white/10">
          <iframe
            src="https://www.youtube.com/embed/dQw4w9WgXcQ"
            title="Demo video"
            allowFullScreen
            className="w-full h-full"
          />
        </div>
      </div>
    </section>
  );
}