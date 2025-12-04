"use client"

import Image from "next/image"
import { Button } from "@/components/ui/button"
import { Fact } from "@/components/shared/Fact"
import { COLORS } from "@/lib/constants"

export function HeroSection() {
  return (
    <section className="px-4 py-12 md:py-24">
      <div className="max-w-7xl mx-auto grid grid-cols-1 md:grid-cols-2 gap-12 items-center">
        <div>
          <h1 className="text-3xl sm:text-4xl md:text-5xl font-bold mb-6 leading-tight">
            Современное образование —<br /> интерактивное, структурированное и доступное
          </h1>
          <p className="text-muted-foreground text-base md:text-lg mb-8 max-w-2xl">
            Осваивайте курсы в удобном формате, получайте автоматическую проверку тестов и персональные рекомендации.
          </p>

          <div className="flex flex-col sm:flex-row gap-4">
            <Button size="lg" className="rounded-full">Начать обучение</Button>
            <Button variant="outline" size="lg" className="rounded-full">Посмотреть курсы</Button>
          </div>

          <div className="mt-12 flex flex-wrap gap-8 text-sm text-muted-foreground">
            <Fact title="100+" text="Интерактивных уроков" colorVar={COLORS.yellowText} />
            <Fact title="15+" text="Качественных курсов" colorVar={COLORS.blueText} />
            <Fact title="60+" text="Обучено студентов" colorVar={COLORS.orangeText} />
          </div>
        </div>

        <div className="relative w-full h-[400px] md:h-[480px]">
          <Image src="/illustrations/hero-illustration.png" alt="Hero Image" fill className="object-contain" />
        </div>
      </div>
    </section>
  )
}
