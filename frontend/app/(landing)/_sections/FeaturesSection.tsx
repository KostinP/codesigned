import { SectionTitle } from "@/components/shared/SectionTitle"
import { FeatureBox } from "@/components/shared/FeatureBox"

export function FeaturesSection() {
  return (
    <section className="bg-[var(--gray-bg)] py-16 px-4">
      <div className="max-w-7xl mx-auto flex flex-col gap-12">
        <SectionTitle
          title="Почему стоит учиться у нас?"
          description="Мы объединяем академическую строгость с технологиями и пониманием того как учатся взрослые."
          className="max-w-3xl"
        />

        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <FeatureBox icon="🎮" title="Геймификация с научным подходом" text="Бейджи, рейтинги и прогресс делают обучение осознанным." />
          <FeatureBox icon="🔬" title="Качественные курсы" text="Создано экспертами и основано на методиках." />
          <FeatureBox icon="👨‍🏫" title="Обратная связь" text="Автоматизация + поддержка преподавателя." />
        </div>
      </div>
    </section>
  )
}
