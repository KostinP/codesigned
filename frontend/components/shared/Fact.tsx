type Props = {
  title: string
  text: string
  colorVar: string
}

export function Fact({ title, text, colorVar }: Props) {
  return (
    <div className="flex flex-col items-start">
      <span className="text-xl font-bold" style={{ color: colorVar }}>{title}</span>
      <span>{text}</span>
    </div>
  )
}
