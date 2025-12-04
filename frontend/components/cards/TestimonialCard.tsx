import Image from 'next/image'
import { Star } from 'lucide-react'

type Props = {
  name: string
  avatar: string
  fact: string
  text: string
  rating: number
}

export function TestimonialCard({ name, avatar, fact, text, rating }: Props) {
  return (
    <div className="bg-white p-6 rounded-2xl border border-border shadow-sm flex flex-col gap-4">
      <div className="flex items-center gap-4">
        <Image src={avatar} alt={name} width={48} height={48} className="rounded-full object-cover" />
        <div>
          <p className="font-semibold text-foreground">{name}</p>
          <p className="text-sm text-muted-foreground">{fact}</p>
        </div>
      </div>
      <p className="text-sm text-muted-foreground">{text}</p>
      <div className="flex gap-1 text-yellow-400">
        {[...Array(5)].map((_, i) => (
          <Star key={i} className={`w-4 h-4 ${i < rating ? 'fill-yellow-400' : 'text-gray-300'}`} />
        ))}
      </div>
    </div>
  )
}
