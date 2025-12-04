import { CourseCategory } from "@/entities/course/models"
import { COLORS } from '@/lib/constants'

export const courseCategories: CourseCategory[] = [
  {
    id: 1,
    icon: "/icons/category-web_development.svg",
    title: "Веб-разработка",
    description: "Создавайте сайты и приложения",
    colorVar: COLORS.pinkText,
  },
  {
    id: 2,
    title: "Дизайн",
    description: "Освойте UI/UX и визуальное мышление",
    icon: "/icons/category-design.svg",
    colorVar: COLORS.blueText,
  },
  {
    id: 3,
    title: 'Маркетинг',
    description: "Анализируйте и продвигайте продукты",
    icon: "/icons/category-marketing.svg",
    colorVar: COLORS.greenText,
  },
  {
    id: 4,
    title: "Геймдев",
    description: "Создавайте игры и симуляции",
    icon: "/icons/category-game_development.svg",
    colorVar: COLORS.greenText,
  },
]
