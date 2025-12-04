export type CourseCategory = {
  id: number
  icon: string
  title: string
  description: string
  colorVar: string
}

export type Course = {
  id: number
  title: string
  image: string
  category: string
  price: string
  rating: number
  lessons: number
  duration: string
  students: number
  author: string
}
