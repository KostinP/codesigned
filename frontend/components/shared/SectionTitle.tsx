type Props = {
  title: string
  description: string
  className?: string
}

export function SectionTitle({ title, description, className = '' }: Props) {
  return (
    <div className={`mb-8 ${className}`}>
      <h2 className="text-2xl md:text-3xl font-bold text-foreground mb-4">{title}</h2>
      <p className="text-muted-foreground text-base leading-relaxed">{description}</p>
    </div>
  )
}
