import Image from 'next/image'

type Props = {
  icon: string
  title: string
  description: string
  colorVar: string
}

export function CourseCard({ icon, title, description, colorVar }: Props) {
  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-border flex flex-col gap-4">
      <div
        className="w-12 h-12 rounded-full flex items-center justify-center"
        style={{ backgroundColor: colorVar }}
      >
        <Image
          src={icon}
          alt={title}
          width={24}
          height={24}
          className="object-contain"
          style={{ width: 'auto', height: 'auto' }}
        />
      </div>
      <h3 className="text-lg font-semibold text-foreground">{title}</h3>
      <p className="text-muted-foreground text-sm leading-relaxed flex-1">{description}</p>
      <div className="mt-auto" />
    </div>
  )
}
