import { courseCategories } from "@/lib/mocks/courseCategories"
import { CourseCard } from "@/components/cards/CourseCard"
import { Button } from "@/components/ui/button"

export function CategoriesSection() {
  return (
    <section className="py-16 px-4">
      <div className="max-w-7xl mx-auto flex flex-col gap-12 items-center text-center">
        <div className="max-w-3xl">
          <h2 className="text-2xl md:text-3xl font-bold mb-4">Основные направления обучения</h2>
          <p className="text-muted-foreground text-base">Выберите платный или бесплатный курс</p>
        </div>

        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6 w-full">
          {courseCategories.map((cat, i) => (
            <CourseCard key={i} {...cat} />
          ))}
        </div>

        <Button variant="outline" size="lg" className="rounded-full">Посмотреть все</Button>
      </div>
    </section>
  )
}
