"use client"

import { useEffect, useState } from "react"
import { CoursePreviewCard } from "@/components/cards/CoursePreviewCard"
import { Button } from "@/components/ui/button"

interface Course {
  id: number
  title: string
  image: string
  description: string
  lessonsCount: number
  duration: string
}

export function PopularCoursesSection() {
  const [courses, setCourses] = useState<Course[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    async function loadCourses() {
      try {
        // TODO: Заменить на реальный API вызов, когда бэкенд будет готов
        // const response = await fetch('/api/courses/popular')
        // const data = await response.json()
        
        // Временные данные до реализации API курсов
        const mockCourses: Course[] = [
          {
            id: 1,
            title: "Основы веб-разработки",
            image: "course1.svg",
            description: "Изучите основы HTML, CSS и JavaScript для создания современных веб-приложений",
            lessonsCount: 12,
            duration: "4ч 30м"
          },
          {
            id: 2,
            title: "Введение в Data Science",
            image: "course2.svg",
            description: "Основы анализа данных, машинного обучения и визуализации результатов",
            lessonsCount: 15,
            duration: "6ч 15м"
          },
          {
            id: 3,
            title: "UI/UX дизайн для начинающих",
            image: "course3.svg",
            description: "Создавайте интуитивно понятные и красивые интерфейсы для веб и мобильных приложений",
            lessonsCount: 10,
            duration: "5ч 00м"
          },
        ]
        
        setCourses(mockCourses)
      } catch (err) {
        setError('Не удалось загрузить курсы')
        console.error('Error loading courses:', err)
      } finally {
        setLoading(false)
      }
    }
    loadCourses()
  }, [])

  return (
    <section className="py-20 px-4 bg-[var(--gray-bg)]">
      <div className="max-w-7xl mx-auto flex flex-col gap-12 items-center text-center">
        <h2 className="text-2xl md:text-3xl font-bold mb-4">Самые популярные курсы</h2>

        {error && (
          <p className="text-destructive">{error}</p>
        )}

        {loading ? (
          <p className="text-muted-foreground">Загрузка курсов...</p>
        ) : (
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6 w-full">
            {courses.map((course) => (
              <CoursePreviewCard 
                key={course.id}
                image={course.image}
                title={course.title}
                description={course.description}
                lessonsCount={course.lessonsCount}
                duration={course.duration}
              />
            ))}
          </div>
        )}

        <Button variant="outline" size="lg" className="rounded-full">
          Смотреть все курсы
        </Button>
      </div>
    </section>
  )
}