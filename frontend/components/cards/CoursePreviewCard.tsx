import Image from 'next/image'
import { Button } from '@/components/ui/button'

type Props = {
  image: string // относительный путь от public/courses
  title: string
  description: string
  lessonsCount: number
  duration: string
}

export function CoursePreviewCard({
  image,
  title,
  description,
  lessonsCount,
  duration,
}: Props) {
  return (
    <div className="bg-white rounded-2xl overflow-hidden shadow-sm border border-border flex flex-col">
      <div className="relative w-full h-40">
        <Image
          src={`/courses/${image}`}
          alt={title}
          fill
          className="object-cover"
        />
      </div>
      <div className="p-4 flex flex-col gap-2 flex-1">
        <h3 className="text-base font-semibold text-foreground">{title}</h3>
        <p className="text-sm text-muted-foreground line-clamp-2">{description}</p>
        <div className="text-xs text-muted-foreground mt-auto flex justify-between pt-2">
          <span>{lessonsCount} уроков</span>
          <span>{duration}</span>
        </div>
        <Button size="sm" className="mt-3 w-full rounded-full">
          Начать
        </Button>
      </div>
    </div>
  )
}
